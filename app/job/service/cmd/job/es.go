package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"time"
)

type ElasticSearch struct {
	client *elasticsearch.Client
}

func (es *ElasticSearch) setNews(m map[string]string) error {
	delete(m, "introduce")
	s, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to marshal news: %s", "es-setNews"))
	}

	req := esapi.IndexRequest{
		Index:      "news",
		DocumentID: m["id"],
		Body:       bytes.NewReader(s),
		Refresh:    "true",
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	res, err := req.Do(ctx, es.client)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("error getting news search create response: %s", "es-setNews"))
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
			return errors.Wrapf(err, fmt.Sprintf("error parsing the response body: %s", "es-setNews"))
		} else {
			return errors.Wrapf(err, fmt.Sprintf("error indexing document to es: %s", "es-setNews"))
		}
	}
	return nil
}

func newElasticsearch(config map[string]interface{}) (*ElasticSearch, error) {
	addr := config["elasticSearch"].(map[string]interface{})["addr"].(string)
	user := config["elasticSearch"].(map[string]interface{})["user"].(string)
	password := config["elasticSearch"].(map[string]interface{})["password"].(string)

	cfg := elasticsearch.Config{
		Username: user,
		Password: password,
		Addresses: []string{
			addr,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("error creating the es client: %s", "newElasticsearch"))
	}

	res, err := es.Info()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("Error getting response: %s", "newElasticsearch"))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.Wrapf(err, fmt.Sprintf("Error: %s", "newElasticsearch"))
	}

	return &ElasticSearch{
		client: es,
	}, nil
}
