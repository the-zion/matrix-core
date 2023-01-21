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

var _ biz.ArticleRepo = (*articleRepo)(nil)

type articleRepo struct {
	data *Data
	log  *log.Helper
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/article")),
	}
}

func (r *articleRepo) GetLastArticleDraft(ctx context.Context, uuid string) (*biz.ArticleDraft, error) {
	draft := &ArticleDraft{}
	err := r.data.db.WithContext(ctx).Select("id", "status").Where("uuid = ?", uuid).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("article draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get last article draft: uuid(%s)", uuid))
	}
	return &biz.ArticleDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *articleRepo) GetArticle(ctx context.Context, id int32) (*biz.Article, error) {
	article := &Article{}
	err := r.data.db.WithContext(ctx).Select("uuid").Where("article_id = ?", id).First(article).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from db: id(%v)", id))
	}
	return &biz.Article{
		ArticleId: id,
		Uuid:      article.Uuid,
	}, nil
}

func (r *articleRepo) GetArticleList(ctx context.Context, page int32) ([]*biz.Article, error) {
	article, err := r.getArticleFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getArticleFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setArticleToCache("article", article)
		})()
	}
	return article, nil
}

func (r *articleRepo) GetArticleAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "article_agree_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge article agree member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *articleRepo) GetArticleCollectJudge(ctx context.Context, id int32, uuid string) (bool, error) {
	ids := strconv.Itoa(int(id))
	judge, err := r.data.redisCli.SIsMember(ctx, "article_collect_"+ids, uuid).Result()
	if err != nil {
		return false, errors.Wrapf(err, fmt.Sprintf("fail to judge article collect member: id(%v), uuid(%s)", id, uuid))
	}
	return judge, nil
}

func (r *articleRepo) getArticleFromDB(ctx context.Context, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid").Where("auth", 1).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from db: page(%v)", page))
	}

	article := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) GetArticleListHot(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	article, err := r.getArticleHotFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.GetArticleHotFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setArticleHotToCache("article_hot", article)
		})()
	}
	return article, nil
}

func (r *articleRepo) GetColumnArticleList(ctx context.Context, id int32) ([]*biz.Article, error) {
	article, err := r.getColumnArticleFromCache(ctx, id)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getColumnArticleFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setColumnArticleToCache(id, article)
		})()
	}
	return article, nil
}

func (r *articleRepo) GetArticleCount(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Article{}).Where("uuid = ?", uuid).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get article count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *articleRepo) GetArticleCountVisitor(ctx context.Context, uuid string) (int32, error) {
	var count int64
	err := r.data.db.WithContext(ctx).Model(&Article{}).Where("uuid = ? and auth = ?", uuid, 1).Count(&count).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get article count from db: uuid(%s)", uuid))
	}
	return int32(count), nil
}

func (r *articleRepo) GetUserArticleList(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	article, err := r.getUserArticleListFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getUserArticleListFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setUserArticleListToCache("user_article_list_"+uuid, article)
		})()
	}
	return article, nil
}

func (r *articleRepo) getUserArticleListFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_article_list_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article list from cache: key(%s), page(%v)", "user_article_list_"+uuid, page))
	}

	article := make([]*biz.Article, 0, len(list))
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

func (r *articleRepo) getUserArticleListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid").Where("uuid = ?", uuid).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article from db: page(%v), uuid(%s)", page, uuid))
	}

	article := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) GetUserArticleListAll(ctx context.Context, uuid string) ([]*biz.Article, error) {
	article, err := r.getUserArticleListAllFromCache(ctx, uuid)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getUserArticleListAllFromDB(ctx, uuid)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserArticleListToCache("user_article_list_all_"+uuid, article)
		})()
	}
	return article, nil
}

func (r *articleRepo) getUserArticleListAllFromCache(ctx context.Context, uuid string) ([]*biz.Article, error) {
	list, err := r.data.redisCli.ZRevRange(ctx, "user_article_list_all_"+uuid, 0, -1).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article list all from cache: key(%s)", "user_article_list_all_"+uuid))
	}

	article := make([]*biz.Article, 0, len(list))
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

func (r *articleRepo) getUserArticleListAllFromDB(ctx context.Context, uuid string) ([]*biz.Article, error) {
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid").Where("uuid = ?", uuid).Order("article_id desc").Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article all from db:, uuid(%s)", uuid))
	}

	article := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) GetUserArticleListVisitor(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	article, err := r.getUserArticleListVisitorFromCache(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size := len(article)
	if size != 0 {
		return article, nil
	}

	article, err = r.getUserArticleListVisitorFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(article)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserArticleListToCache("user_article_list_visitor_"+uuid, article)
		})()
	}
	return article, nil
}

func (r *articleRepo) getUserArticleListVisitorFromCache(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "user_article_list_visitor_"+uuid, index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article list visitor from cache: key(%s), page(%v)", "user_article_list_visitor_"+uuid, page))
	}

	article := make([]*biz.Article, 0, len(list))
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

func (r *articleRepo) getUserArticleListVisitorFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid").Where("uuid = ? and auth = ?", uuid, 1).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article visitor from db: page(%v), uuid(%s)", page, uuid))
	}

	article := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) setUserArticleListToCache(key string, article []*biz.Article) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(article))
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.ArticleId),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v), err(%v)", article, err)
	}
}

func (r *articleRepo) GetArticleHotFromDB(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ArticleStatistic, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid", "agree").Where("auth", 1).Order("agree desc, article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article statistic from db: page(%v)", page))
	}

	article := make([]*biz.ArticleStatistic, 0, len(list))
	for _, item := range list {
		article = append(article, &biz.ArticleStatistic{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
			Agree:     item.Agree,
		})
	}
	return article, nil
}

func (r *articleRepo) GetArticleStatistic(ctx context.Context, id int32, uuid string) (*biz.ArticleStatistic, error) {
	var statistic *biz.ArticleStatistic
	key := "article_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		statistic, err = r.getArticleStatisticFromCache(ctx, key, uuid)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}

	statistic, err = r.getArticleStatisticFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
	go r.data.Recover(newCtx, func(ctx context.Context) {
		r.setArticleStatisticToCache(key, statistic)
	})()

	if statistic.Auth == 2 && statistic.Uuid != uuid {
		return nil, errors.Errorf("fail to get article statistic from cache: no auth")
	}

	return statistic, nil
}

func (r *articleRepo) getArticleStatisticFromDB(ctx context.Context, id int32) (*biz.ArticleStatistic, error) {
	as := &ArticleStatistic{}
	err := r.data.db.WithContext(ctx).Select("uuid", "agree", "collect", "view", "comment", "auth").Where("article_id = ?", id).First(as).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get statistic from db: id(%v)", id))
	}
	return &biz.ArticleStatistic{
		Uuid:    as.Uuid,
		Agree:   as.Agree,
		Collect: as.Collect,
		View:    as.View,
		Comment: as.Comment,
		Auth:    as.Auth,
	}, nil
}

func (r *articleRepo) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*biz.ArticleStatistic, error) {
	exists, unExists, err := r.articleListStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

	articleListStatistic := make([]*biz.ArticleStatistic, 0, cap(exists))
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(exists) == 0 {
			return nil
		}
		return r.getArticleListStatisticFromCache(ctx, exists, &articleListStatistic)
	}))
	g.Go(r.data.GroupRecover(ctx, func(ctx context.Context) error {
		if len(unExists) == 0 {
			return nil
		}
		return r.getArticleListStatisticFromDb(ctx, unExists, &articleListStatistic)
	}))

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return articleListStatistic, nil
}

func (r *articleRepo) articleListStatisticExist(ctx context.Context, ids []int32) ([]int32, []int32, error) {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "article_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if article statistic exist from cache: ids(%v)", ids))
	}

	exists := make([]int32, 0, len(cmd))
	unExists := make([]int32, 0, len(cmd))
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

func (r *articleRepo) getArticleListStatisticFromCache(ctx context.Context, exists []int32, articleListStatistic *[]*biz.ArticleStatistic) error {
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, id := range exists {
			pipe.HMGet(ctx, "article_"+strconv.Itoa(int(id)), "agree", "collect", "view", "comment")
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get article list statistic from cache: ids(%v)", exists))
	}

	for index, item := range cmd {
		val := []int32{0, 0, 0, 0}
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
		*articleListStatistic = append(*articleListStatistic, &biz.ArticleStatistic{
			ArticleId: exists[index],
			Agree:     val[0],
			Collect:   val[1],
			View:      val[2],
			Comment:   val[3],
		})
	}
	return nil
}

func (r *articleRepo) getArticleListStatisticFromDb(ctx context.Context, unExists []int32, articleListStatistic *[]*biz.ArticleStatistic) error {
	list := make([]*ArticleStatistic, 0, cap(unExists))
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid", "agree", "comment", "collect", "view", "auth").Where("article_id IN ?", unExists).Find(&list).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get article statistic list from db: ids(%v)", unExists))
	}

	for _, item := range list {
		*articleListStatistic = append(*articleListStatistic, &biz.ArticleStatistic{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
			Agree:     item.Agree,
			Comment:   item.Comment,
			Collect:   item.Collect,
			View:      item.View,
			Auth:      item.Auth,
		})
	}

	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setArticleListStatisticToCache(list)
		})()
	}

	return nil
}

func (r *articleRepo) setArticleListStatisticToCache(commentList []*ArticleStatistic) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "article_" + strconv.Itoa(int(item.ArticleId))
			pipe.HSetNX(ctx, key, "uuid", item.Uuid)
			pipe.HSetNX(ctx, key, "agree", item.Agree)
			pipe.HSetNX(ctx, key, "comment", item.Comment)
			pipe.HSetNX(ctx, key, "collect", item.Collect)
			pipe.HSetNX(ctx, key, "view", item.View)
			pipe.HSetNX(ctx, key, "auth", item.Auth)
			pipe.Expire(ctx, key, time.Minute*30)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set article statistic to cache, err(%s)", err.Error())
	}
}

func (r *articleRepo) GetArticleDraftList(ctx context.Context, uuid string) ([]*biz.ArticleDraft, error) {
	draftList := make([]*ArticleDraft, 0)
	err := r.data.db.WithContext(ctx).Select("id").Where("uuid = ? and status = ?", uuid, 3).Order("id desc").Find(&draftList).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get draft list : uuid(%s)", uuid))
	}
	reply := make([]*biz.ArticleDraft, 0, len(draftList))
	for _, item := range draftList {
		reply = append(reply, &biz.ArticleDraft{
			Id: int32(item.ID),
		})
	}
	return reply, nil
}

func (r *articleRepo) GetArticleSearch(ctx context.Context, page int32, search, time string) ([]*biz.ArticleSearch, int32, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	reply := make([]*biz.ArticleSearch, 0)

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
							"fields": []string{"text", "tags", "title"},
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
					"text": map[string]interface{}{
						"fragment_size":       300,
						"number_of_fragments": 1,
						"no_match_size":       300,
						"pre_tags":            "<span style='color:red'>",
						"post_tags":           "</span>",
					},
					"title": map[string]interface{}{
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
					"text": map[string]interface{}{
						"fragment_size":       300,
						"number_of_fragments": 1,
						"no_match_size":       300,
						"pre_tags":            "<span style='color:red'>",
						"post_tags":           "</span>",
					},
					"title": map[string]interface{}{
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
		r.data.elasticSearch.es.Search.WithIndex("article"),
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
			return nil, 0, errors.Errorf(fmt.Sprintf("error search article from  es: reason(%v), page(%v), search(%s), time(%s)", e, page, search, time))
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

		article := &biz.ArticleSearch{
			Id:     int32(id),
			Tags:   hit.(map[string]interface{})["_source"].(map[string]interface{})["tags"].(string),
			Update: hit.(map[string]interface{})["_source"].(map[string]interface{})["update"].(string),
			Cover:  hit.(map[string]interface{})["_source"].(map[string]interface{})["cover"].(string),
			Uuid:   hit.(map[string]interface{})["_source"].(map[string]interface{})["uuid"].(string),
		}

		if text, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["text"]; ok {
			article.Text = text.([]interface{})[0].(string)
		}

		if title, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["title"]; ok {
			article.Title = title.([]interface{})[0].(string)
		}

		reply = append(reply, article)
	}
	res.Body.Close()
	return reply, int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), nil
}

func (r *articleRepo) GetArticleAuth(ctx context.Context, id int32) (int32, error) {
	article := &Article{}
	err := r.data.db.WithContext(ctx).Select("auth").Where("article_id = ?", id).First(article).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get article auth from db: id(%v)", id))
	}
	return article.Auth, nil
}

func (r *articleRepo) GetUserArticleAgree(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userArticleAgreeExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserArticleAgreeFromCache(ctx, uuid)
	} else {
		return r.getUserArticleAgreeFromDb(ctx, uuid)
	}
}

func (r *articleRepo) userArticleAgreeExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_article_agree_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user article agree exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *articleRepo) getUserArticleAgreeFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_article_agree_" + uuid
	agreeSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article agree from cache: uuid(%s)", uuid))
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

func (r *articleRepo) getUserArticleAgreeFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*ArticleAgree, 0)
	err := r.data.db.WithContext(ctx).Select("article_id").Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article agree from db: uuid(%s)", uuid))
	}

	agreeMap := make(map[int32]bool, 0)
	for _, item := range list {
		agreeMap[item.ArticleId] = true
	}
	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserArticleAgreeToCache(uuid, list)
		})()
	}
	return agreeMap, nil
}

func (r *articleRepo) setUserArticleAgreeToCache(uuid string, agreeList []*ArticleAgree) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0, len(agreeList))
		key := "user_article_agree_" + uuid
		for _, item := range agreeList {
			set = append(set, item.ArticleId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user article agree to cache: uuid(%s), agreeList(%v), error(%s) ", uuid, agreeList, err.Error())
	}
}

func (r *articleRepo) GetUserArticleCollect(ctx context.Context, uuid string) (map[int32]bool, error) {
	exist, err := r.userArticleCollectExist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if exist == 1 {
		return r.getUserArticleCollectFromCache(ctx, uuid)
	} else {
		return r.getUserArticleCollectFromDb(ctx, uuid)
	}
}

func (r *articleRepo) userArticleCollectExist(ctx context.Context, uuid string) (int32, error) {
	exist, err := r.data.redisCli.Exists(ctx, "user_article_collect_"+uuid).Result()
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to check if user article collect exist from cache: uuid(%s)", uuid))
	}
	return int32(exist), nil
}

func (r *articleRepo) getUserArticleCollectFromCache(ctx context.Context, uuid string) (map[int32]bool, error) {
	key := "user_article_collect_" + uuid
	collectSet, err := r.data.redisCli.SMembers(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article collect from cache: uuid(%s)", uuid))
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

func (r *articleRepo) getUserArticleCollectFromDb(ctx context.Context, uuid string) (map[int32]bool, error) {
	list := make([]*ArticleCollect, 0)
	err := r.data.db.WithContext(ctx).Select("article_id").Where("uuid = ? and status = ?", uuid, 1).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article collect from db: uuid(%s)", uuid))
	}

	collectMap := make(map[int32]bool, 0)
	for _, item := range list {
		collectMap[item.ArticleId] = true
	}
	if len(list) != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setUserArticleCollectToCache(uuid, list)
		})()
	}
	return collectMap, nil
}

func (r *articleRepo) setUserArticleCollectToCache(uuid string, collectList []*ArticleCollect) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		set := make([]interface{}, 0, len(collectList))
		key := "user_article_collect_" + uuid
		for _, item := range collectList {
			set = append(set, item.ArticleId)
		}
		pipe.SAdd(ctx, key, set...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set user article collect to cache: uuid(%s), collectList(%v), error(%s) ", uuid, collectList, err.Error())
	}
}

func (r *articleRepo) SetArticleImageIrregular(ctx context.Context, review *biz.ImageReview) (*biz.ImageReview, error) {
	ar := &ArticleReview{
		ArticleId: review.CreationId,
		Kind:      review.Kind,
		Uid:       review.Uid,
		Uuid:      review.Uuid,
		JobId:     review.JobId,
		Url:       review.Url,
		Label:     review.Label,
		Result:    review.Result,
		Category:  review.Category,
		SubLabel:  review.SubLabel,
		Score:     review.Score,
	}
	err := r.data.DB(ctx).Select("ArticleId", "Kind", "Uid", "Uuid", "JobId", "Url", "Label", "Result", "Category", "SubLabel", "Score").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to add article image review record: review(%v)", review))
	}
	review.Id = int32(ar.ID)
	review.CreateAt = time.Now().Format("2006-01-02")
	return review, nil
}

func (r *articleRepo) SetArticleContentIrregular(ctx context.Context, review *biz.TextReview) (*biz.TextReview, error) {
	ar := &ArticleContentReview{
		ArticleId: review.CreationId,
		Title:     review.Title,
		Kind:      review.Kind,
		Uuid:      review.Uuid,
		JobId:     review.JobId,
		Label:     review.Label,
		Result:    review.Result,
		Section:   review.Section,
	}
	err := r.data.DB(ctx).Select("ArticleId", "Title", "Kind", "Uuid", "JobId", "Label", "Result", "Section").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to add article content review record: review(%v)", review))
	}
	review.Id = int32(ar.ID)
	review.CreateAt = time.Now().Format("2006-01-02")
	return review, nil
}

func (r *articleRepo) SetArticleImageIrregularToCache(ctx context.Context, review *biz.ImageReview) error {
	marshal, err := review.MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set article image irregular to json: json.Marshal(%v)", review))
	}
	keys := []string{"article_image_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = r.data.redisCli.EvalSha(ctx, "8f6205011a2b264278a7c5bc0a2bcd1006ac6e5d", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set article image irregular to cache: review(%v)", review))
	}
	return nil
}

func (r *articleRepo) SetArticleContentIrregularToCache(ctx context.Context, review *biz.TextReview) error {
	marshal, err := review.MarshalJSON()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set article content irregular to json: json.Marshal(%v)", review))
	}
	keys := []string{"article_content_irregular_" + review.Uuid}
	values := []interface{}{marshal}
	_, err = r.data.redisCli.EvalSha(ctx, "8f6205011a2b264278a7c5bc0a2bcd1006ac6e5d", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set article content irregular to cache: review(%v)", review))
	}
	return nil
}

func (r *articleRepo) GetCollectionsIdFromArticleCollect(ctx context.Context, id int32) (int32, error) {
	collect := &Collect{}
	err := r.data.db.WithContext(ctx).Select("collections_id").Where("creations_id = ? and mode = ?", id, 1).First(collect).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to get collections id  from db: creationsId(%v)", id))
	}
	return collect.CollectionsId, nil
}

func (r *articleRepo) GetArticleImageReview(ctx context.Context, page int32, uuid string) ([]*biz.ImageReview, error) {
	key := "article_image_irregular_" + uuid
	review, err := r.getArticleImageReviewFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getArticleImageReviewFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setArticleImageReviewToCache(key, review)
		})()
	}
	return review, nil
}

func (r *articleRepo) getArticleImageReviewFromCache(ctx context.Context, page int32, key string) ([]*biz.ImageReview, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article image irregular list from cache: key(%s), page(%v)", key, page))
	}

	review := make([]*biz.ImageReview, 0, len(list))
	for _index, item := range list {
		var imageReview = &biz.ImageReview{}
		err = imageReview.UnmarshalJSON([]byte(item))
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: imageReview(%v)", item))
		}
		review = append(review, &biz.ImageReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: imageReview.CreationId,
			Kind:       imageReview.Kind,
			Uid:        imageReview.Uid,
			Uuid:       imageReview.Uuid,
			CreateAt:   imageReview.CreateAt,
			JobId:      imageReview.JobId,
			Url:        imageReview.Url,
			Label:      imageReview.Label,
			Result:     imageReview.Result,
			Category:   imageReview.Category,
			SubLabel:   imageReview.SubLabel,
			Score:      imageReview.Score,
		})
	}
	return review, nil
}

func (r *articleRepo) getArticleImageReviewFromDB(ctx context.Context, page int32, uuid string) ([]*biz.ImageReview, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ArticleReview, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "kind", "uid", "uuid", "job_id", "created_at", "url", "label", "result", "category", "sub_label", "score").Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article image review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.ImageReview, 0, len(list))
	for _index, item := range list {
		review = append(review, &biz.ImageReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: item.ArticleId,
			Kind:       item.Kind,
			Uid:        item.Uid,
			Uuid:       item.Uuid,
			JobId:      item.JobId,
			CreateAt:   item.CreatedAt.Format("2006-01-02"),
			Url:        item.Url,
			Label:      item.Label,
			Result:     item.Result,
			Category:   item.Category,
			SubLabel:   item.SubLabel,
			Score:      item.Score,
		})
	}
	return review, nil
}

func (r *articleRepo) setArticleImageReviewToCache(key string, review []*biz.ImageReview) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0, len(review))
		for _, item := range review {
			m, err := item.MarshalJSON()
			if err != nil {
				return errors.Wrapf(err, "fail to marshal avatar review: imageReview(%v)", review)
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article image review to cache: imageReview(%v), err(%v)", review, err)
	}
}

func (r *articleRepo) GetArticleContentReview(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	key := "article_content_irregular_" + uuid
	review, err := r.getArticleContentReviewFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getArticleContentReviewFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setArticleContentReviewToCache(key, review)
		})()
	}
	return review, nil
}

func (r *articleRepo) getArticleContentReviewFromCache(ctx context.Context, page int32, key string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article content irregular list from cache: key(%s), page(%v)", key, page))
	}

	review := make([]*biz.TextReview, 0, len(list))
	for _index, item := range list {
		var textReview = &biz.TextReview{}
		err = textReview.UnmarshalJSON([]byte(item))
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

func (r *articleRepo) getArticleContentReviewFromDB(ctx context.Context, page int32, uuid string) ([]*biz.TextReview, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ArticleContentReview, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "kind", "title", "uuid", "job_id", "created_at", "label", "result", "section").Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article content review from db: page(%v), uuid(%s)", page, uuid))
	}

	review := make([]*biz.TextReview, 0, len(list))
	for _index, item := range list {
		review = append(review, &biz.TextReview{
			Id:         int32(_index+1) + (page-1)*20,
			CreationId: item.ArticleId,
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

func (r *articleRepo) setArticleContentReviewToCache(key string, review []*biz.TextReview) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0, len(review))
		for _, item := range review {
			m, err := item.MarshalJSON()
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: contentReview(%v)", review))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article content review to cache: contentReview(%v), err(%v)", review, err)
	}
}

func (r *articleRepo) CreateArticle(ctx context.Context, id, auth int32, uuid string) error {
	article := &Article{
		ArticleId: id,
		Uuid:      uuid,
		Auth:      auth,
	}
	err := r.data.DB(ctx).Select("ArticleId", "Uuid", "Auth").Create(article).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a article: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) DeleteArticle(ctx context.Context, id int32, uuid string) error {
	article := &Article{}
	err := r.data.DB(ctx).Where("article_id = ? and uuid = ?", id, uuid).Delete(article).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete a article: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *articleRepo) CreateArticleDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &ArticleDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create an article draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *articleRepo) DeleteArticleStatistic(ctx context.Context, id int32, uuid string) error {
	statistic := &ArticleStatistic{}
	err := r.data.DB(ctx).Where("article_id = ? and uuid = ?", id, uuid).Delete(statistic).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete an article statistic: uuid(%s)", uuid))
	}
	return nil
}

func (r *articleRepo) CreateArticleFolder(ctx context.Context, id int32, uuid string) error {
	name := "article/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create an article folder: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CreateArticleStatistic(ctx context.Context, id, auth int32, uuid string) error {
	as := &ArticleStatistic{
		ArticleId: id,
		Uuid:      uuid,
		Auth:      auth,
	}
	err := r.data.DB(ctx).Select("ArticleId", "Uuid", "Auth").Create(as).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a article statistic: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CreateArticleCache(ctx context.Context, id, auth int32, uuid, mode string) error {
	ids := strconv.Itoa(int(id))
	articleStatistic := "article_" + ids
	article := "article"
	articleHot := "article_hot"
	leaderboard := "leaderboard"
	userArticleList := "user_article_list_" + uuid
	userArticleListVisitor := "user_article_list_visitor_" + uuid
	creationUser := "creation_user_" + uuid
	creationUserVisitor := "creation_user_visitor_" + uuid
	keys := []string{articleStatistic, article, articleHot, leaderboard, userArticleList, userArticleListVisitor, creationUser, creationUserVisitor}
	values := []interface{}{uuid, auth, id, ids + "%" + uuid, mode, ids + "%" + uuid + "%article"}
	_, err := r.data.redisCli.EvalSha(ctx, "f1dcbca8b2c0cebc82fb4f1ab03fe38a1e724806", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create(update) article cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) UpdateArticleCache(ctx context.Context, id, auth int32, uuid, mode string) error {
	return r.CreateArticleCache(ctx, id, auth, uuid, mode)
}

func (r *articleRepo) DeleteArticleCache(ctx context.Context, id, auth int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	keys := []string{"article", "article_hot", "leaderboard", "article_" + ids, "article_collect_" + ids, "user_article_list_" + uuid, "user_article_list_visitor_" + uuid, "creation_user_" + uuid, "creation_user_visitor_" + uuid}
	values := []interface{}{ids + "%" + uuid, ids + "%" + uuid + "%article", auth}
	_, err := r.data.redisCli.EvalSha(ctx, "a00173f1947bc4a40039ed55303d230cd2daa749", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete article cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *articleRepo) FreezeArticleCos(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "article/" + uuid + "/" + ids + "/content"
	_, err := r.data.cosCli.Object.Delete(ctx, key)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to freeze article: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *articleRepo) CreateArticleSearch(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "article/" + uuid + "/" + ids + "/search"
	resp, err := r.data.cosCli.Object.Get(ctx, key, &cos.ObjectGetOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get article from cos: id(%v), uuid(%s)", id, uuid))
	}

	article, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read request body: id(%v), uuid(%s)", id, uuid))
	}

	resp.Body.Close()

	req := esapi.IndexRequest{
		Index:      "article",
		DocumentID: strconv.Itoa(int(id)),
		Body:       bytes.NewReader(article),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("error getting article search create response: id(%v), uuid(%s)", id, uuid))
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

func (r *articleRepo) AddArticleComment(ctx context.Context, id int32) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ?", id).Update("comment", gorm.Expr("comment + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article comment: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) AddArticleCommentToCache(ctx context.Context, id int32, uuid string) error {
	key := "article_" + strconv.Itoa(int(id))
	keys := []string{key}
	var values []interface{}
	_, err := r.data.redisCli.EvalSha(ctx, "17cef96499553bacd3611f03630542e2b475a42c", keys, values...).Result()
	if err != nil {
		r.log.Errorf("fail to add article comment to cache: id(%v), uuid(%s), err(%v)", id, uuid, err)
	}
	return nil
}

func (r *articleRepo) AddCreationUserArticle(ctx context.Context, uuid string, auth int32) error {
	cu := &CreationUser{
		Uuid:    uuid,
		Article: 1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"article": gorm.Expr("article + ?", 1)}),
	}).Create(cu).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation article: uuid(%v)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := &CreationUserVisitor{
		Uuid:    uuid,
		Article: 1,
	}
	err = r.data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"article": gorm.Expr("article + ?", 1)}),
	}).Create(cuv).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add creation article visitor: uuid(%v)", uuid))
	}

	return nil
}

func (r *articleRepo) ReduceArticleComment(ctx context.Context, id int32) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and comment > 0", id).Update("comment", gorm.Expr("comment - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article comment: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) ReduceArticleCommentToCache(ctx context.Context, id int32, uuid string) error {
	key := "article_" + strconv.Itoa(int(id))
	keys := []string{key}
	_, err := r.data.redisCli.EvalSha(ctx, "e9942f8b94d123b097211f27b5bf8fc896b1fa45", keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article comment to cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *articleRepo) ReduceCreationUserArticle(ctx context.Context, auth int32, uuid string) error {
	cu := CreationUser{}
	err := r.data.DB(ctx).Model(&cu).Where("uuid = ? and article > 0", uuid).Update("article", gorm.Expr("article - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user article: uuid(%s)", uuid))
	}

	if auth == 2 {
		return nil
	}

	cuv := CreationUserVisitor{}
	err = r.data.DB(ctx).Model(&cuv).Where("uuid = ? and article > 0", uuid).Update("article", gorm.Expr("article - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user article visitor: uuid(%s)", uuid))
	}
	return nil
}

func (r *articleRepo) DeleteArticleSearch(ctx context.Context, id int32, uuid string) error {
	req := esapi.DeleteRequest{
		Index:      "article",
		DocumentID: strconv.Itoa(int(id)),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting article search delete response: id(%v), uuid(%s)", id, uuid))
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

func (r *articleRepo) EditArticleCos(ctx context.Context, id int32, uuid string) error {
	err := r.EditArticleCosContent(ctx, id, uuid)
	if err != nil {
		return err
	}

	err = r.EditArticleCosIntroduce(ctx, id, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *articleRepo) EditArticleCosContent(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "article/" + uuid + "/" + ids + "/content-edit"
	dest := "article/" + uuid + "/" + ids + "/content"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy article from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) EditArticleCosIntroduce(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "article/" + uuid + "/" + ids + "/introduce-edit"
	dest := "article/" + uuid + "/" + ids + "/introduce"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy article from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) EditArticleSearch(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	key := "article/" + uuid + "/" + ids + "/search"
	resp, err := r.data.cosCli.Object.Get(ctx, key, &cos.ObjectGetOptions{})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get article from cos: id(%v), uuid(%s)", id, uuid))
	}

	article, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read request body: id(%v), uuid(%s)", id, uuid))
	}

	resp.Body.Close()

	req := esapi.IndexRequest{
		Index:      "article",
		DocumentID: strconv.Itoa(int(id)),
		Body:       bytes.NewReader(article),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, r.data.elasticSearch.es)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("Error getting article search edit response: id(%v), uuid(%s)", id, uuid))
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

func (r *articleRepo) DeleteArticleDraft(ctx context.Context, id int32, uuid string) error {
	ad := &ArticleDraft{}
	ad.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(ad).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete article draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) ArticleDraftMark(ctx context.Context, id int32, uuid string) error {
	err := r.data.db.WithContext(ctx).Model(&ArticleDraft{}).Where("id = ? and uuid = ?", id, uuid).Update("status", 3).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 3: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) SendArticle(ctx context.Context, id int32, uuid string) (*biz.ArticleDraft, error) {
	ad := &ArticleDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&ArticleDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 3).Updates(ad).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.ArticleDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *articleRepo) SendReviewToMq(ctx context.Context, review *biz.ArticleReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *articleRepo) SendScoreToMq(ctx context.Context, score int32, uuid, mode string) error {
	scoreMap := &biz.SendScoreMap{
		Uuid:  uuid,
		Score: score,
		Mode:  mode,
	}
	data, err := scoreMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send score to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *articleRepo) SendArticleToMq(ctx context.Context, article *biz.Article, mode string) error {
	articleMap := &biz.SendArticleMap{
		Uuid: article.Uuid,
		Id:   article.ArticleId,
		Auth: article.Auth,
		Mode: mode,
	}
	data, err := articleMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{article.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article to mq: %v", article))
	}
	return nil
}

func (r *articleRepo) SendStatisticToMq(ctx context.Context, id, collectionsId int32, uuid, userUuid, mode string) error {
	statisticMap := &biz.SendStatisticMap{
		Id:            id,
		CollectionsId: collectionsId,
		Uuid:          uuid,
		UserUuid:      userUuid,
		Mode:          mode,
	}
	data, err := statisticMap.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article statistic to mq: map(%v)", statisticMap))
	}
	return nil
}

func (r *articleRepo) SetArticleAgree(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article agree: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetUserArticleAgree(ctx context.Context, id int32, userUuid string) error {
	aa := &ArticleAgree{
		ArticleId: id,
		Uuid:      userUuid,
		Status:    1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(aa).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user article agree: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) SetArticleAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	hotKey := fmt.Sprintf("article_hot")
	statisticKey := fmt.Sprintf("article_%v", id)
	boardKey := fmt.Sprintf("leaderboard")
	userKey := fmt.Sprintf("user_article_agree_%s", userUuid)
	keys := []string{hotKey, statisticKey, boardKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%article"), id}
	_, err := r.data.redisCli.EvalSha(ctx, "7c94d6c820881b0b69ec949e5565ea0557be65ae", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user article agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid))
	}
	return nil
}

func (r *articleRepo) SetArticleView(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("view", gorm.Expr("view + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article view: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetArticleViewToCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	keys := []string{"article_" + ids}
	values := []interface{}{1}
	_, err := r.data.redisCli.EvalSha(ctx, "6abf0f5d9afbfe7ed552297f59b9d898c6b06747", keys, values...).Result()

	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article view to cache: id(%v), uuid(%s)", id, uuid))
	}
	return nil
}

func (r *articleRepo) SetCollectionsArticleCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
	collect := &Collect{
		CollectionsId: collectionsId,
		Uuid:          userUuid,
		CreationsId:   id,
		Mode:          1,
		Status:        1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to collect an article: article_id(%v), collectionsId(%v), userUuid(%s)", id, collectionsId, userUuid))
	}
	return nil
}

func (r *articleRepo) SetCollectionArticle(ctx context.Context, collectionsId int32, userUuid string) error {
	c := Collections{}
	err := r.data.DB(ctx).Model(&c).Where("collections_id = ? and uuid = ?", collectionsId, userUuid).Update("article", gorm.Expr("article + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add collection article collect: collectionsId(%v), userUuid(%s)", collectionsId, userUuid))
	}
	return nil
}

func (r *articleRepo) SetUserArticleCollect(ctx context.Context, id int32, userUuid string) error {
	ac := &ArticleCollect{
		ArticleId: id,
		Uuid:      userUuid,
		Status:    1,
	}
	err := r.data.DB(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{"status": 1}),
	}).Create(ac).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add user article collect: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) SetCreationUserCollect(ctx context.Context, userUuid string) error {
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

func (r *articleRepo) SetArticleCollect(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article collect: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetArticleCollectToCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("article_%v", id)
	collectKey := fmt.Sprintf("collections_%v_article", collectionsId)
	collectionsKey := fmt.Sprintf("collections_%v", collectionsId)
	creationKey := fmt.Sprintf("creation_user_%s", userUuid)
	userKey := fmt.Sprintf("user_article_collect_%s", userUuid)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid)}
	_, err := r.data.redisCli.EvalSha(ctx, "958864163ea9699b2c0ee77a49ea456f0db8b45a", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article collect to cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
	}
	return nil
}

func (r *articleRepo) SetUserArticleAgreeToCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_article_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "16415c44ae0544f8c7f85d841521813d41c35994", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user article agree to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) SetUserArticleCollectToCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_article_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "16415c44ae0544f8c7f85d841521813d41c35994", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set user article collect to cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelArticleAgree(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ? and agree > 0", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article agree: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CancelUserArticleAgree(ctx context.Context, id int32, userUuid string) error {
	aa := ArticleAgree{}
	err := r.data.DB(ctx).Model(&aa).Where("article_id = ? and uuid = ?", id, userUuid).Update("status", 2).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user article agree: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelArticleAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	hotKey := "article_hot"
	boardKey := "leaderboard"
	statisticKey := fmt.Sprintf("article_%v", id)
	userKey := fmt.Sprintf("user_article_agree_%s", userUuid)

	keys := []string{hotKey, boardKey, statisticKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), fmt.Sprintf("%v%s%s%s", id, "%", uuid, "%article"), id}
	_, err := r.data.redisCli.EvalSha(ctx, "c2c68100b9b682b94faff1a9b182159921d0b8e4", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelCollectionsArticleCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&Collect{}).Where("creations_id = ? and mode = ? and uuid = ?", id, 1, userUuid).Updates(collect).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect: article_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelUserArticleCollect(ctx context.Context, id int32, userUuid string) error {
	ac := &ArticleCollect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&ArticleCollect{}).Where("article_id = ? and uuid = ?", id, userUuid).Updates(ac).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user article collect: article_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelCollectionArticle(ctx context.Context, id int32, uuid string) error {
	collections := &Collections{}
	err := r.data.DB(ctx).Model(collections).Where("collections_id = ? and uuid = ? and article > 0", id, uuid).Update("article", gorm.Expr("article - ?", 1)).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel collections article: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) ReduceCreationUserCollect(ctx context.Context, uuid string) error {
	cu := &CreationUser{}
	err := r.data.DB(ctx).Model(cu).Where("uuid = ? and collect > 0", uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce creation user collect: uuid(%v)", uuid))
	}
	return nil
}

func (r *articleRepo) CancelArticleCollect(ctx context.Context, id int32, uuid string) error {
	as := &ArticleStatistic{}
	err := r.data.DB(ctx).Model(as).Where("article_id = ? and uuid = ? and collect > 0", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CancelArticleCollectFromCache(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	statisticKey := fmt.Sprintf("article_%v", id)
	collectKey := fmt.Sprintf("collections_%v_article", collectionsId)
	collectionsKey := fmt.Sprintf("collections_%v", collectionsId)
	creationKey := fmt.Sprintf("creation_user_%s", userUuid)
	userKey := fmt.Sprintf("user_article_collect_%s", userUuid)
	keys := []string{statisticKey, collectKey, collectionsKey, creationKey, userKey}
	values := []interface{}{fmt.Sprintf("%v%s%s", id, "%", uuid), id}
	_, err := r.data.redisCli.EvalSha(ctx, "51c5960205bc0710a94f1520e748701724a72608", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect from cache: id(%v), collectionsId(%v), uuid(%s), userUuid(%s)", id, collectionsId, uuid, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelUserArticleAgreeFromCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_article_agree_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "9aab0f65688566d3d54b7851570412ad816a80bb", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user article agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelUserArticleCollectFromCache(ctx context.Context, id int32, userUuid string) error {
	keys := []string{"user_article_collect_" + userUuid}
	values := []interface{}{strconv.Itoa(int(id))}
	_, err := r.data.redisCli.EvalSha(ctx, "d7d8a559ccbad93b64d6e4ef0b9130d2b5992224", keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel user article agree from cache: id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) SendArticleStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error {
	achievement := &biz.SendArticleStatisticMap{
		Uuid:     uuid,
		UserUuid: userUuid,
		Mode:     mode,
	}
	data, err := achievement.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article statistic to mq: uuid(%s)", uuid))
	}
	return nil
}

func (r *articleRepo) SendArticleImageIrregularToMq(ctx context.Context, review *biz.ImageReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article image review to mq: %v", err))
	}
	return nil
}

func (r *articleRepo) SendArticleContentIrregularToMq(ctx context.Context, review *biz.TextReview) error {
	data, err := review.MarshalJSON()
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "matrix",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.mqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article content review to mq: %v", err))
	}
	return nil
}

func (r *articleRepo) getArticleFromCache(ctx context.Context, page int32) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "article", index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from cache: key(%s), page(%v)", "article", page))
	}

	article := make([]*biz.Article, 0, len(list))
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

func (r *articleRepo) getArticleHotFromCache(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "article_hot", index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article hot from cache: key(%s), page(%v)", "article_hot", page))
	}

	article := make([]*biz.ArticleStatistic, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		article = append(article, &biz.ArticleStatistic{
			ArticleId: int32(id),
			Uuid:      member[1],
		})
	}
	return article, nil
}

func (r *articleRepo) getColumnArticleFromCache(ctx context.Context, id int32) ([]*biz.Article, error) {
	ids := strconv.Itoa(int(id))
	list, err := r.data.redisCli.ZRevRange(ctx, "column_includes_"+ids, 0, -1).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column article from cache: columnId(%v)", id))
	}

	article := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		member := strings.Split(item, "%")
		articleId, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		article = append(article, &biz.Article{
			ArticleId: int32(articleId),
			Uuid:      member[1],
		})
	}
	return article, nil
}

func (r *articleRepo) getColumnArticleFromDB(ctx context.Context, id int32) ([]*biz.Article, error) {
	list := make([]*ColumnInclusion, 0)
	err := r.data.db.WithContext(ctx).Select("article_id", "uuid").Where("column_id = ? and status = ?", id, 1).Order("updated_at desc").Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column article from db: columnId(%v)", id))
	}

	article := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) setArticleToCache(key string, article []*biz.Article) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(article))
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.ArticleId),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v), err(%v)", article, err)
	}
}

func (r *articleRepo) setArticleHotToCache(key string, article []*biz.ArticleStatistic) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(article))
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v), err(%v)", article, err)
	}
}

func (r *articleRepo) setColumnArticleToCache(id int32, article []*biz.Article) {
	ids := strconv.Itoa(int(id))
	length := len(article)
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(article))
		for index, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(length - index),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(ctx, "column_includes_"+ids, z...)
		pipe.Expire(ctx, "column_includes_"+ids, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set column article to cache: article(%v), err(%v)", article, err)
	}
}

func (r *articleRepo) getArticleStatisticFromCache(ctx context.Context, key, uuid string) (*biz.ArticleStatistic, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "uuid", "agree", "collect", "view", "comment", "auth").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article statistic form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0, 0}
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
	if val[4] == 2 && statistic[0].(string) != uuid {
		return nil, errors.Errorf("fail to get article statistic from cache: no auth")
	}
	return &biz.ArticleStatistic{
		Uuid:    statistic[0].(string),
		Agree:   val[0],
		Collect: val[1],
		View:    val[2],
		Comment: val[3],
	}, nil
}

func (r *articleRepo) setArticleStatisticToCache(key string, statistic *biz.ArticleStatistic) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View, "comment", statistic.Comment)
		pipe.Expire(ctx, key, time.Minute*30)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article statistic to cache, err(%s)", err.Error())
	}
}
