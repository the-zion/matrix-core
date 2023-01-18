package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func getConfig() (map[string]interface{}, error) {
	config := map[string]interface{}{}
	url := fmt.Sprintf("http://%s/nacos/v1/cs/configs?dataId=%s&group=%s&tenant=%s&accessToken=%s",
		os.Getenv("NACOS_ADDR"),
		os.Getenv("NACOS_DATAID"),
		os.Getenv("NACOS_GROUP"),
		os.Getenv("NACOS_NAMESPACE"),
		os.Getenv("NACOS_ACCESS_TOKEN"))
	method := "GET"
	client := &http.Client{}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to new request: %s", "getConfig"))
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to load config: %s", "getConfig"))
	}

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news: %s", "getConfig"))
	}
	return config, nil
}
