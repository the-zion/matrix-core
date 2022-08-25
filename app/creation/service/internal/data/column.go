package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
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
	var script = redis.NewScript(`
					local key = KEYS[1]
					local member = KEYS[2]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("ZADD", key, change, member)
					end
					return 0
	`)
	keys := []string{"column_includes_" + ids, articleIds + "%" + uuid}
	values := []interface{}{float64(time.Now().Unix())}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column includes to cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *columnRepo) AddCreationUserColumn(ctx context.Context, uuid string, auth int32) error {
	cu := &CreationUser{
		Uuid:   uuid,
		Column: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"`column`": gorm.Expr("`column` + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation column: uuid(%v)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := &CreationUserVisitor{
		Uuid:   uuid,
		Column: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"`column`": gorm.Expr("`column` + ?", 1)}),
	}).Create(cuv).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation column visitor: uuid(%v)", uuid))
	}
	return nil
}

func (r *columnRepo) ReduceCreationUserColumn(ctx context.Context, auth int32, uuid string) error {
	cu := CreationUser{}
	err := r.data.DB(ctx).Model(&cu).Where("uuid = ? and `column` > 0", uuid).Update("`column`", gorm.Expr("`column` - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user column: uuid(%s)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := CreationUserVisitor{}
	err = r.data.DB(ctx).Model(&cuv).Where("uuid = ? and `column` > 0", uuid).Update("`column`", gorm.Expr("`column` - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user column visitor: uuid(%s)", uuid))
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

func (r *columnRepo) SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error {
	scoreMap := map[string]interface{}{}
	scoreMap["uuid"] = uuid
	scoreMap["score"] = score
	scoreMap["mode"] = mode

	data, err := json.Marshal(scoreMap)
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
		return errors.Wrapf(err, fmt.Sprintf("fail to send score to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *columnRepo) SendStatisticToMq(ctx context.Context, id, collectionsId int32, uuid, userUuid, mode string) error {
	statisticMap := map[string]interface{}{}
	statisticMap["id"] = id
	statisticMap["collectionsId"] = collectionsId
	statisticMap["uuid"] = uuid
	statisticMap["userUuid"] = userUuid
	statisticMap["mode"] = mode

	data, err := json.Marshal(statisticMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "column",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.columnMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send column statistic to mq: map(%v)", statisticMap))
	}
	return nil
}

func (r *columnRepo) SendColumnIncludesToMq(ctx context.Context, id, articleId int32, uuid, mode string) error {
	includesMap := map[string]interface{}{}
	includesMap["id"] = id
	includesMap["articleId"] = articleId
	includesMap["uuid"] = uuid
	includesMap["mode"] = mode

	data, err := json.Marshal(includesMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "column",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.columnMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column includes to mq: map(%v)", includesMap))
	}
	return nil
}

func (r *columnRepo) SendColumnSubscribeToMq(ctx context.Context, id int32, uuid, mode string) error {
	subscribeMap := map[string]interface{}{}
	subscribeMap["id"] = id
	subscribeMap["uuid"] = uuid
	subscribeMap["mode"] = mode

	data, err := json.Marshal(subscribeMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "column",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.columnMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set column subscribe to mq: map(%v)", subscribeMap))
	}
	return nil
}

func (r *columnRepo) CreateColumnCache(ctx context.Context, id, auth int32, uuid string) error {
	exists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, "column")
		pipe.Exists(ctx, "column_hot")
		pipe.Exists(ctx, "leaderboard")
		pipe.Exists(ctx, "user_column_list_"+uuid)
		pipe.Exists(ctx, "user_column_list_visitor_"+uuid)
		pipe.Exists(ctx, "creation_user_"+uuid)
		pipe.Exists(ctx, "creation_user_visitor_"+uuid)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if column exist from cache: id(%v),uuid(%s)", id, uuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, int32(exist))
	}

	ids := strconv.Itoa(int(id))
	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "column_"+ids, "uuid", uuid)
		pipe.HSetNX(ctx, "column_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "column_"+ids, "collect", 0)
		pipe.HSetNX(ctx, "column_"+ids, "view", 0)

		if exists[3] == 1 {
			pipe.ZAddNX(ctx, "user_column_list_"+uuid, &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid,
			})
		}

		if exists[5] == 1 {
			pipe.HIncrBy(ctx, "creation_user_"+uuid, "column", 1)
		}

		if auth == 2 {
			return nil
		}

		if exists[0] == 1 {
			pipe.ZAddNX(ctx, "column", &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid,
			})
		}

		if exists[1] == 1 {
			pipe.ZAddNX(ctx, "column_hot", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid,
			})
		}

		if exists[2] == 1 {
			pipe.ZAddNX(ctx, "leaderboard", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid + "%column",
			})
		}

		if exists[4] == 1 {
			pipe.ZAddNX(ctx, "user_column_list_visitor_"+uuid, &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid + "%column",
			})
		}

		if exists[6] == 1 {
			pipe.HIncrBy(ctx, "creation_user_visitor_"+uuid, "column", 1)
		}
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
	ids := strconv.Itoa(int(id))
	key := "column/" + uuid + "/" + ids + "/search"
	resp, err := r.data.cosCli.Object.Get(ctx, key, &cos.ObjectGetOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get column from cos: id(%v), uuid(%s)", id, uuid))
	}

	column, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read request body: id(%v), uuid(%s)", id, uuid))
	}

	resp.Body.Close()

	req := esapi.IndexRequest{
		Index:      "column",
		DocumentID: strconv.Itoa(int(id)),
		Body:       bytes.NewReader(column),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("error getting column search create response: id(%v), uuid(%s)", id, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: id(%v), uuid(%s)", id, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error indexing document to es: reason(%v), id(%v), uuid(%s)", e, id, uuid))
		}
	}
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
	list, err := r.data.redisCli.ZRevRange(ctx, "column", index*10, index*10+9).Result()
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
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range column {
			z = append(z, &redis.Z{
				Score:  float64(item.ColumnId),
				Member: strconv.Itoa(int(item.ColumnId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
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

	column, err = r.GetColumnHotFromDB(ctx, page)
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
	column, err := r.getUserColumnListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(column)
	if size != 0 {
		return column, nil
	}

	column, err = r.getUserColumnListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(column)
	if size != 0 {
		go r.setUserColumnListToCache("user_column_list_"+uuid, column)
	}
	return column, nil
}

func (r *columnRepo) getUserColumnListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_column_list_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column list visitor from cache: key(%s), page(%v)", "user_column_list_"+uuid, page))
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

func (r *columnRepo) getUserColumnListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
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
	column, err := r.getUserColumnListVisitorFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(column)
	if size != 0 {
		return column, nil
	}

	column, err = r.getUserColumnListVisitorFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(column)
	if size != 0 {
		go r.setUserColumnListToCache("user_column_list_visitor_"+uuid, column)
	}
	return column, nil
}

func (r *columnRepo) getUserColumnListVisitorFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_column_list_visitor_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column list visitor from cache: key(%s), page(%v)", "user_column_list_visitor_"+uuid, page))
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

func (r *columnRepo) getUserColumnListVisitorFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Column, error) {
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

func (r *columnRepo) setUserColumnListToCache(key string, column []*biz.Column) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range column {
			z = append(z, &redis.Z{
				Score:  float64(item.ColumnId),
				Member: strconv.Itoa(int(item.ColumnId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		pipe.Expire(context.Background(), key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set column to cache: column(%v)", column)
	}
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
	list, err := r.data.redisCli.ZRevRange(ctx, "column_hot", index*10, index*10+9).Result()
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

func (r *columnRepo) GetColumnHotFromDB(ctx context.Context, page int32) ([]*biz.ColumnStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ColumnStatistic, 0)
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("agree desc, column_id desc").Offset(index * 10).Limit(10).Find(&list).Error
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
	columnListStatistic := make([]*biz.ColumnStatistic, 0)
	exists, unExists, err := r.columnListStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getColumnListStatisticFromCache(ctx, exists, &columnListStatistic)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getColumnListStatisticFromDb(ctx, unExists, &columnListStatistic)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return columnListStatistic, nil
}

func (r *columnRepo) columnListStatisticExist(ctx context.Context, ids []int32) ([]int32, []int32, error) {
	exists := make([]int32, 0)
	unExists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "column_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if column statistic exist from cache: ids(%v)", ids))
	}

	for index, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		if exist == 1 {
			exists = append(exists, ids[index])
		} else {
			unExists = append(unExists, ids[index])
		}
	}
	return exists, unExists, nil
}

func (r *columnRepo) getColumnListStatisticFromCache(ctx context.Context, exists []int32, columnListStatistic *[]*biz.ColumnStatistic) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range exists {
			pipe.HMGet(ctx, "column_"+strconv.Itoa(int(id)), "agree", "collect", "view")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get column list statistic from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		val := []int32{0, 0, 0}
		for _index, count := range item.(*redis.SliceCmd).Val() {
			if count == nil {
				break
			}
			num, err := strconv.ParseInt(count.(string), 10, 32)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
			}
			val[_index] = int32(num)
		}
		*columnListStatistic = append(*columnListStatistic, &biz.ColumnStatistic{
			ColumnId: exists[index],
			Agree:    val[0],
			Collect:  val[1],
			View:     val[2],
		})
	}
	return nil
}

func (r *columnRepo) getColumnListStatisticFromDb(ctx context.Context, unExists []int32, columnListStatistic *[]*biz.ColumnStatistic) error {
	list := make([]*ColumnStatistic, 0)
	err := r.data.db.WithContext(ctx).Where("column_id IN ?", unExists).Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get column statistic list from db: ids(%v)", unExists))
	}

	for _, item := range list {
		*columnListStatistic = append(*columnListStatistic, &biz.ColumnStatistic{
			ColumnId: item.ColumnId,
			Uuid:     item.Uuid,
			Agree:    item.Agree,
			Collect:  item.Collect,
			View:     item.View,
		})
	}

	if len(list) != 0 {
		go r.setColumnListStatisticToCache(list)
	}

	return nil
}

func (r *columnRepo) setColumnListStatisticToCache(commentList []*ColumnStatistic) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "column_" + strconv.Itoa(int(item.ColumnId))
			pipe.HSetNX(context.Background(), key, "uuid", item.Uuid)
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.HSetNX(context.Background(), key, "collect", item.Collect)
			pipe.HSetNX(context.Background(), key, "view", item.View)
			pipe.Expire(context.Background(), key, time.Hour*8)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set column statistic to cache, err(%s)", err.Error())
	}
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
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View).Err()
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
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

func (r *columnRepo) DeleteColumnCache(ctx context.Context, id, auth int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	var script = redis.NewScript(`
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]
					local key7 = KEYS[7]
					local key8 = KEYS[8]
					local key9 = KEYS[9]

                    local member = ARGV[1]
					local leadMember = ARGV[2]
					local auth = ARGV[3]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key2, member)
					redis.call("ZREM", key3, leadMember)
					redis.call("DEL", key4)
					redis.call("DEL", key5)

                    redis.call("ZREM", key6, member)

                    local exist = redis.call("EXISTS", key8)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key8, "column"))
						if number > 0 then
  							redis.call("HINCRBY", key8, "column", -1)
						end
					end

                    if auth == 2 then
						return 0
					end

					redis.call("ZREM", key7, member)

					local exist = redis.call("EXISTS", key9)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key9, "column"))
						if number > 0 then
  							redis.call("HINCRBY", key9, "column", -1)
						end
					end

					return 0
	`)
	keys := []string{"column", "column_hot", "leaderboard", "column_" + ids, "column_collect_" + ids, "user_column_list_" + uuid, "user_column_list_visitor_" + uuid, "creation_user_" + uuid, "creation_user_visitor_" + uuid}
	values := []interface{}{ids + "%" + uuid, ids + "%" + uuid + "%column", auth}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
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
	subscribe, err := r.getUserSubscribeListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(subscribe)
	if size != 0 {
		return subscribe, nil
	}

	subscribe, err = r.getUserSubscribeListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(subscribe)
	if size != 0 {
		go r.setUserSubscribeListToCache("user_column_subscribe_list_"+uuid, subscribe)
	}
	return subscribe, nil
}

func (r *columnRepo) getUserSubscribeListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Subscribe, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_column_subscribe_list_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column subscribe list from cache: key(%s), page(%v)", "user_column_subscribe_list_"+uuid, page))
	}

	subscribe := make([]*biz.Subscribe, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		subscribe = append(subscribe, &biz.Subscribe{
			ColumnId: int32(id),
			AuthorId: member[1],
		})
	}
	return subscribe, nil
}

func (r *columnRepo) getUserSubscribeListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Subscribe, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Subscribe, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Order("column_id desc").Offset(index * 10).Limit(10).Find(&list).Error
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

func (r *columnRepo) setUserSubscribeListToCache(key string, subscribe []*biz.Subscribe) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range subscribe {
			z = append(z, &redis.Z{
				Score:  float64(item.ColumnId),
				Member: strconv.Itoa(int(item.ColumnId)) + "%" + item.AuthorId,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		pipe.Expire(context.Background(), key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set column subscribe to cache: column(%v)", subscribe)
	}
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

func (r *columnRepo) GetColumnSearch(ctx context.Context, page int32, search, time string) ([]*biz.ColumnSearch, int32, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	reply := make([]*biz.ColumnSearch, 0)

	var buf bytes.Buffer
	var body map[string]interface{}
	if search != "" {
		body = map[string]interface{}{
			"from":    index * 10,
			"size":    10,
			"_source": []string{"update", "tags", "cover", "uuid"},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"multi_match": map[string]interface{}{
							"query":  search,
							"fields": []string{"introduce", "tags", "name"},
						}},
					},
					"filter": map[string]interface{}{
						"range": map[string]interface{}{
							"update": map[string]interface{}{},
						},
					},
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					"introduce": map[string]interface{}{
						"fragment_size":       300,
						"number_of_fragments": 1,
						"no_match_size":       300,
						"pre_tags":            "<span style='color:red'>",
						"post_tags":           "</span>",
					},
					"name": map[string]interface{}{
						"pre_tags":      "<span style='color:red'>",
						"post_tags":     "</span>",
						"no_match_size": 100,
					},
				},
			},
		}
	} else {
		body = map[string]interface{}{
			"from":    index * 10,
			"size":    10,
			"_source": []string{"update", "tags", "cover", "uuid"},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"match_all": map[string]interface{}{}},
					},
					"filter": map[string]interface{}{
						"range": map[string]interface{}{
							"update": map[string]interface{}{},
						},
					},
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					"introduce": map[string]interface{}{
						"fragment_size":       300,
						"number_of_fragments": 1,
						"no_match_size":       300,
						"pre_tags":            "<span style='color:red'>",
						"post_tags":           "</span>",
					},
					"name": map[string]interface{}{
						"pre_tags":      "<span style='color:red'>",
						"post_tags":     "</span>",
						"no_match_size": 100,
					},
				},
			},
			"sort": []map[string]interface{}{
				{"_id": "desc"},
			},
		}
	}

	switch time {
	case "1day":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1d"
		break
	case "1week":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1w"
		break
	case "1month":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1M"
		break
	case "1year":
		body["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].(map[string]interface{})["range"].(map[string]interface{})["update"].(map[string]interface{})["gte"] = "now-1y"
		break
	}
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error encoding query: page(%v), search(%s), time(%s)", page, search, time))
	}

	res, err := r.data.elasticSearch.es.Search(
		r.data.elasticSearch.es.Search.WithContext(ctx),
		r.data.elasticSearch.es.Search.WithIndex("column"),
		r.data.elasticSearch.es.Search.WithBody(&buf),
		r.data.elasticSearch.es.Search.WithTrackTotalHits(true),
		r.data.elasticSearch.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error getting response from es: page(%v), search(%s), time(%s)", page, search, time))
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, errors.Wrapf(err, fmt.Sprintf("error parsing the response body: page(%v), search(%s), time(%s)", page, search, time))
		} else {
			return nil, 0, errors.Errorf(fmt.Sprintf("error search column from es: reason(%v), page(%v), search(%s), time(%s)", e, page, search, time))
		}
	}

	result := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error parsing the response body: page(%v), search(%s), time(%s)", page, search, time))
	}

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		id, err := strconv.ParseInt(hit.(map[string]interface{})["_id"].(string), 10, 32)
		if err != nil {
			return nil, 0, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: page(%v), search(%s), time(%s)", page, search, time))
		}

		reply = append(reply, &biz.ColumnSearch{
			Id:        int32(id),
			Tags:      hit.(map[string]interface{})["_source"].(map[string]interface{})["tags"].(string),
			Update:    hit.(map[string]interface{})["_source"].(map[string]interface{})["update"].(string),
			Cover:     hit.(map[string]interface{})["_source"].(map[string]interface{})["cover"].(string),
			Uuid:      hit.(map[string]interface{})["_source"].(map[string]interface{})["uuid"].(string),
			Introduce: hit.(map[string]interface{})["highlight"].(map[string]interface{})["introduce"].([]interface{})[0].(string),
			Name:      hit.(map[string]interface{})["highlight"].(map[string]interface{})["name"].([]interface{})[0].(string),
		})
	}
	res.Body.Close()
	return reply, int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), nil
}

func (r *columnRepo) GetColumnAuth(ctx context.Context, id int32) (int32, error) {
	column := &Column{}
	err := r.data.db.WithContext(ctx).Where("column_id = ?", id).First(column).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get column auth from db: id(%v)", id))
	}
	return column.Auth, nil
}

func (r *columnRepo) DeleteColumnSearch(ctx context.Context, id int32, uuid string) error {
	req := esapi.DeleteRequest{
		Index:      "column",
		DocumentID: strconv.Itoa(int(id)),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting column search delete response: id(%v), uuid(%s)", id, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: id(%v), uuid(%s)", id, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error delete document to es: reason(%v), id(%v), uuid(%s)", e, id, uuid))
		}
	}
	return nil
}

func (r *columnRepo) EditColumnSearch(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "column/" + uuid + "/" + ids + "/search"
	resp, err := r.data.cosCli.Object.Get(ctx, key, &cos.ObjectGetOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get column from cos: id(%v), uuid(%s)", id, uuid))
	}

	column, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read request body: id(%v), uuid(%s)", id, uuid))
	}

	resp.Body.Close()

	req := esapi.IndexRequest{
		Index:      "column",
		DocumentID: strconv.Itoa(int(id)),
		Body:       bytes.NewReader(column),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting column search edit response: id(%v), uuid(%s)", id, uuid))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: id(%v), uuid(%s)", id, uuid))
		} else {
			return errors.Errorf(fmt.Sprintf("error indexing document to es: reason(%v), id(%v), uuid(%s)", e, id, uuid))
		}
	}
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

func (r *columnRepo) SetUserColumnAgree(ctx context.Context, id int32, userUuid string) error {
	cg := &ColumnAgree{
		ColumnId: id,
		Uuid:     userUuid,
		Status:   1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(cg).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user column agree: id(%v), userUuid(%s)", id, userUuid))
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
	hotKey := fmt.Sprintf("column_hot")
	statisticKey := fmt.Sprintf("column_%v", id)
	boardKey := fmt.Sprintf("leaderboard")
	userKey := fmt.Sprintf("user_column_agree_%s", userUuid)
	exists := make([]int64, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, hotKey)
		pipe.Exists(ctx, statisticKey)
		pipe.Exists(ctx, boardKey)
		pipe.Exists(ctx, userKey)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if cache about user column agree exist: id(%v), uuid(%s), userUuid(%s) ", id, uuid, userUuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, exist)
	}

	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if exists[0] == 1 {
			pipe.ZIncrBy(ctx, hotKey, 1, fmt.Sprintf("%v%s%s", id, "%", uuid))
		}
		if exists[1] == 1 {
			pipe.HIncrBy(ctx, statisticKey, "agree", 1)
		}
		if exists[2] == 1 {
			pipe.ZIncrBy(ctx, boardKey, 1, fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%column"))
		}
		if exists[3] == 1 {
			pipe.SAdd(ctx, userKey, id)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user column agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid))
	}
	return nil
}

func (r *columnRepo) SendColumnStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error {
	achievement := map[string]interface{}{}
	achievement["uuid"] = uuid
	achievement["userUuid"] = userUuid
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

func (r *columnRepo) SetColumnCollectToCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("column_%v", id)
	collectKey := fmt.Sprintf("collections_%v_column", collectionsId)
	collectionsKey := fmt.Sprintf("collections_%v", collectionsId)
	creationKey := fmt.Sprintf("creation_user_%s", userUuid)
	userKey := fmt.Sprintf("user_column_collect_%s", userUuid)
	exists := make([]int64, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, statisticKey)
		pipe.Exists(ctx, collectKey)
		pipe.Exists(ctx, collectionsKey)
		pipe.Exists(ctx, creationKey)
		pipe.Exists(ctx, userKey)
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if cache about user column collect exist: id(%v), collectionsId(%v), uuid(%s), userUuid(%s) ", id, collectionsId, uuid, userUuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, exist)
	}

	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		if exists[0] == 1 {
			pipe.HIncrBy(ctx, statisticKey, "collect", 1)
		}
		if exists[1] == 1 {
			pipe.ZAddNX(ctx, collectKey, &redis.Z{
				Score:  0,
				Member: fmt.Sprintf("%v%s%s", id, "%", uuid),
			})
		}
		if exists[2] == 1 {
			pipe.HIncrBy(ctx, collectionsKey, "column", 1)
		}
		if exists[3] == 1 {
			pipe.HIncrBy(ctx, creationKey, "collect", 1)
		}
		if exists[4] == 1 {
			pipe.SAdd(ctx, userKey, id)
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add column collect to cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
	}
	return nil
}

func (r *columnRepo) SetUserColumnAgreeToCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`)
	keys := []string{"user_column_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user column agree to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) SetUserColumnCollectToCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`)
	keys := []string{"user_column_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user column collect to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) SetCollectionsColumnCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
	collect := &Collect{
		CollectionsId: collectionsId,
		Uuid:          userUuid,
		CreationsId:   id,
		Mode:          2,
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

func (r *columnRepo) SetCollectionColumn(ctx context.Context, collectionsId int32, userUuid string) error {
	c := Collections{}
	err := r.data.DB(ctx).Model(&c).Where("id = ? and uuid = ?", collectionsId, userUuid).Update("`column`", gorm.Expr("`column` + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add collection column collect: collectionsId(%v), userUuid(%s)", collectionsId, userUuid))
	}
	return nil
}

func (r *columnRepo) SetUserColumnCollect(ctx context.Context, id int32, userUuid string) error {
	cc := &ColumnCollect{
		ColumnId: id,
		Uuid:     userUuid,
		Status:   1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(cc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user column collect: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) SetCreationUserCollect(ctx context.Context, userUuid string) error {
	cu := &CreationUser{
		Uuid:    userUuid,
		Collect: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"collect": gorm.Expr("collect + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation collect: uuid(%v)", userUuid))
	}
	return nil
}

func (r *columnRepo) CancelColumnAgree(ctx context.Context, id int32, uuid string) error {
	cs := ColumnStatistic{}
	err := r.data.DB(ctx).Model(&cs).Where("column_id = ? and uuid = ? and agree > 0", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column agree: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) CancelColumnAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	hotKey := "column_hot"
	boardKey := "leaderboard"
	statisticKey := fmt.Sprintf("column_%v", id)
	userKey := fmt.Sprintf("user_column_agree_%s", userUuid)

	var script = redis.NewScript(`
					local hotKey = KEYS[1]
                    local member = ARGV[1]
					local hotKeyExist = redis.call("EXISTS", hotKey)
					if hotKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", hotKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", hotKey, -1, member)
						end
					end

					local boardKey = KEYS[2]
                    local member = ARGV[2]
					local boardKeyExist = redis.call("EXISTS", boardKey)
					if boardKeyExist == 1 then
						local score = tonumber(redis.call("ZSCORE", boardKey, member))
						if score > 0 then
  							redis.call("ZINCRBY", boardKey, -1, member)
						end
					end

					local statisticKey = KEYS[3]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "agree", -1)
						end
					end

					local userKey = KEYS[4]
					local commentId = ARGV[3]
					redis.call("SREM", userKey, commentId)
					return 0
	`)
	keys := []string{hotKey, boardKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%column"), id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid))
	}
	return nil
}

func (r *columnRepo) CancelColumnCollect(ctx context.Context, id int32, uuid string) error {
	cs := &ColumnStatistic{}
	err := r.data.DB(ctx).Model(cs).Where("column_id = ? and uuid = ? and collect > 0", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column collect: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) CancelColumnCollectFromCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("column_%v", id)
	collectKey := fmt.Sprintf("collections_%v_column", collectionsId)
	collectionsKey := fmt.Sprintf("collections_%v", collectionsId)
	creationKey := fmt.Sprintf("creation_user_%s", userUuid)
	userKey := fmt.Sprintf("user_column_collect_%s", userUuid)
	var script = redis.NewScript(`
					local statisticKey = KEYS[1]
					local statisticKeyExist = redis.call("EXISTS", statisticKey)
					if statisticKeyExist == 1 then
						local number = tonumber(redis.call("HGET", statisticKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", statisticKey, "collect", -1)
						end
					end

					local collectKey = KEYS[2]
					local columnMember = ARGV[1]
					redis.call("ZREM", collectKey, columnMember)

					local collectionsKey = KEYS[3]
					local collectionsKeyExist = redis.call("EXISTS", collectionsKey)
					if collectionsKeyExist == 1 then
						local number = tonumber(redis.call("HGET", collectionsKey, "column"))
						if number > 0 then
  							redis.call("HINCRBY", collectionsKey, "column", -1)
						end
					end

					local creationKey = KEYS[4]
					local creationKeyExist = redis.call("EXISTS", creationKey)
					if creationKeyExist == 1 then
						local number = tonumber(redis.call("HGET", creationKey, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", creationKey, "collect", -1)
						end
					end

					local userKey = KEYS[5]
					local columnId = ARGV[2]
					redis.call("SREM", userKey, columnId)

					return 0
	`)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column collect from cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
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

func (r *columnRepo) SetUserColumnSubscribeToCache(ctx context.Context, id int32, uuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end
					return 0
	`)
	keys := []string{"user_column_subscribe_" + uuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user column subscribe to cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) SetColumnSubscribeToCache(ctx context.Context, id int32, author, uuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SADD", key, change)
					end

					local key2 = KEYS[2]
					local member = KEYS[3]
					local value = redis.call("EXISTS", key2)
					if value == 1 then
  						redis.call("ZADD", key2, change, member)
					end

					local key3 = KEYS[4]
					local value = redis.call("EXISTS", key3)
					if value == 1 then
  						redis.call("HINCRBY", key3, "subscribe", 1)
					end

					return 0
	`)
	keys := []string{"user_column_subscribe_" + uuid, "user_column_subscribe_list_" + uuid, fmt.Sprintf("%v%s%s", id, "%", author), "creation_user_" + uuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set column subscribe to cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) AddUserCreationSubscribe(ctx context.Context, uuid string) error {
	cu := &CreationUser{
		Uuid:      uuid,
		Subscribe: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"subscribe": gorm.Expr("subscribe + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation subscribe: uuid(%v)", uuid))
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

func (r *columnRepo) CancelUserColumnAgreeFromCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`)
	keys := []string{"user_column_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user column agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) CancelUserColumnAgree(ctx context.Context, id int32, userUuid string) error {
	ca := ColumnAgree{}
	err := r.data.DB(ctx).Model(&ca).Where("column_id = ? and uuid = ?", id, userUuid).Update("status", 2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user column agree: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) CancelUserColumnCollectFromCache(ctx context.Context, id int32, userUuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("SREM", key, change)
					end
					return 0
	`)
	keys := []string{"user_column_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user column agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) CancelCollectionsColumnCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&Collect{}).Where("creations_id = ? and mode = ? and uuid = ?", id, 2, userUuid).Updates(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column collect: column_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) CancelCollectionColumn(ctx context.Context, id int32, uuid string) error {
	collections := &Collections{}
	err := r.data.DB(ctx).Model(collections).Where("id = ? and uuid = ? and `column` > 0", id, uuid).Update("`column`", gorm.Expr("`column` - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel collections column: id(%v)", id))
	}
	return nil
}

func (r *columnRepo) CancelUserColumnSubscribeFromCache(ctx context.Context, id int32, uuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					redis.call("SREM", key, change)

					return 0
	`)
	keys := []string{"user_column_subscribe_" + uuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user column subscribe from cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *columnRepo) CancelColumnSubscribeFromCache(ctx context.Context, id int32, author, uuid string) error {
	var script = redis.NewScript(`
					local key = KEYS[1]
                    local change = ARGV[1]
					redis.call("SREM", key, change)

					local key2 = KEYS[2]
					local member = KEYS[3]
					redis.call("ZREM", key2, member)

					local key3 = KEYS[4]
					local key3Exist = redis.call("EXISTS", key3)
					if key3Exist == 1 then
						local number = tonumber(redis.call("HGET", key3, "subscribe"))
						if number > 0 then
  							redis.call("HINCRBY", key3, "subscribe", -1)
						end
					end
					return 0
	`)
	keys := []string{"user_column_subscribe_" + uuid, "user_column_subscribe_list_" + uuid, fmt.Sprintf("%v%s%s", id, "%", author), "creation_user_" + uuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel column subscribe from cache: id(%v), author(%s), uuid(%s)", id, author, uuid))
	}
	return nil
}

func (r *columnRepo) ReduceCreationUserCollect(ctx context.Context, uuid string) error {
	cu := &CreationUser{}
	err := r.data.DB(ctx).Model(cu).Where("uuid = ? and collect > 0", uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user collect: uuid(%v)", uuid))
	}
	return nil
}

func (r *columnRepo) ReduceUserCreationSubscribe(ctx context.Context, uuid string) error {
	cu := &CreationUser{}
	err := r.data.DB(ctx).Model(cu).Where("uuid = ? and subscribe > 0", uuid).Update("subscribe", gorm.Expr("subscribe - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce user creation subscribe: uuid(%v)", uuid))
	}
	return nil
}

func (r *columnRepo) CancelUserColumnCollect(ctx context.Context, id int32, userUuid string) error {
	cc := &ColumnCollect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&ColumnCollect{}).Where("column_id = ? and uuid = ?", id, userUuid).Updates(cc).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user column collect: column_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *columnRepo) GetCollectionsIdFromColumnCollect(ctx context.Context, id int32) (int32, error) {
	collect := &Collect{}
	err := r.data.db.WithContext(ctx).Where("creations_id = ? and mode = ?", id, 2).First(collect).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections id  from db: creationsId(%v)", id))
	}
	return collect.CollectionsId, nil
}

func (r *columnRepo) GetUserColumnAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userColumnAgreeExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserColumnAgreeFromCache(ctx, uuid)
	} else {
		return r.getUserColumnAgreeFromDb(ctx, uuid)
	}
}

func (r *columnRepo) userColumnAgreeExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_column_agree_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user column agree exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *columnRepo) getUserColumnAgreeFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_column_agree_" + uuid
	agreeSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column agree from cache: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range agreeSet {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s), uuid(%s), key(%s)", id, uuid, key))
		}
		agreeMap[int32(id)] = true
	}
	return agreeMap, nil
}

func (r *columnRepo) getUserColumnAgreeFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*ColumnAgree, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column agree from db: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range list {
		agreeMap[item.ColumnId] = true
	}
	if len(list) != 0 {
		go r.setUserColumnAgreeToCache(uuid, list)
	}
	return agreeMap, nil
}

func (r *columnRepo) setUserColumnAgreeToCache(uuid string, agreeList []*ColumnAgree) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0)
		key := "user_column_agree_" + uuid
		for _, item := range agreeList {
			set = append(set, item.ColumnId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user column agree to cache: uuid(%s), agreeList(%v), error(%s) ", uuid, agreeList, err.Error())
	}
}

func (r *columnRepo) GetUserColumnCollect(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userColumnCollectExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserColumnCollectFromCache(ctx, uuid)
	} else {
		return r.getUserColumnCollectFromDb(ctx, uuid)
	}
}

func (r *columnRepo) userColumnCollectExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_column_collect_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user column collect exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *columnRepo) getUserColumnCollectFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_column_collect_" + uuid
	collectSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column collect from cache: uuid(%s)", uuid))
	}

	collectMap := make(map[int32]bool, 0)
	for _, item := range collectSet {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s), uuid(%s), key(%s)", id, uuid, key))
		}
		collectMap[int32(id)] = true
	}
	return collectMap, nil
}

func (r *columnRepo) getUserColumnCollectFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*ColumnCollect, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column collect from db: uuid(%s)", uuid))
	}

	collectMap := make(map[int32]bool, 0)
	for _, item := range list {
		collectMap[item.ColumnId] = true
	}
	if len(list) != 0 {
		go r.setUserColumnCollectToCache(uuid, list)
	}
	return collectMap, nil
}

func (r *columnRepo) setUserColumnCollectToCache(uuid string, collectList []*ColumnCollect) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0)
		key := "user_column_collect_" + uuid
		for _, item := range collectList {
			set = append(set, item.ColumnId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user column collect to cache: uuid(%s), collectList(%v), error(%s) ", uuid, collectList, err.Error())
	}
}

func (r *columnRepo) GetUserSubscribeColumn(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userSubscribeColumnExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserColumnSubscribeFromCache(ctx, uuid)
	} else {
		return r.getUserColumnSubscribeFromDb(ctx, uuid)
	}
}

func (r *columnRepo) userSubscribeColumnExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_column_subscribe_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user column subscribe exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *columnRepo) getUserColumnSubscribeFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_column_subscribe_" + uuid
	subscribeSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column subscribe from cache: uuid(%s)", uuid))
	}
	subscribeMap := make(map[int32]bool, 0)
	for _, item := range subscribeSet {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s), uuid(%s), key(%s)", id, uuid, key))
		}
		subscribeMap[int32(id)] = true
	}
	return subscribeMap, nil
}

func (r *columnRepo) getUserColumnSubscribeFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*Subscribe, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user column subscribe from db: uuid(%s)", uuid))
	}

	subscribeMap := make(map[int32]bool, 0)
	for _, item := range list {
		subscribeMap[item.ColumnId] = true
	}
	if len(list) != 0 {
		go r.setUserColumnSubscribeToCache(uuid, list)
	}
	return subscribeMap, nil
}

func (r *columnRepo) setUserColumnSubscribeToCache(uuid string, subscribeList []*Subscribe) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0)
		key := "user_column_subscribe_" + uuid
		for _, item := range subscribeList {
			set = append(set, item.ColumnId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user column subscribe to cache: uuid(%s), subscribeList(%v), error(%s) ", uuid, subscribeList, err.Error())
	}
}

func (r *columnRepo) GetAuthorFromSubscribe(ctx context.Context, id int32) (string, error) {
	c := &Column{}
	err := r.data.db.WithContext(ctx).Where("column_id = ?", id).First(c).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get column author rom db: id(%v)", id))
	}
	return c.Uuid, nil
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
