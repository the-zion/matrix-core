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
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"time"
)

var _ biz.CreationRepo = (*creationRepo)(nil)

type creationRepo struct {
	data *Data
	log  *log.Helper
}

func NewCreationRepo(data *Data, logger log.Logger) biz.CreationRepo {
	return &creationRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/creation")),
	}
}

func (r *creationRepo) GetLeaderBoard(ctx context.Context) ([]*biz.LeaderBoard, error) {
	return r.getLeaderBoardFromCache(ctx)
}

func (r *creationRepo) getLeaderBoardFromCache(ctx context.Context) ([]*biz.LeaderBoard, error) {
	list, err := r.data.redisCli.ZRevRange(ctx, "leaderboard", 0, 9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get leader board from cache")
	}

	board := make([]*biz.LeaderBoard, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		board = append(board, &biz.LeaderBoard{
			Id:   int32(id),
			Uuid: member[1],
			Mode: member[2],
		})
	}
	return board, nil
}

func (r *creationRepo) GetLastCollectionsDraft(ctx context.Context, uuid string) (*biz.CollectionsDraft, error) {
	draft := &CollectionsDraft{}
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("collections draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
	}
	return &biz.CollectionsDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *creationRepo) GetCollectionsContentReview(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	key := "collections_content_irregular_" + uuid
	review, err := r.getCollectionsContentReviewFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getCollectionsContentReviewFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setCollectionsContentReviewToCache(key, review)
		})()
	}
	return review, nil
}

func (r *creationRepo) getCollectionsContentReviewFromCache(ctx context.Context, page int32, key string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get creationRepo content irregular list from cache: key(%s), page(%v)", key, page))
	}

	review := make([]*biz.TextReview, 0)
	for _index, item := range list {
		var textReview = &biz.TextReview{}
		err = json.Unmarshal([]byte(item), textReview)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: contentReview(%v)", item))
		}
		review = append(review, &biz.TextReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: textReview.CreationId,
			Title:      textReview.Title,
			Kind:       textReview.Kind,
			Uuid:       textReview.Uuid,
			CreateAt:   textReview.CreateAt,
			JobId:      textReview.JobId,
			Label:      textReview.Label,
			Result:     textReview.Result,
			Section:    textReview.Section,
		})
	}
	return review, nil
}

func (r *creationRepo) getCollectionsContentReviewFromDB(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*CollectionsContentReview, 0)
	err := r.data.db.WithContext(ctx).Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections content review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.TextReview, 0)
	for _index, item := range list {
		review = append(review, &biz.TextReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: item.CollectionsId,
			Kind:       item.Kind,
			Title:      item.Title,
			Uuid:       item.Uuid,
			JobId:      item.JobId,
			CreateAt:   item.CreatedAt.Format("2006-01-02"),
			Label:      item.Label,
			Result:     item.Result,
			Section:    item.Section,
		})
	}
	return review, nil
}

func (r *creationRepo) setCollectionsContentReviewToCache(key string, review []*biz.TextReview) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0)
		for _, item := range review {
			m, err := json.Marshal(item)
			if err != nil {
				r.log.Errorf("fail to marshal avatar review: contentReview(%v)", review)
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set collections content review to cache: contentReview(%v)", review)
	}
}

func (r *creationRepo) GetCollectArticleList(ctx context.Context, id, page int32) ([]*biz.Article, error) {
	key := fmt.Sprintf("collections_%v_article", id)
	articleList, err := r.getCollectArticleListFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(articleList)
	if size != 0 {
		return articleList, nil
	}

	articleList, err = r.getCollectArticleListFromDB(ctx, id, page)
	if err != nil {
		return nil, err
	}

	size = len(articleList)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setCollectArticleListToCache(key, articleList)
		})()
	}
	return articleList, nil
}

func (r *creationRepo) getCollectArticleListFromCache(ctx context.Context, page int32, key string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collect article list from cache: key(%s), page(%v)", key, page))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		article = append(article, &biz.Article{
			ArticleId: int32(id),
			Uuid:      member[1],
		})
	}
	return article, nil
}

func (r *creationRepo) getCollectArticleListFromDB(ctx context.Context, id, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collect, 0)
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ? and status = ?", id, 1, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collect article list from db: id(%v), page(%v)", id, page))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.CreationsId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *creationRepo) setCollectArticleListToCache(key string, article []*biz.Article) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.ArticleId),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user collect article to cache: article(%v)", article)
	}
}

func (r *creationRepo) GetCollectArticleCount(ctx context.Context, id int32) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ? and status = ?", id, 1, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect article count from db: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectTalkList(ctx context.Context, id, page int32) ([]*biz.Talk, error) {
	key := fmt.Sprintf("collections_%v_talk", id)
	talkList, err := r.getCollectTalkListFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(talkList)
	if size != 0 {
		return talkList, nil
	}

	talkList, err = r.getCollectTalkListFromDB(ctx, id, page)
	if err != nil {
		return nil, err
	}

	size = len(talkList)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setCollectTalkListToCache(key, talkList)
		})()
	}
	return talkList, nil
}

func (r *creationRepo) getCollectTalkListFromCache(ctx context.Context, page int32, key string) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collect talk list from cache: key(%s), page(%v)", key, page))
	}

	talk := make([]*biz.Talk, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		talk = append(talk, &biz.Talk{
			TalkId: int32(id),
			Uuid:   member[1],
		})
	}
	return talk, nil
}

func (r *creationRepo) getCollectTalkListFromDB(ctx context.Context, id, page int32) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collect, 0)
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ? and status = ?", id, 3, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collect talk from db: id(%v), page(%v)", id, page))
	}

	talk := make([]*biz.Talk, 0)
	for _, item := range list {
		talk = append(talk, &biz.Talk{
			TalkId: item.CreationsId,
			Uuid:   item.Uuid,
		})
	}
	return talk, nil
}

func (r *creationRepo) setCollectTalkListToCache(key string, talk []*biz.Talk) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user collect talk to cache: article(%v)", talk)
	}
}

func (r *creationRepo) GetCollectTalkCount(ctx context.Context, id int32) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ? and status = ?", id, 2, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect talk count from db: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectColumnList(ctx context.Context, id, page int32) ([]*biz.Column, error) {
	key := fmt.Sprintf("collections_%v_column", id)
	columnList, err := r.getCollectColumnListFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(columnList)
	if size != 0 {
		return columnList, nil
	}

	columnList, err = r.getCollectColumnListFromDB(ctx, id, page)
	if err != nil {
		return nil, err
	}

	size = len(columnList)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setCollectColumnListToCache(key, columnList)
		})()
	}
	return columnList, nil
}

func (r *creationRepo) getCollectColumnListFromCache(ctx context.Context, page int32, key string) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, key, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collect column list from cache: key(%s), page(%v)", key, page))
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

func (r *creationRepo) getCollectColumnListFromDB(ctx context.Context, id, page int32) ([]*biz.Column, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collect, 0)
	err := r.data.db.WithContext(ctx).Where("collections_id = ? and mode = ? and status = ?", id, 2, 1).Order("updated_at desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collect column list from db: id(%v), page(%v)", id, page))
	}

	column := make([]*biz.Column, 0)
	for _, item := range list {
		column = append(column, &biz.Column{
			ColumnId: item.CreationsId,
			Uuid:     item.Uuid,
		})
	}
	return column, nil
}

func (r *creationRepo) setCollectColumnListToCache(key string, column []*biz.Column) {
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
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user collect column to cache: column(%v)", column)
	}
}

func (r *creationRepo) GetCollectColumnCount(ctx context.Context, id int32) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and mode = ? and status = ?", id, 3, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect column count from db: id(%v)", id))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollections(ctx context.Context, id int32, uuid string) (*biz.Collections, error) {
	var collections *biz.Collections
	key := "collections_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		collections, err = r.getCollectionsFromCache(ctx, key, uuid)
		if err != nil {
			return nil, err
		}
		return collections, nil
	}

	collections, err = r.getCollectionsFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	go r.data.Recover(context.Background(), func(ctx context.Context) {
		r.setCollectionsToCache(key, collections)
	})()

	if collections.Auth == 2 && collections.Uuid != uuid {
		return nil, errors.Errorf("fail to get collection: no auth")
	}

	return collections, nil
}

func (r *creationRepo) getCollectionsFromCache(ctx context.Context, key, uuid string) (*biz.Collections, error) {
	collections, err := r.data.redisCli.HMGet(ctx, key, "uuid", "auth", "article", "column", "talk").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0}
	for _index, count := range collections[1:] {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	if val[0] == 2 && collections[0].(string) != uuid {
		return nil, errors.Errorf("fail to get collection from cache: no auth")
	}
	return &biz.Collections{
		Uuid:    collections[0].(string),
		Auth:    val[0],
		Article: val[1],
		Column:  val[2],
		Talk:    val[3],
	}, nil
}

func (r *creationRepo) getCollectionsFromDB(ctx context.Context, id int32) (*biz.Collections, error) {
	collections := &Collections{}
	err := r.data.db.WithContext(ctx).Where("collections_id = ?", id).First(collections).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get collections from db: id(%v)", id))
	}
	return &biz.Collections{
		Uuid:    collections.Uuid,
		Article: collections.Article,
		Column:  collections.Column,
		Talk:    collections.Talk,
		Auth:    collections.Auth,
	}, nil
}

func (r *creationRepo) setCollectionsToCache(key string, collections *biz.Collections) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "uuid", collections.Uuid, "auth", collections.Auth, "article", collections.Article, "column", collections.Column, "talk", collections.Talk)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set collections to cache, err(%s)", err.Error())
	}
}

func (r *creationRepo) GetCollectionListInfo(ctx context.Context, ids []int32) ([]*biz.Collections, error) {
	collectionsListInfo := make([]*biz.Collections, 0)
	exists, unExists, err := r.collectionsListInfoExist(ctx, ids)
	if err != nil {
		return nil, err
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getCollectionsListInfoFromCache(ctx, exists, &collectionsListInfo)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getCollectionsListInfoFromDb(ctx, unExists, &collectionsListInfo)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return collectionsListInfo, nil
}

func (r *creationRepo) collectionsListInfoExist(ctx context.Context, ids []int32) ([]int32, []int32, error) {
	exists := make([]int32, 0)
	unExists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "collection_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if collection exist from cache: ids(%v)", ids))
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

func (r *creationRepo) getCollectionsListInfoFromCache(ctx context.Context, exists []int32, collectionsListInfo *[]*biz.Collections) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range exists {
			pipe.HMGet(ctx, "collection_"+strconv.Itoa(int(id)), "name", "introduce")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get collections list info from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		val := []string{"", ""}
		for _index, info := range item.(*redis.SliceCmd).Val() {
			if info == nil {
				break
			}
			val[_index] = info.(string)
		}
		*collectionsListInfo = append(*collectionsListInfo, &biz.Collections{
			CollectionsId: exists[index],
		})
	}
	return nil
}

func (r *creationRepo) getCollectionsListInfoFromDb(ctx context.Context, unExists []int32, collectionsListInfo *[]*biz.Collections) error {
	list := make([]*Collections, 0)
	err := r.data.db.WithContext(ctx).Where("id IN ?", unExists).Order("id desc").Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get collections list info from db: ids(%v)", unExists))
	}

	for _, item := range list {
		*collectionsListInfo = append(*collectionsListInfo, &biz.Collections{
			CollectionsId: item.CollectionsId,
			Uuid:          item.Uuid,
			Auth:          item.Auth,
		})
	}

	if len(list) != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setCollectionsListInfoToCache(list)
		})()
	}
	return nil
}

func (r *creationRepo) setCollectionsListInfoToCache(collectionsList []*Collections) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range collectionsList {
			key := "collection_" + strconv.Itoa(int(item.CollectionsId))
			pipe.HSetNX(ctx, key, "uuid", item.Uuid)
			pipe.HSetNX(ctx, key, "auth", item.Auth)
			pipe.Expire(ctx, key, time.Hour*8)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set collections info to cache, err(%s)", err.Error())
	}
}

func (r *creationRepo) GetCollectCount(ctx context.Context, id int32) (int64, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collect{}).Where("collections_id = ? and status = ?", id, 1).Limit(1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collect count from db: collections_id(%v)", id))
	}
	return count, nil
}

func (r *creationRepo) GetCollectionsList(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	collections, err := r.getUserCollectionsListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(collections)
	if size != 0 {
		return collections, nil
	}

	collections, err = r.getUserCollectionsListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(collections)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserCollectionsListToCache("user_collections_list_"+uuid, collections)
		})()
	}
	return collections, nil
}

func (r *creationRepo) getUserCollectionsListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_collections_list_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collections list from cache: key(%s), page(%v)", "user_collections_list_"+uuid, page))
	}

	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", item))
		}
		collections = append(collections, &biz.Collections{
			CollectionsId: int32(id),
		})
	}
	return collections, nil
}

func (r *creationRepo) getUserCollectionsListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collections, 0)
	handle := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("collections_id desc").Offset(index * 10).Limit(10).Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections from db: uuid(%s)", uuid))
	}
	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		collections = append(collections, &biz.Collections{
			CollectionsId: item.CollectionsId,
		})
	}
	return collections, nil
}

func (r *creationRepo) GetCollectionsListByVisitor(ctx context.Context, uuid string, page int32) ([]*biz.Collections, error) {
	collections, err := r.getUserCollectionsListVisitorFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(collections)
	if size != 0 {
		return collections, nil
	}

	collections, err = r.getUserCollectionsListVisitorFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(collections)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserCollectionsListToCache("user_collections_list_visitor_"+uuid, collections)
		})()
	}
	return collections, nil
}

func (r *creationRepo) getUserCollectionsListVisitorFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_collections_list_visitor_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collections list visitor from cache: key(%s), page(%v)", "user_collections_list_visitor_"+uuid, page))
	}

	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", item))
		}
		collections = append(collections, &biz.Collections{
			CollectionsId: int32(id),
		})
	}
	return collections, nil
}

func (r *creationRepo) getUserCollectionsListVisitorFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Collections, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Collections, 0)
	handle := r.data.db.WithContext(ctx).Where("uuid = ? and auth = ?", uuid, 1).Order("collections_id desc").Offset(index * 10).Limit(10).Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections visitor from db: uuid(%s)", uuid))
	}
	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		collections = append(collections, &biz.Collections{
			CollectionsId: item.CollectionsId,
		})
	}
	return collections, nil
}

func (r *creationRepo) setUserCollectionsListToCache(key string, collections []*biz.Collections) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range collections {
			z = append(z, &redis.Z{
				Score:  float64(item.CollectionsId),
				Member: strconv.Itoa(int(item.CollectionsId)),
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set collections to cache: collections(%v)", collections)
	}
}

func (r *creationRepo) GetCollectionsListAll(ctx context.Context, uuid string) ([]*biz.Collections, error) {
	collections, err := r.getUserCollectionsListAllFromCache(ctx, uuid)
	if err != nil {
		return nil, err
	}

	size := len(collections)
	if size != 0 {
		return collections, nil
	}

	collections, err = r.getUserCollectionsListAllFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	size = len(collections)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserCollectionsListToCache("user_collections_list_all_"+uuid, collections)
		})()
	}
	return collections, nil
}

func (r *creationRepo) getUserCollectionsListAllFromCache(ctx context.Context, uuid string) ([]*biz.Collections, error) {
	list, err := r.data.redisCli.ZRevRange(ctx, "user_collections_list_all_"+uuid, 0, -1).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user collections list all from cache: key(%s)", "user_collections_list_all_"+uuid))
	}

	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		id, err := strconv.ParseInt(item, 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", item))
		}
		collections = append(collections, &biz.Collections{
			CollectionsId: int32(id),
		})
	}
	return collections, nil
}

func (r *creationRepo) getUserCollectionsListAllFromDB(ctx context.Context, uuid string) ([]*biz.Collections, error) {
	list := make([]*Collections, 0)
	handle := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("collections_id desc").Find(&list)
	err := handle.Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get collections from db: uuid(%s)", uuid))
	}
	collections := make([]*biz.Collections, 0)
	for _, item := range list {
		collections = append(collections, &biz.Collections{
			CollectionsId: item.CollectionsId,
		})
	}
	return collections, nil
}

func (r *creationRepo) GetCollectionsCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collections{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCollectionsVisitorCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Collections{}).Where("uuid = ? and auth = ?", uuid, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections visitor count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *creationRepo) GetCreationUser(ctx context.Context, uuid string) (*biz.CreationUser, error) {
	var creationUser *biz.CreationUser
	key := "creation_user_" + uuid
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		creationUser, err = r.getCreationUserFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return creationUser, nil
	}

	creationUser, err = r.getCreationUserFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	go r.data.Recover(context.Background(), func(ctx context.Context) {
		r.setCreationUserToCache(key, creationUser)
	})()

	return creationUser, nil
}

func (r *creationRepo) getCreationUserFromCache(ctx context.Context, key string) (*biz.CreationUser, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "article", "column", "talk", "collections", "collect", "subscribe").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user creation form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0, 0}
	for _index, count := range statistic {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.CreationUser{
		Article:     val[0],
		Column:      val[1],
		Talk:        val[2],
		Collections: val[3],
		Collect:     val[4],
		Subscribe:   val[5],
	}, nil
}

func (r *creationRepo) getCreationUserFromDB(ctx context.Context, uuid string) (*biz.CreationUser, error) {
	cuv := &CreationUser{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(cuv).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get user creation from db: uuid(%v)", uuid))
	}
	return &biz.CreationUser{
		Article:     cuv.Article,
		Column:      cuv.Column,
		Collections: cuv.Collections,
		Talk:        cuv.Talk,
		Collect:     cuv.Collect,
		Subscribe:   cuv.Subscribe,
	}, nil
}

func (r *creationRepo) setCreationUserToCache(key string, creationUser *biz.CreationUser) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "article", creationUser.Article, "talk", creationUser.Talk, "column", creationUser.Column, "collections", creationUser.Collections, "collect", creationUser.Collect, "subscribe", creationUser.Subscribe)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user creation to cache, err(%s)", err.Error())
	}
}

func (r *creationRepo) GetCreationUserVisitor(ctx context.Context, uuid string) (*biz.CreationUser, error) {
	var creationUser *biz.CreationUser
	key := "creation_user_visitor_" + uuid
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		creationUser, err = r.getCreationUserVisitorFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return creationUser, nil
	}

	creationUser, err = r.getCreationUserVisitorFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	go r.data.Recover(context.Background(), func(ctx context.Context) {
		r.setCreationUserVisitorToCache(key, creationUser)
	})()

	return creationUser, nil
}

func (r *creationRepo) GetUserTimeLineList(ctx context.Context, page int32, uuid string) ([]*biz.TimeLine, error) {
	timeline, err := r.getUserTimeLineListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(timeline)
	if size != 0 {
		return timeline, nil
	}

	timeline, err = r.getUserTimeLineListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(timeline)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserTimeLineListToCache("user_timeline_list_"+uuid, timeline)
		})()
	}
	return timeline, nil
}

func (r *creationRepo) getUserTimeLineListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.TimeLine, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_timeline_list_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user timeline list from cache: key(%s), page(%v)", "user_timeline_list_"+uuid, page))
	}

	timeline := make([]*biz.TimeLine, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		creationId, err := strconv.ParseInt(member[1], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: creationId(%s)", member[1]))
		}
		mode, err := strconv.ParseInt(member[2], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: mode(%s)", member[2]))
		}
		timeline = append(timeline, &biz.TimeLine{
			Id:         int32(id),
			CreationId: int32(creationId),
			Mode:       int32(mode),
			Uuid:       member[3],
		})
	}
	return timeline, nil
}

func (r *creationRepo) getUserTimeLineListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.TimeLine, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*TimeLine, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user timeline list from db: page(%v), uuid(%s)", page, uuid))
	}

	timeline := make([]*biz.TimeLine, 0)
	for _, item := range list {
		timeline = append(timeline, &biz.TimeLine{
			Id:         int32(item.ID),
			Uuid:       item.Uuid,
			CreationId: item.CreationsId,
			Mode:       item.Mode,
		})
	}
	return timeline, nil
}

func (r *creationRepo) setUserTimeLineListToCache(key string, timeline []*biz.TimeLine) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range timeline {
			z = append(z, &redis.Z{
				Score:  float64(item.Id),
				Member: strconv.Itoa(int(item.Id)) + "%" + strconv.Itoa(int(item.CreationId)) + "%" + strconv.Itoa(int(item.Mode)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set timeline to cache: timeline(%v)", timeline)
	}
}

func (r *creationRepo) GetCollectionsAuth(ctx context.Context, id int32) (int32, error) {
	collections := &Collections{}
	err := r.data.db.WithContext(ctx).Where("collections_id = ?", id).First(collections).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections auth from db: id(%v)", id))
	}
	return collections.Auth, nil
}

func (r *creationRepo) getCreationUserVisitorFromCache(ctx context.Context, key string) (*biz.CreationUser, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "article", "column", "talk", "collections").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user creation form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0}
	for _index, count := range statistic {
		if count == nil {
			break
		}
		num, err := strconv.ParseInt(count.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: count(%v)", count))
		}
		val[_index] = int32(num)
	}
	return &biz.CreationUser{
		Article:     val[0],
		Column:      val[1],
		Talk:        val[2],
		Collections: val[3],
	}, nil
}

func (r *creationRepo) getCreationUserVisitorFromDB(ctx context.Context, uuid string) (*biz.CreationUser, error) {
	cuv := &CreationUserVisitor{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).First(cuv).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get user creation from db: uuid(%v)", uuid))
	}
	return &biz.CreationUser{
		Article:     cuv.Article,
		Column:      cuv.Column,
		Collections: cuv.Collections,
		Talk:        cuv.Talk,
	}, nil
}

func (r *creationRepo) setCreationUserVisitorToCache(key string, creationUser *biz.CreationUser) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "article", creationUser.Article, "talk", creationUser.Talk, "column", creationUser.Column, "collections", creationUser.Collections)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user creation to cache, err(%s)", err.Error())
	}
}

func (r *creationRepo) SendCollections(ctx context.Context, id int32, uuid string) (*biz.CollectionsDraft, error) {
	cd := &CollectionsDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&CollectionsDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 1).Updates(cd).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.CollectionsDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *creationRepo) SendCollectionsContentIrregularToMq(ctx context.Context, review *biz.TextReview) error {
	m := make(map[string]interface{}, 0)
	m["creation_id"] = review.CreationId
	m["result"] = review.Result
	m["uuid"] = review.Uuid
	m["job_id"] = review.JobId
	m["label"] = review.Label
	m["title"] = review.Title
	m["kind"] = review.Kind
	m["section"] = review.Section
	m["mode"] = review.Mode
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "collections",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.collectionsMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send collections content review to mq: %v", err))
	}
	return nil
}

func (r *creationRepo) SetCollectionsContentIrregular(ctx context.Context, review *biz.TextReview) (*biz.TextReview, error) {
	ar := &CollectionsContentReview{
		CollectionsId: review.CreationId,
		Title:         review.Title,
		Kind:          review.Kind,
		Uuid:          review.Uuid,
		JobId:         review.JobId,
		Label:         review.Label,
		Result:        review.Result,
		Section:       review.Section,
	}
	err := r.data.DB(ctx).Select("CollectionsId", "Title", "Kind", "Uuid", "JobId", "Label", "Result", "Section").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to add collections content review record: review(%v)", review))
	}
	review.Id = int32(ar.ID)
	review.CreateAt = time.Now().Format("2006-01-02")
	return review, nil
}

func (r *creationRepo) SetCollectionsContentIrregularToCache(ctx context.Context, review *biz.TextReview) error {
	marshal, err := json.Marshal(review)
	if err != nil {
		r.log.Errorf("fail to set collections content irregular to json: json.Marshal(%v), error(%v)", review, err)
	}
	var script = redis.NewScript(`
					local key = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					return 0
	`)
	keys := []string{"collections_content_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set collections content irregular to cache: review(%v)", review))
	}
	return nil
}

func (r *creationRepo) CreateCollectionsDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &CollectionsDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create a collections draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *creationRepo) CreateCollectionsFolder(ctx context.Context, id int32, uuid string) error {
	name := "collections/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a collections folder: id(%v)", id))
	}
	return nil
}

func (r *creationRepo) CreateCollections(ctx context.Context, id, auth int32, uuid string) error {
	collections := &Collections{
		CollectionsId: id,
		Uuid:          uuid,
		Auth:          auth,
	}
	err := r.data.DB(ctx).Select("CollectionsId", "Uuid", "Auth").Create(collections).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a collections: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *creationRepo) CreateCollectionsCache(ctx context.Context, id, auth int32, uuid, mode string) error {
	ids := strconv.Itoa(int(id))
	userCollectionsList := "user_collections_list_" + uuid
	userCollectionsListVisitor := "user_collections_list_visitor_" + uuid
	creationUser := "creation_user_" + uuid
	creationUserVisitor := "creation_user_visitor_" + uuid
	userCollectionsListAll := "user_collections_list_all_" + uuid
	collectionsStatistic := "collections_" + ids
	var script = redis.NewScript(`
					local collectionsStatistic = KEYS[1]
					local userCollectionsList = KEYS[2]
					local userCollectionsListVisitor = KEYS[3]
					local creationUser = KEYS[4]
					local creationUserVisitor = KEYS[5]
					local userCollectionsListAll = KEYS[6]

					local uuid = ARGV[1]
					local auth = ARGV[2]
					local id = ARGV[3]
					local ids = ARGV[4]
					local mode = ARGV[5]

					local userCollectionsListExist = redis.call("EXISTS", userCollectionsList)
					local userCollectionsListVisitorExist = redis.call("EXISTS", userCollectionsListVisitor)
					local creationUserExist = redis.call("EXISTS", creationUser)
					local creationUserVisitorExist = redis.call("EXISTS", creationUserVisitor)
					local userCollectionsListAllExist = redis.call("EXISTS", userCollectionsListAll)

					redis.call("HSETNX", collectionsStatistic, "uuid", uuid)
					redis.call("HSETNX", collectionsStatistic, "auth", auth)
					redis.call("HSETNX", collectionsStatistic, "article", 0)
					redis.call("HSETNX", collectionsStatistic, "column", 0)
					redis.call("HSETNX", collectionsStatistic, "talk", 0)

					if userCollectionsListExist == 1 then
						redis.call("ZADD", userCollectionsList, id, ids)
					end

					if userCollectionsListAllExist == 1 then
						redis.call("ZADD", userCollectionsListAll, id, ids)
					end

					if (creationUserExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUser, "collections", 1)
					end

					if auth == "2" then
						return 0
					end

					if userCollectionsListVisitorExist == 1 then
						redis.call("ZADD", userCollectionsListVisitor, id, ids)
					end

					if (creationUserVisitorExist == 1) and (mode == "create") then
						redis.call("HINCRBY", creationUserVisitor, "collections", 1)
					end
					return 0
	`)
	keys := []string{collectionsStatistic, userCollectionsList, userCollectionsListVisitor, creationUser, creationUserVisitor, userCollectionsListAll}
	values := []interface{}{uuid, auth, id, ids, mode}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create(update) collections cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *creationRepo) AddCreationUserCollections(ctx context.Context, uuid string, auth int32) error {
	cu := &CreationUser{
		Uuid:        uuid,
		Collections: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"collections": gorm.Expr("collections + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation collections: uuid(%v)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := &CreationUserVisitor{
		Uuid:        uuid,
		Collections: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"collections": gorm.Expr("collections + ?", 1)}),
	}).Create(cuv).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation collections visitor: uuid(%v)", uuid))
	}

	return nil
}

func (r *creationRepo) EditCollectionsCos(ctx context.Context, id int32, uuid string) error {
	err := r.EditCollectionsCosContent(ctx, id, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *creationRepo) EditCollectionsCosContent(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "collections/" + uuid + "/" + ids + "/content-edit"
	dest := "collections/" + uuid + "/" + ids + "/content"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy collections from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *creationRepo) UpdateCollectionsCache(ctx context.Context, id, auth int32, uuid, mode string) error {
	return r.CreateCollectionsCache(ctx, id, auth, uuid, mode)
}

func (r *creationRepo) DeleteCollections(ctx context.Context, id int32, uuid string) error {
	collections := &Collections{}
	err := r.data.DB(ctx).Where("collections_id = ? and uuid = ?", id, uuid).Delete(collections).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an collection from db: id(%v), uuid(%v)", id, uuid))
	}
	return nil
}

func (r *creationRepo) DeleteCollect(ctx context.Context, id int32) error {
	collect := &Collect{}
	err := r.data.DB(ctx).Where("collections_id = ?", id).Delete(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an collection from db: id(%v)", id))
	}
	return nil
}

func (r *creationRepo) DeleteCollectionsDraft(ctx context.Context, id int32, uuid string) error {
	cd := &CollectionsDraft{}
	cd.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(cd).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete collections draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *creationRepo) DeleteCreationCache(ctx context.Context, id, auth int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	var script = redis.NewScript(`
					local key1 = KEYS[1]
					local key2 = KEYS[2]
					local key3 = KEYS[3]
					local key4 = KEYS[4]
					local key5 = KEYS[5]
					local key6 = KEYS[6]

					local member = ARGV[1]
					local auth = ARGV[2]

                    redis.call("ZREM", key1, member)
					redis.call("ZREM", key5, member)
					redis.call("DEL", key6)

                    local exist = redis.call("EXISTS", key3)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key3, "collections"))
						if number > 0 then
  							redis.call("HINCRBY", key3, "collections", -1)
						end
					end

                    if auth == "2" then
						return 0
					end

					redis.call("ZREM", key2, member)

					local exist = redis.call("EXISTS", key4)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", key4, "collections"))
						if number > 0 then
  							redis.call("HINCRBY", key4, "collections", -1)
						end
					end
					return 0
	`)
	keys := []string{"user_collections_list_" + uuid, "user_collections_list_visitor_" + uuid, "creation_user_" + uuid, "creation_user_visitor_" + uuid, "user_collections_list_all_" + uuid, "collections_" + ids}
	values := []interface{}{ids, auth}
	_, err := script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete collections cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *creationRepo) ReduceCreationUserCollections(ctx context.Context, auth int32, uuid string) error {
	cu := CreationUser{}
	err := r.data.DB(ctx).Model(&cu).Where("uuid = ? and collections > 0", uuid).Update("collections", gorm.Expr("collections - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user collections: uuid(%s)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := CreationUserVisitor{}
	err = r.data.DB(ctx).Model(&cuv).Where("uuid = ? and collections > 0", uuid).Update("collections", gorm.Expr("collections - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user collections visitor: uuid(%s)", uuid))
	}
	return nil
}

func (r *creationRepo) SetRecord(ctx context.Context, id, mode int32, uuid, operation, ip string) error {
	record := &Record{
		CreationId: id,
		Mode:       mode,
		Uuid:       uuid,
		Operation:  operation,
		Ip:         ip,
	}
	err := r.data.DB(ctx).Create(record).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add an record to db: id(%v), uuid(%s), ip(%s)", id, uuid, ip))
	}
	return nil
}

func (r *creationRepo) SetLeaderBoardToCache(ctx context.Context, boardList []*biz.LeaderBoard) {
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range boardList {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.Id)) + "%" + item.Uuid + "%" + item.Mode,
			})
		}
		pipe.ZAddNX(context.Background(), "leaderboard", z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set LeaderBoard to cache: LeaderBoard(%v)", boardList)
	}
}

func (r *creationRepo) SendReviewToMq(ctx context.Context, review *biz.CollectionsReview) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "collections_review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.collectionsReviewMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *creationRepo) SendCollectionsToMq(ctx context.Context, collections *biz.Collections, mode string) error {
	collectionsMap := map[string]interface{}{}
	collectionsMap["uuid"] = collections.Uuid
	collectionsMap["id"] = collections.CollectionsId
	collectionsMap["auth"] = collections.Auth
	collectionsMap["mode"] = mode

	data, err := json.Marshal(collectionsMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "collections",
		Body:  data,
	}
	msg.WithKeys([]string{collections.Uuid})
	_, err = r.data.collectionsMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send collections to mq: %v", collections))
	}
	return nil
}
