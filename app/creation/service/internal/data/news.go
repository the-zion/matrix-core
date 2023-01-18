package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"strconv"
	"time"
)

var _ biz.NewsRepo = (*newRepo)(nil)

type newRepo struct {
	data *Data
	log  *log.Helper
}

func NewNewsRepo(data *Data, logger log.Logger) biz.NewsRepo {
	return &newRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/new")),
	}
}

func (r *newRepo) GetNews(ctx context.Context, page int32) ([]*biz.News, error) {
	news, err := r.getNewsFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(news)
	if size != 0 {
		return news, nil
	}

	news, err = r.getNewsFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(news)
	if size != 0 {
		newCtx, _ := context.WithTimeout(context.Background(), time.Second*2)
		go r.data.Recover(newCtx, func(ctx context.Context) {
			r.setArticleToCache("news", news)
		})()
	}
	return news, nil
}

func (r *newRepo) getNewsFromCache(ctx context.Context, page int32) ([]*biz.News, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "news", index*10, index*10+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get news from cache: key(%s), page(%v)", "news", page))
	}

	news := make([]*biz.News, 0, len(list))
	for _, item := range list {
		news = append(news, &biz.News{
			Id: item,
		})
	}
	return news, nil
}

func (r *newRepo) getNewsFromDB(ctx context.Context, page int32) ([]*biz.News, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*News, 0)
	err := r.data.db.WithContext(ctx).Select("id").Order("article_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get news from db: page(%v)", page))
	}

	news := make([]*biz.News, 0, len(list))
	for _, item := range list {
		news = append(news, &biz.News{
			Id: item.Id,
		})
	}
	return news, nil
}

func (r *newRepo) setArticleToCache(key string, news []*biz.News) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0, len(news))
		for _, item := range news {
			id, err := strconv.Atoi(item.Id)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("cannot convert id %s from string to int", item.Id))
			}
			z = append(z, &redis.Z{
				Score:  float64(id),
				Member: item.Id,
			})
		}
		pipe.ZAddNX(ctx, key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set news to cache: news(%v), err(%v)", news, err)
	}
}

func (r *newRepo) GetNewsSearch(ctx context.Context, page int32, search, time string) ([]*biz.NewsSearch, int32, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	reply := make([]*biz.NewsSearch, 0)

	var buf bytes.Buffer
	var body map[string]interface{}
	if search != "" {
		body = map[string]interface{}{
			"from":    index * 10,
			"size":    10,
			"_source": []string{"update", "tags", "author", "url", "cover"},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"multi_match": map[string]interface{}{
							"query":  search,
							"fields": []string{"tags", "title", "content"},
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
					"content": map[string]interface{}{
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
			"_source": []string{"update", "tags", "author", "url", "cover"},
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
					"content": map[string]interface{}{
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
		r.data.elasticSearch.es.Search.WithIndex("news"),
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
			return nil, 0, errors.Errorf(fmt.Sprintf("error search news from  es: reason(%v), page(%v), search(%s), time(%s)", e, page, search, time))
		}
	}

	result := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, errors.Wrapf(err, fmt.Sprintf("error parsing the response body: page(%v), search(%s), time(%s)", page, search, time))
	}

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		news := &biz.NewsSearch{
			Id:     hit.(map[string]interface{})["_id"].(string),
			Tags:   hit.(map[string]interface{})["_source"].(map[string]interface{})["tags"].(string),
			Update: hit.(map[string]interface{})["_source"].(map[string]interface{})["update"].(string),
			Author: hit.(map[string]interface{})["_source"].(map[string]interface{})["author"].(string),
			Url:    hit.(map[string]interface{})["_source"].(map[string]interface{})["url"].(string),
			Cover:  hit.(map[string]interface{})["_source"].(map[string]interface{})["cover"].(string),
		}

		if content, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["content"]; ok {
			news.Content = content.([]interface{})[0].(string)
		}

		if title, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["title"]; ok {
			news.Title = title.([]interface{})[0].(string)
		}

		reply = append(reply, news)
	}
	res.Body.Close()
	return reply, int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), nil
}
