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
	"unicode/utf8"
)

func getContent(config map[string]interface{}, m map[string]string) error {
	url := config["gugudata"].(map[string]interface{})["content_url"].(string)
	key := config["gugudata"].(map[string]interface{})["content_key"].(string)
	method := "POST"
	data := map[string]interface{}{}

	payload := strings.NewReader(fmt.Sprintf("appkey=%s&url=%s&contentwithhtml=false&htmlsourcecontent=false", key, m["url"]))

	client := &http.Client{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to new request: %s", "getContent"))
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get news content: %s", "getContent"))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read data from body: %s", "getContent"))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news: %s", "getContent"))
	}

	if data["Data"] == nil {
		err = errors.New(fmt.Sprintf("%v", data))
		return errors.Wrapf(err, fmt.Sprintf("fail to get news content: %s", "getContent"))
	}

	m["content"] = regexp.MustCompile("\\s+").ReplaceAllString(data["Data"].(map[string]interface{})["Content"].(string), "")
	m["content"] = regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]").ReplaceAllString(m["content"], "")
	if m["content"] != "" {
		length := utf8.RuneCountInString(m["content"])
		if length > 250 {
			var n, i int
			for i = range m["content"] {
				if n == 250 {
					break
				}
				n++
			}
			m["introduce"] = m["content"][:i]
		} else {
			m["introduce"] = m["content"]
		}
	}
	return nil
}
