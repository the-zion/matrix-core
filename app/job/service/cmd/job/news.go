package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type News struct {
	Id      string
	Update  string
	Title   string
	Author  string
	Content string
	Tags    string
	Cover   string
	Url     string
}

func getNews(config map[string]interface{}) ([]*News, error) {
	addr := config["gugudata"].(map[string]interface{})["url"].(string)

	data := map[string]interface{}{}
	url := addr + "1"
	method := "GET"
	client := &http.Client{}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new request: %s", "getNews"))
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get news: %s", "getNews"))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to read data from body: %s", "getNews"))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news: %s", "getNews"))
	}

	news := make([]*News, 0, len(data["Data"].([]interface{})))
	for _, item := range data["Data"].([]interface{}) {
		title := regexp.MustCompile("\\s+").ReplaceAllString(item.(map[string]interface{})["ArticleTitle"].(string), "")
		title = regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]").ReplaceAllString(title, "")
		news = append(news, &News{
			Id:     item.(map[string]interface{})["ArticleId"].(string),
			Update: item.(map[string]interface{})["CreateDateTime"].(string),
			Title:  title,
			Author: item.(map[string]interface{})["ArticleAuthor"].(string),
			Tags:   getTags(item.(map[string]interface{})["Tags"].([]interface{})),
			Url:    item.(map[string]interface{})["ArticleSourceUrl"].(string),
		})
	}
	return news, nil
}

func getTags(tags []interface{}) string {
	box := make([]string, 0, len(tags))
	for _, item := range tags {
		box = append(box, item.(map[string]interface{})["TagName"].(string))
	}
	return strings.Join(box, `;`)
}
