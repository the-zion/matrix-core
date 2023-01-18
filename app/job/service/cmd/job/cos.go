package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
)

type CosClient struct {
	client *cos.Client
}

func (s *CosClient) setNews(m map[string]string) error {
	content, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news content: %s", "cos-setNews"))
	}

	name := "news/" + m["id"] + "/content"
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	_, err = s.client.Object.Put(ctx, name, bytes.NewReader(content), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add content to cos: %s", "cos-setNews"))
	}

	introduce, err := json.Marshal(map[string]string{
		"id":        m["id"],
		"cover":     m["cover"],
		"update":    m["update"],
		"title":     m["title"],
		"author":    m["author"],
		"tags":      m["tags"],
		"url":       m["url"],
		"introduce": m["introduce"],
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news content: %s", "cos-setNews"))
	}

	name = "news/" + m["id"] + "/introduce"
	ctx, _ = context.WithTimeout(context.Background(), time.Second*60)
	_, err = s.client.Object.Put(ctx, name, bytes.NewReader(introduce), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add introduce to cos: %s", "cos-setNews"))
	}
	return err
}

func newCosServiceClient(config map[string]interface{}) (*CosClient, error) {
	addr := config["cos"].(map[string]interface{})["url"].(string)
	secretId := config["cos"].(map[string]interface{})["secret_id"].(string)
	secretKey := config["cos"].(map[string]interface{})["secret_key"].(string)

	u, err := url.Parse(addr)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to init cos server: %s", "newCosServiceClient"))
	}
	b := &cos.BaseURL{BucketURL: u}
	return &CosClient{
		client: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  secretId,
				SecretKey: secretKey,
			},
		})}, nil
}
