package main

import (
	"fmt"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"os"
	"time"
)

type Log struct {
	client  *cls.AsyncProducerClient
	topicId string
}

func (l *Log) SendLog(content string) {
	callBack := &Callback{}
	log := cls.NewCLSLog(time.Now().Unix(), map[string]string{"error": content})
	err := l.client.SendLog(l.topicId, log, callBack)
	if err != nil {
		fmt.Println(err)
	}
}

type Callback struct {
}

func (callback *Callback) Success(result *cls.Result) {
	attemptList := result.GetReservedAttempts()
	for _, attempt := range attemptList {
		fmt.Printf("%+v \n", attempt)
	}
}

func (callback *Callback) Fail(result *cls.Result) {
	fmt.Println(result.IsSuccessful())
	fmt.Println(result.GetErrorCode())
	fmt.Println(result.GetErrorMessage())
	fmt.Println(result.GetReservedAttempts())
	fmt.Println(result.GetRequestId())
	fmt.Println(result.GetTimeStampMs())
}

func logInit() (*Log, error) {
	producerConfig := cls.GetDefaultAsyncProducerClientConfig()
	producerConfig.Endpoint = os.Getenv("TENCENT_LOG_HOST")
	producerConfig.AccessKeyID = os.Getenv("TENCENT_LOG_ACCESSKEY")
	producerConfig.AccessKeySecret = os.Getenv("TENCENT_LOG_ACCESSSECRET")
	topicId := os.Getenv("TENCENT_LOG_TOPIC_ID")
	producerInstance, err := cls.NewAsyncProducerClient(producerConfig)
	if err != nil {
		fmt.Print(fmt.Sprintf("fail to init log: %v", err))
		return nil, err
	}

	producerInstance.Start()
	return &Log{
		client:  producerInstance,
		topicId: topicId,
	}, nil
}
