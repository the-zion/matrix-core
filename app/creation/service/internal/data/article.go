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
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("article draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
	}
	return &biz.ArticleDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *articleRepo) GetArticle(ctx context.Context, id int32) (*biz.Article, error) {
	article := &Article{}
	err := r.data.db.WithContext(ctx).Where("article_id = ?", id).First(article).Error
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
		go r.setArticleToCache("article", article)
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
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article from db: page(%v)", page))
	}

	article := make([]*biz.Article, 0)
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
		go r.setArticleHotToCache("article_hot", article)
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
		go r.setColumnArticleToCache(id, article)
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
		go r.setUserArticleListToCache("user_article_list_"+uuid, article)
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
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article list from cache: key(%s), page(%v)", "user_article_list_", page))
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

func (r *articleRepo) getUserArticleListFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article from db: page(%v), uuid(%s)", page, uuid))
	}

	article := make([]*biz.Article, 0)
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
		go r.setUserArticleListToCache("user_article_list_visitor_"+uuid, article)
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
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article list visitor from cache: key(%s), page(%v)", "user_article_list_visitor_", page))
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

func (r *articleRepo) getUserArticleListVisitorFromDB(ctx context.Context, page int32, uuid string) ([]*biz.Article, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Article, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and auth = ?", uuid, 1).Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user article visitor from db: page(%v), uuid(%s)", page, uuid))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) setUserArticleListToCache(key string, article []*biz.Article) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.ArticleId),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		pipe.Expire(context.Background(), key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v)", article)
	}
}

func (r *articleRepo) GetArticleHotFromDB(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*ArticleStatistic, 0)
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("agree desc, article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article statistic from db: page(%v)", page))
	}

	article := make([]*biz.ArticleStatistic, 0)
	for _, item := range list {
		article = append(article, &biz.ArticleStatistic{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
			Agree:     item.Agree,
		})
	}
	return article, nil
}

func (r *articleRepo) GetArticleStatistic(ctx context.Context, id int32) (*biz.ArticleStatistic, error) {
	var statistic *biz.ArticleStatistic
	key := "article_" + strconv.Itoa(int(id))
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to judge if key exist or not from cache: key(%s)", key))
	}

	if exist == 1 {
		statistic, err = r.getArticleStatisticFromCache(ctx, key)
		if err != nil {
			return nil, err
		}
		return statistic, nil
	}

	statistic, err = r.getArticleStatisticFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	go r.setArticleStatisticToCache(key, statistic)

	return statistic, nil
}

func (r *articleRepo) getArticleStatisticFromDB(ctx context.Context, id int32) (*biz.ArticleStatistic, error) {
	as := &ArticleStatistic{}
	err := r.data.db.WithContext(ctx).Where("article_id = ?", id).First(as).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("faile to get statistic from db: id(%v)", id))
	}
	return &biz.ArticleStatistic{
		Uuid:    as.Uuid,
		Agree:   as.Agree,
		Collect: as.Collect,
		View:    as.View,
		Comment: as.Comment,
	}, nil
}

func (r *articleRepo) GetArticleListStatistic(ctx context.Context, ids []int32) ([]*biz.ArticleStatistic, error) {
	articleListStatistic := make([]*biz.ArticleStatistic, 0)
	exists, unExists, err := r.articleListStatisticExist(ctx, ids)
	if err != nil {
		return nil, err
	}

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
	exists := make([]int32, 0)
	unExists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range ids {
			pipe.Exists(ctx, "article_"+strconv.Itoa(int(item)))
		}
		return nil
	})
	if err != nil {
		return nil, nil, errors.Wrapf(err, fmt.Sprintf("fail to check if article statistic exist from cache: ids(%v)", ids))
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
	list := make([]*ArticleStatistic, 0)
	err := r.data.db.WithContext(ctx).Where("article_id IN ?", unExists).Find(&list).Error
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
		})
	}

	if len(list) != 0 {
		go r.setArticleListStatisticToCache(list)
	}

	return nil
}

func (r *articleRepo) setArticleListStatisticToCache(commentList []*ArticleStatistic) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for _, item := range commentList {
			key := "article_" + strconv.Itoa(int(item.ArticleId))
			pipe.HSetNX(context.Background(), key, "uuid", item.Uuid)
			pipe.HSetNX(context.Background(), key, "agree", item.Agree)
			pipe.HSetNX(context.Background(), key, "comment", item.Comment)
			pipe.HSetNX(context.Background(), key, "collect", item.Collect)
			pipe.HSetNX(context.Background(), key, "view", item.View)
			pipe.Expire(context.Background(), key, time.Hour*8)
		}
		return nil
	})

	if err != nil {
		r.log.Errorf("fail to set article statistic to cache, err(%s)", err.Error())
	}
}

func (r *articleRepo) GetArticleDraftList(ctx context.Context, uuid string) ([]*biz.ArticleDraft, error) {
	reply := make([]*biz.ArticleDraft, 0)
	draftList := make([]*ArticleDraft, 0)
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 3).Order("id desc").Find(&draftList).Error
	if err != nil {
		return reply, errors.Wrapf(err, fmt.Sprintf("fail to get draft list : uuid(%s)", uuid))
	}
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

		reply = append(reply, &biz.ArticleSearch{
			Id:     int32(id),
			Tags:   hit.(map[string]interface{})["_source"].(map[string]interface{})["tags"].(string),
			Update: hit.(map[string]interface{})["_source"].(map[string]interface{})["update"].(string),
			Cover:  hit.(map[string]interface{})["_source"].(map[string]interface{})["cover"].(string),
			Uuid:   hit.(map[string]interface{})["_source"].(map[string]interface{})["uuid"].(string),
			Text:   hit.(map[string]interface{})["highlight"].(map[string]interface{})["text"].([]interface{})[0].(string),
			Title:  hit.(map[string]interface{})["highlight"].(map[string]interface{})["title"].([]interface{})[0].(string),
		})
	}
	res.Body.Close()
	return reply, int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), nil
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

func (r *articleRepo) CreateArticleCache(ctx context.Context, id, auth int32, uuid string) error {
	exists := make([]int32, 0)
	cmd, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Exists(ctx, "article")
		pipe.Exists(ctx, "article_hot")
		pipe.Exists(ctx, "leaderboard")
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to check if article exist from cache: id(%v),uuid(%s)", id, uuid))
	}

	for _, item := range cmd {
		exist := item.(*redis.IntCmd).Val()
		exists = append(exists, int32(exist))
	}

	ids := strconv.Itoa(int(id))
	_, err = r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "article_"+ids, "uuid", uuid)
		pipe.HSetNX(ctx, "article_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "article_"+ids, "collect", 0)
		pipe.HSetNX(ctx, "article_"+ids, "view", 0)
		pipe.HSetNX(ctx, "article_"+ids, "comment", 0)

		if auth == 2 {
			return nil
		}

		if exists[0] == 1 {
			pipe.ZAddNX(ctx, "article", &redis.Z{
				Score:  float64(id),
				Member: ids + "%" + uuid,
			})
		}

		if exists[1] == 1 {
			pipe.ZAddNX(ctx, "article_hot", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid,
			})
		}

		if exists[2] == 1 {
			pipe.ZAddNX(ctx, "leaderboard", &redis.Z{
				Score:  0,
				Member: ids + "%" + uuid + "%article",
			})
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create(update) article cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *articleRepo) UpdateArticleCache(ctx context.Context, id, auth int32, uuid string) error {
	return r.CreateArticleCache(ctx, id, auth, uuid)
}

func (r *articleRepo) DeleteArticleCache(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, "article", ids+"%"+uuid)
		pipe.ZRem(ctx, "article_hot", ids+"%"+uuid)
		pipe.ZRem(ctx, "leaderboard", ids+"%"+uuid+"%article")
		pipe.Del(ctx, "article_"+ids)
		pipe.Del(ctx, "article_collect_"+ids)
		return nil
	})
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
	exist, err := r.data.redisCli.Exists(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to check if article statistic exist from cache: id(%v), uuid(%s)", id, uuid)
	}

	if exist == 0 {
		return nil
	}

	_, err = r.data.redisCli.HIncrBy(ctx, key, "comment", 1).Result()
	if err != nil {
		r.log.Errorf("fail to add article comment to cache: id(%v), uuid(%s)", id, uuid)
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
	var incrBy = redis.NewScript(`
					local key = KEYS[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
						local number = tonumber(redis.call("HGET", key, "comment"))
						if number > 0 then
  							redis.call("HINCRBY", key, "comment", -1)
						end
					end
					return 0
	`)
	keys := []string{key}
	_, err := incrBy.Run(ctx, r.data.redisCli, keys).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reduce article comment to cache: id(%v), uuid(%s)", id, uuid))
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
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "article_review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.articleReviewMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *articleRepo) SendArticleToMq(ctx context.Context, article *biz.Article, mode string) error {
	articleMap := map[string]interface{}{}
	articleMap["uuid"] = article.Uuid
	articleMap["id"] = article.ArticleId
	articleMap["auth"] = article.Auth
	articleMap["mode"] = mode

	data, err := json.Marshal(articleMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "article",
		Body:  data,
	}
	msg.WithKeys([]string{article.Uuid})
	_, err = r.data.articleMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send article to mq: %v", article))
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

func (r *articleRepo) SetArticleAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "agree", 1)
		pipe.ZIncrBy(ctx, "article_hot", 1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", 1, ids+"%"+uuid+"%article")
		pipe.SAdd(ctx, "article_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add article agree to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
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
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "view", 1)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add article agree to cache: id(%v), uuid(%s)", id, uuid)
	}
	return nil
}

func (r *articleRepo) SetArticleUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error {
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

func (r *articleRepo) SetArticleCollect(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect + ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add article collect: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) SetArticleCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "collect", 1)
		pipe.SAdd(ctx, "article_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to add article collect to cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) CancelArticleAgree(ctx context.Context, id int32, uuid string) error {
	as := ArticleStatistic{}
	err := r.data.DB(ctx).Model(&as).Where("article_id = ? and uuid = ?", id, uuid).Update("agree", gorm.Expr("agree - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article agree: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CancelArticleAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "agree", -1)
		pipe.ZIncrBy(ctx, "article_hot", -1, ids+"%"+uuid)
		pipe.ZIncrBy(ctx, "leaderboard", -1, ids+"%"+uuid+"%article")
		pipe.SRem(ctx, "article_agree_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel article agree from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) CancelArticleUserCollect(ctx context.Context, id int32, userUuid string) error {
	collect := &Collect{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&Collect{}).Where("creations_id = ? and mode = ? and uuid = ?", id, 1, userUuid).Updates(collect).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect: article_id(%v), userUuid(%s)", id, userUuid))
	}
	return nil
}

func (r *articleRepo) CancelArticleCollect(ctx context.Context, id int32, uuid string) error {
	as := &ArticleStatistic{}
	err := r.data.DB(ctx).Model(as).Where("article_id = ? and uuid = ?", id, uuid).Update("collect", gorm.Expr("collect - ?", 1)).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to cancel article collect: id(%v)", id))
	}
	return nil
}

func (r *articleRepo) CancelArticleCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HIncrBy(ctx, "article_"+ids, "collect", -1)
		pipe.SRem(ctx, "article_collect_"+ids, userUuid)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to cancel article collect from cache: id(%v), uuid(%s), userUuid(%s)", id, uuid, userUuid)
	}
	return nil
}

func (r *articleRepo) SendArticleStatisticToMq(ctx context.Context, uuid, mode string) error {
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
		return errors.Wrapf(err, fmt.Sprintf("fail to send article statistic to mq: uuid(%s)", uuid))
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

func (r *articleRepo) getArticleHotFromCache(ctx context.Context, page int32) ([]*biz.ArticleStatistic, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "article_hot", index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article hot from cache: key(%s), page(%v)", "article_hot", page))
	}

	article := make([]*biz.ArticleStatistic, 0)
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

	article := make([]*biz.Article, 0)
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
	err := r.data.db.WithContext(ctx).Where("column_id = ? and status = ?", id, 1).Order("updated_at desc").Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get column article from db: columnId(%v)", id))
	}

	article := make([]*biz.Article, 0)
	for _, item := range list {
		article = append(article, &biz.Article{
			ArticleId: item.ArticleId,
			Uuid:      item.Uuid,
		})
	}
	return article, nil
}

func (r *articleRepo) setArticleToCache(key string, article []*biz.Article) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.ArticleId),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v)", article)
	}
}

func (r *articleRepo) setArticleHotToCache(key string, article []*biz.ArticleStatistic) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(item.Agree),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set article to cache: article(%v)", article)
	}
}

func (r *articleRepo) setColumnArticleToCache(id int32, article []*biz.Article) {
	ids := strconv.Itoa(int(id))
	length := len(article)
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for index, item := range article {
			z = append(z, &redis.Z{
				Score:  float64(length - index),
				Member: strconv.Itoa(int(item.ArticleId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), "column_includes_"+ids, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set column article to cache: article(%v)", article)
	}
}

func (r *articleRepo) getArticleStatisticFromCache(ctx context.Context, key string) (*biz.ArticleStatistic, error) {
	statistic, err := r.data.redisCli.HMGet(ctx, key, "uuid", "agree", "collect", "view", "comment").Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get article statistic form cache: key(%s)", key))
	}
	val := []int32{0, 0, 0, 0}
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
	return &biz.ArticleStatistic{
		Uuid:    statistic[0].(string),
		Agree:   val[0],
		Collect: val[1],
		View:    val[2],
		Comment: val[3],
	}, nil
}

func (r *articleRepo) setArticleStatisticToCache(key string, statistic *biz.ArticleStatistic) {
	err := r.data.redisCli.HMSet(context.Background(), key, "uuid", statistic.Uuid, "agree", statistic.Agree, "collect", statistic.Collect, "view", statistic.View, "comment", statistic.Comment).Err()
	if err != nil {
		r.log.Errorf("fail to set article statistic to cache, err(%s)", err.Error())
	}
}
