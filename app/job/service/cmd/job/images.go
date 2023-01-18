package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func getImages(config map[string]interface{}, m map[string]string, cosCli *CosClient) error {
	url := config["gugudata"].(map[string]interface{})["image_url"].(string)
	key := config["gugudata"].(map[string]interface{})["image_key"].(string)
	limitType := config["gugudata"].(map[string]interface{})["image_limit_type"].(string)
	limitValue := config["gugudata"].(map[string]interface{})["image_limit_value"].(string)
	method := "POST"
	data := map[string]interface{}{}

	payload := strings.NewReader(fmt.Sprintf("appkey=%s&url=%s&limittype=%s&limitvalue=%s&imagewithtag=false&htmlsourcecontent=false", key, m["url"], limitType, limitValue))

	client := &http.Client{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to new request: %s", "getImages"))
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to get news images: %s", "getImages"))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to read data from body: %s", "getImages"))
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to unmarshal news: %s", "getImages"))
	}

	if data["Data"] == nil {
		err = errors.New(fmt.Sprintf("%v", data))
		return errors.Wrapf(err, fmt.Sprintf("fail to get news images: %s", "getImages"))
	}

	images := data["Data"].(map[string]interface{})["ImagesUrl"].([]interface{})
	if len(images) == 0 {
		return nil
	}

	err = getImageAndSendToCos(images[0].(string), m["id"], cosCli)
	if err != nil {
		return err
	}

	m["cover"] = m["id"]
	return nil
}

func getImageAndSendToCos(url string, id string, cosCli *CosClient) error {
	imageUrl := url
	method := "GET"
	client := &http.Client{}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	req, err := http.NewRequestWithContext(ctx, method, imageUrl, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to new an image request: %s", "getImageAndSendToCos"))
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to new an image request: %s", "getImageAndSendToCos"))
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("the status code is not 200: %v", res.StatusCode))
		return errors.Wrapf(err, fmt.Sprintf("fail to get news image: %s", "getImageAndSendToCos"))
	}

	err = sendImageToCos(res.Body, id, cosCli)
	if err != nil {
		return err
	}
	return nil
}

func sendImageToCos(image io.Reader, id string, cosCli *CosClient) error {
	key := "news/" + id + "/cover.png"
	operation := "imageMogr2/format/webp/interlace/0/quality/100"
	pic := &cos.PicOperations{
		IsPicInfo: 1,
		Rules: []cos.PicOperationsRules{
			{
				FileId: "cover.webp",
				Rule:   operation,
			},
		},
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   "image/jpeg",
			XOptionHeader: &http.Header{},
		},
	}
	opt.XOptionHeader.Add("Pic-Operations", cos.EncodePicOperations(pic))
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	_, err := cosCli.client.Object.Put(ctx, key, image, opt)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send news image to cos: %s", "sendImageToCos"))
	}
	return nil
}
