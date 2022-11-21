package data

import (
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
		fmt.Println(err)
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
