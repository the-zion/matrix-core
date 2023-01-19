package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var (
	addr      string
	password  string
	scriptBox = map[string]string{
		"AddMailBoxSystemNotificationToCache": `
					local key = KEYS[1]
					local value = ARGV[1]
					local uuid = ARGV[2]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					redis.call("HINCRBY", "message_system", uuid, 1)
					return 0
	`,
	}
)

func NewRedis(logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "lua", "redis"))

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          4,
		DialTimeout: time.Second * 2,
		PoolSize:    10,
		Password:    password,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		l.Fatalf("redis connect error: %v", err)
	}
	return client
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)
	r := NewRedis(logger)
	ScriptLoad(r)
}

func ScriptLoad(r redis.Cmdable) {
	for key, value := range scriptBox {
		result, err := r.ScriptLoad(context.Background(), value).Result()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(key, result)
	}
}

func init() {
	flag.StringVar(&addr, "addr",
		"ip:port",
		"redis addr, eg: -addr ip:port")
	flag.StringVar(&password, "password",
		"abc",
		"redis password, eg: -password abc")
}
