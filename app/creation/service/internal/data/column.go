package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

var _ biz.ColumnRepo = (*columnRepo)(nil)

type columnRepo struct {
	data *Data
	log  *log.Helper
}

func NewColumnRepo(data *Data, logger log.Logger) biz.ColumnRepo {
	return &columnRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/column")),
	}
}

func (r *columnRepo) GetColumn(ctx context.Context, id int32) (*biz.Column, error) {
	column := &Column{}
	err := r.data.db.WithContext(ctx).Where("column_id = ?", id).First(column).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column from db: id(%v)", id))
	}
	return &biz.Column{
		ColumnId: id,
		Uuid:     column.Uuid,
	}, nil
}

func (r *columnRepo) GetLastColumnDraft(ctx context.Context, uuid string) (*biz.ColumnDraft, error) {
	draft := &ColumnDraft{}
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("column draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
	}
	return &biz.ColumnDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *columnRepo) CreateColumnDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &ColumnDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create an column draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *columnRepo) CreateColumnFolder(ctx context.Context, id int32, uuid string) error {
	name := "column/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create an column folder: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) SendColumn(ctx context.Context, id int32, uuid string) (*biz.ColumnDraft, error) {
	cd := &ColumnDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&ColumnDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 1).Updates(cd).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.ColumnDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *columnRepo) CreateColumn(ctx context.Context, id, auth int32, uuid string) error {
	column := &Column{
		ColumnId: id,
		Uuid:     uuid,
		Auth:     auth,
	}
	err := r.data.DB(ctx).Select("ColumnId", "Uuid", "Auth").Create(column).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create an column: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) AddColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error {
	inclusion := &ColumnInclusion{
		ColumnId:  id,
		ArticleId: articleId,
		Uuid:      uuid,
		Status:    1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(inclusion).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add an article to column: uuid(%s), id(%v), articleId(%v)", uuid, id, articleId))
	}
	return nil
}

func (r *columnRepo) AddColumnIncludesToCache(ctx context.Context, id, articleId int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	articleIds := strconv.Itoa(int(articleId))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {

		pipe.ZAddNX(ctx, "column_includes_"+ids, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: articleIds + "%" + uuid,
		})
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column includes to cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) DeleteColumnIncludes(ctx context.Context, id, articleId int32, uuid string) error {
	cc := &ColumnInclusion{
		Status: 2,
	}
	err := r.data.db.WithContext(ctx).Model(&ColumnInclusion{}).Where("column_id = ? and article_id = ? and uuid = ?", id, articleId, uuid).Updates(cc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an article from column: uuid(%s), id(%v), articleId(%v)", uuid, id, articleId))
	}
	return nil
}

func (r *columnRepo) DeleteColumnIncludesFromCache(ctx context.Context, id, articleId int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	articleIds := strconv.Itoa(int(articleId))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, "column_includes_"+ids, articleIds+"%"+uuid)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete column includes from cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) DeleteColumnDraft(ctx context.Context, id int32, uuid string) error {
	cd := &ColumnDraft{}
	cd.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(cd).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete column draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) CreateColumnStatistic(ctx context.Context, id, auth int32, uuid string) error {
	cs := &ColumnStatistic{
		ColumnId: id,
		Uuid:     uuid,
		Auth:     auth,
	}
	err := r.data.DB(ctx).Select("ColumnId", "Uuid", "Auth").Create(cs).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create an column statistic: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) SendColumnToMq(ctx context.Context, column *biz.Column, mode string) error {
	columnMap := map[string]interface{}{}
	columnMap["uuid"] = column.Uuid
	columnMap["id"] = column.ColumnId
	columnMap["auth"] = column.Auth
	columnMap["mode"] = mode

	data, err := json.Marshal(columnMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "column",
		Body:  data,
	}
	msg.WithKeys([]string{column.Uuid})
	_, err = r.data.columnMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send column to mq: %v", column))
	}
	return nil
}

func (r *columnRepo) SendReviewToMq(ctx context.Context, review *biz.ColumnReview) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "column_review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.columnReviewMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *columnRepo) CreateColumnCache(ctx context.Context, id, auth int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "column_"+ids, "uuid", uuid)
		pipe.HSetNX(ctx, "column_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "column_"+ids, "collect", 0)
		pipe.HSetNX(ctx, "column_"+ids, "view", 0)

		if auth == 2 {
			return nil
		}

		pipe.ZAddNX(ctx, "column", &redis.Z{
			Score:  float64(id),
			Member: ids + "%" + uuid,
		})
		pipe.ZAddNX(ctx, "column_hot", &redis.Z{
			Score:  0,
			Member: ids + "%" + uuid,
		})
		pipe.ZAddNX(ctx, "leaderboard", &redis.Z{
			Score:  0,
			Member: ids + "%" + uuid + "%column",
		})
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create(update) column cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) UpdateColumnCache(ctx context.Context, id, auth int32, uuid string) error {
	return r.CreateColumnCache(ctx, id, auth, uuid)
}

func (r *columnRepo) EditColumnCos(ctx context.Context, id int32, uuid string) error {
	err := r.EditColumnCosContent(ctx, id, uuid)
	if err != nil {
		return err
	}

	err = r.EditColumnCosIntroduce(ctx, id, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *columnRepo) EditColumnCosContent(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "column/" + uuid + "/" + ids + "/content-edit"
	dest := "column/" + uuid + "/" + ids + "/content"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy column from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) EditColumnCosIntroduce(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "column/" + uuid + "/" + ids + "/introduce-edit"
	dest := "column/" + uuid + "/" + ids + "/introduce"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy column from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) CreateColumnSearch(ctx context.Context, id int32, uuid string) error {
	return nil
}

func (r *columnRepo) GetColumnList(ctx context.Context, page int32) ([]*biz.Column, error) {
	column, err := r.getColumnFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(column)
	if size != 0 {
		return column, nil
	}

	column, err = r.getColumnFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(column)
	if size != 0 {
		go r.setColumnToCache("column", column)
	}
	return column, nil
}

func (r *columnRepo) getColumnFromCache(ctx context.Context, page int32) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "column", index*10, index+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column from cache: key(%s), page(%v)", "column", page))
	}

	column := make([]*biz.Column, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		column = append(column, &biz.Column{
			ColumnId: int32(id),
			Uuid:     member[1],
		})
	}
	return column, nil
}

func (r *columnRepo) getColumnFromDB(ctx context.Context, page int32) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Column, 0)
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("column_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column from db: page(%v)", page))
	}

	column := make([]*biz.Column, 0)
	for _, item := range list {
		column = append(column, &biz.Column{
			ColumnId: item.ColumnId,
			Uuid:     item.Uuid,
		})
	}
	return column, nil
}

func (r *columnRepo) setColumnToCache(key string, column []*biz.Column) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range column {
			z = append(z, &redis.Z{
				Score:  float64(item.ColumnId),
				Member: strconv.Itoa(int(item.ColumnId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set column to cache: column(%v)", column)
	}
}

func (r *columnRepo) GetColumnListHot(ctx context.Context, page int32) ([]*biz.ColumnStatistic, error) {
	column, err := r.getColumnHotFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(column)
	if size != 0 {
		return column, nil
	}

	column, err = r.getColumnHotFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(column)
	if size != 0 {
		go r.setColumnHotToCache("column_hot", column)
	}
	return column, nil
}

func (r *columnRepo) GetUserColumnList(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Column, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("column_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column from db: page(%v), uuid(%s)", page, uuid))
	}

	column := make([]*biz.Column, 0)
	for _, item := range list {
		column = append(column, &biz.Column{
			ColumnId: item.ColumnId,
			Uuid:     item.Uuid,
		})
	}
	return column, nil
}

func (r *columnRepo) GetUserColumnListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Column, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and auth = ?", uuid, 1).Order("column_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column visitor from db: page(%v), uuid(%s)", page, uuid))
	}

	column := make([]*biz.Column, 0)
	for _, item := range list {
		column = append(column, &biz.Column{
			ColumnId: item.ColumnId,
			Uuid:     item.Uuid,
		})
	}
	return column, nil
}

func (r *columnRepo) GetColumnCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Column{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get column count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *columnRepo) GetColumnCountVisitor(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Column{}).Where("uuid = ? and auth = ?", uuid, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get column count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *columnRepo) getColumnHotFromCache(ctx context.Context, page int32) ([]*biz.ColumnStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "column_hot", index*10, index+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column hot from cache: key(%s), page(%v)", "column_hot", page))
	}

	column := make([]*biz.ColumnStatistic, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		column = append(column, &biz.ColumnStatistic{
			ColumnId: int32(id),
			Uuid:     member[1],
		})
	}
	return column, nil
}

func (r *columnRepo) getColumnHotFromDB(ctx context.Context, page int32) ([]*biz.ColumnStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ColumnStatistic, 0)
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("agree desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column statistic from db: page(%v)", page))
	}

	column := make([]*biz.ColumnStatistic, 0)
	for _, item := range list {
		column = append(column, &biz.ColumnStatistic{
			ColumnId: item.ColumnId,
			Uuid:     item.Uuid,
			Agree:    item.Agree,
		})
	}
	return column, nil
}

func (r *columnRepo) setColumnHotToCache(key string, column []*biz.ColumnStatistic) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range column {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.ColumnId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set column to cache: column(%v)", column)
	}
}

func (r *columnRepo) GetColumnListStatistic(ctx context.Context, ids []int32) ([]*biz.ColumnStatistic, error) {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range ids {
			pipe.HMGet(ctx, "column_"+strconv.Itoa(int(id)), "agree", "collect", "view")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column list statistic from cache: ids(%v)", ids))
	}

	statistic := make([]*biz.ColumnStatistic, 0)
	for index, item := range cmd {
		val := []int32{0, 0, 0}
		for _index, count := range item.(*redis.SliceCmd).Val() {
			if count == nil {
				break
			}
			num, err := strconv.ParseInt(count.(string), 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
			}
			val[_index] = int32(num)
		}
		statistic = append(statistic, &biz.ColumnStatistic{
			ColumnId: ids[index],
			Agree:    val[0],
			Collect:  val[1],
			View:     val[2],
		})
	}
	return statistic, nil
}

func (r *columnRepo) GetColumnStatistic(ctx context.Context, id int32) (*biz.ColumnStatistic, error) {
	var statistic *biz.ColumnStatistic
	key := "column_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		statistic, err = r.getColumnStatisticFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}

	statistic, err = r.getColumnStatisticFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	go r.setColumnStatisticToCache(key, statistic)

	return statistic, nil
}

func (r *columnRepo) getColumnStatisticFromCache(ctx context.Context, key string) (*biz.ColumnStatistic, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "uuid", "agree", "collect", "view").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column statistic form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0}
	for _index, count := range statistic[1:] {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.ColumnStatistic{
		Uuid:    statistic[0].(string),
		Agree:   val[0],
		Collect: val[1],
		View:    val[2],
	}, nil
}

func (r *columnRepo) getColumnStatisticFromDB(ctx context.Context, id int32) (*biz.ColumnStatistic, error) {
	cs := &ColumnStatistic{}
	err := r.data.db.WithContext(ctx).Where("column_id = ?", id).First(cs).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get statistic from db: id(%v)", id))
	}
	return &biz.ColumnStatistic{
		Uuid:    cs.Uuid,
		Agree:   cs.Agree,
		Collect: cs.Collect,
		View:    cs.View,
	}, nil
}

func (r *columnRepo) setColumnStatisticToCache(key string, statistic *biz.ColumnStatistic) {
	err := r.data.redisCli.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View).Err()
	if err != nil {
		r.log.Errorf("fail to set column statistic to cache, err(%s)", err.Error())
	}
}

func (r *columnRepo) DeleteColumn(ctx context.Context, id int32, uuid string) error {
	column := &Column{}
	err := r.data.DB(ctx).Where("column_id = ? and uuid = ?", id, uuid).Delete(column).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an column: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) DeleteColumnStatistic(ctx context.Context, id int32, uuid string) error {
	statistic := &ColumnStatistic{}
	err := r.data.DB(ctx).Where("column_id = ? and uuid = ?", id, uuid).Delete(statistic).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an column statistic: uuid(%s)", uuid))
	}
	return nil
}

func (r *columnRepo) DeleteColumnCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, "column", ids+"%"+uuid)
		pipe.ZRem(ctx, "column_hot", ids+"%"+uuid)
		pipe.ZRem(ctx, "leaderboard", ids+"%"+uuid+"%column")
		pipe.Del(ctx, "column_"+ids)
		pipe.Del(ctx, "column_collect_"+ids)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete column cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) FreezeColumnCos(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "column/" + uuid + "/" + ids + "/content"
	_, err := r.data.cosCli.Object.Delete(ctx, key)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to freeze column: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) GetColumnAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "column_agree_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge column agree member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *columnRepo) GetColumnCollectJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "column_collect_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge column collect member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *columnRepo) GetSubscribeList(ctx context.Context, page int32, uuid string) ([]*biz.Subscribe, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Subscribe, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("column_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get subscribe column from db: page(%v)", page))
	}

	subscribe := make([]*biz.Subscribe, 0)
	for _, item := range list {
		subscribe = append(subscribe, &biz.Subscribe{
			ColumnId: item.ColumnId,
			AuthorId: item.AuthorId,
		})
	}
	return subscribe, nil
}

func (r *columnRepo) GetSubscribeListCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Subscribe{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get subscribe column list count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *columnRepo) GetColumnSubscribes(ctx context.Context, uuid string, ids []int32) ([]*biz.Subscribe, error) {
	list := make([]*Subscribe, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and column_id IN ?", uuid, ids).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column subscribe from db: uuid(%s), ids(%v)", uuid, ids))
	}

	subscribes := make([]*biz.Subscribe, 0)
	for _, item := range list {
		subscribes = append(subscribes, &biz.Subscribe{
			ColumnId: item.ColumnId,
			Status:   item.Status,
		})
	}
	return subscribes, nil
}

func (r *columnRepo) DeleteColumnSearch(ctx context.Context, id int32, uuid string) error {
	return nil
}

func (r *columnRepo) EditColumnSearch(ctx context.Context, id int32, uuid string) error {
	return nil
}

func (r *columnRepo) SetColumnAgree(ctx context.Context, id int32, uuid string) error {
	cs := ColumnStatistic{}
	err := r.data.DB(ctx).Model(&cs).Where("column_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column agree: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) SetColumnView(ctx context.Context, id int32, uuid string) error {
	cs := ColumnStatistic{}
	err := r.data.DB(ctx).Model(&cs).Where("column_id = ? and uuid = ?", id, uuid).Update("view", gorm.Expr("view + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column view: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) SetColumnViewToCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "column_"+ids, "view", 1)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add column agree to cache: id(%v), uuid(%s)", id, uuid)
	}
	return nil
}

func (r *columnRepo) SetColumnAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "column_"+ids, "agree", 1)
		pipe.ZIncrBy(ctx, "column_hot", 1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", 1, ids+"%"+uuid+"%column")
		pipe.SAdd(ctx, "column_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add column agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *columnRepo) SendColumnStatisticToMq(ctx context.Context, uuid, mode string) error {
	achievement := map[string]interface{}{}
	achievement["uuid"] = uuid
	achievement["mode"] = mode

	data, err := json.Marshal(achievement)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "achievement",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.achievementMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send column statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *columnRepo) SetColumnUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
	collect := &Collect{
		CollectionsId: collectionsId,
		Uuid:          userUuid,
		CreationsId:   id,
		Mode:          3,
		Status:        1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to collect an column: column_id(%v), collectionsId(%v), userUuid(%s)", id, collectionsId, userUuid))
	}
	return nil
}

func (r *columnRepo) SetColumnCollect(ctx context.Context, id int32, uuid string) error {
	cs := ColumnStatistic{}
	err := r.data.DB(ctx).Model(&cs).Where("column_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column collect: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) SetColumnCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "column_"+ids, "collect", 1)
		pipe.SAdd(ctx, "column_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add column collect to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *columnRepo) CancelColumnAgree(ctx context.Context, id int32, uuid string) error {
	cs := ColumnStatistic{}
	err := r.data.DB(ctx).Model(&cs).Where("column_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column agree: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) CancelColumnAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "column_"+ids, "agree", -1)
		pipe.ZIncrBy(ctx, "column_hot", -1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", -1, ids+"%"+uuid+"%column")
		pipe.SRem(ctx, "column_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel column agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *columnRepo) CancelColumnCollect(ctx context.Context, id int32, uuid string) error {
	cs := &ColumnStatistic{}
	err := r.data.DB(ctx).Model(cs).Where("column_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column collect: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) CancelColumnCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "column_"+ids, "collect", -1)
		pipe.SRem(ctx, "column_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel column collect from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *columnRepo) CancelColumnUserCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{}
	err := r.data.DB(ctx).Where("creations_id = ? and uuid = ? and mode = ?", id, userUuid, 3).Delete(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column collect: column_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) SubscribeColumn(ctx context.Context, id int32, author, uuid string) error {
	sub := &Subscribe{
		ColumnId: id,
		AuthorId: author,
		Uuid:     uuid,
		Status:   1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(sub).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to subscribe column: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) CancelSubscribeColumn(ctx context.Context, id int32, uuid string) error {
	sub := &Subscribe{}
	err := r.data.DB(ctx).Model(sub).Where("column_id = ? and uuid = ?", id, uuid).Update("status", 2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel subscribe column: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) SubscribeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	f := &Subscribe{}
	err := r.data.db.WithContext(ctx).Where("column_id = ? and uuid = ? and status = ?", id, uuid, 1).First(f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to get column subscribe judge rom db: id(%v), uuid(%s)", id, uuid))
	}
	return true, nil
}
