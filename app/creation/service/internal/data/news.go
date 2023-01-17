package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

func (r *newRepo) GetNews(_ context.Context, page int32) ([]*biz.News, error) {
	data := map[string]interface{}{}
	url := r.data.newsCli.url + strconv.Itoa(int(page))
	method := "GET"
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new a request: page(%v)", page))
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get news: page(%v)", page))
	}

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news: page(%v)", page))
	}

	news := make([]*biz.News, 0, len(data["Data"].([]interface{})))
	for _, item := range data["Data"].([]interface{}) {
		news = append(news, &biz.News{
			Id:     item.(map[string]interface{})["ArticleId"].(string),
			Update: item.(map[string]interface{})["CreateDateTime"].(string),
			Title:  item.(map[string]interface{})["ArticleTitle"].(string),
			Author: item.(map[string]interface{})["ArticleAuthor"].(string),
			Tags:   r.getTags(item.(map[string]interface{})["Tags"].([]interface{})),
			Url:    item.(map[string]interface{})["ArticleSourceUrl"].(string),
		})
	}
	return news, nil
}

func (r *newRepo) getTags(tags []interface{}) string {
	box := make([]string, 0, len(tags))
	for _, item := range tags {
		box = append(box, item.(map[string]interface{})["TagName"].(string))
	}
	return strings.Join(box, `;`)
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
			"_source": []string{"update", "tags", "author", "url"},
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{"multi_match": map[string]interface{}{
							"query":  search,
							"fields": []string{"tags", "title"},
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
			"_source": []string{"update", "tags", "author", "url"},
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
		}

		if title, ok := hit.(map[string]interface{})["highlight"].(map[string]interface{})["title"]; ok {
			news.Title = title.([]interface{})[0].(string)
		}

		reply = append(reply, news)
	}
	res.Body.Close()
	return reply, int32(result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), nil
}
