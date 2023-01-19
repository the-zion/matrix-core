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
		"SetUserMedalToCache": `
					local key = KEYS[1]
                    local change = ARGV[1]
					local value = redis.call("EXISTS", key)
					if value == 1 then
  						redis.call("HSET", key, change, 1)
					end
					return 0
	`,
		"SetAchievementAgreeToCache": `
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "agree", 1)
					end
					return 0
	`,
		"CancelAchievementAgreeFromCache": `
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", uuid, "agree"))
						if number > 0 then
  							redis.call("HINCRBY", uuid, "agree", -1)
						end
					end
					return 0
	`,
		"SetAchievementViewToCache": `
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "view", 1)
					end
					return 0
	`,
		"SetAchievementCollectToCache": `
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "collect", 1)
					end
					return 0
	`,
		"CancelAchievementCollectFromCache": `
					local uuid = KEYS[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", uuid, "collect"))
						if number > 0 then
  							redis.call("HINCRBY", uuid, "collect", -1)
						end
					end
					return 0
	`,
		"SetAchievementFollowToCache": `
					local follow = KEYS[1]
					local exist = redis.call("EXISTS", follow)
					if exist == 1 then
						redis.call("HINCRBY", follow, "followed", 1)
					end

					local followed = KEYS[2]
					local exist = redis.call("EXISTS", followed)
					if exist == 1 then
						redis.call("HINCRBY", followed, "follow", 1)
					end
					return 0
	`,
		"CancelAchievementFollowFromCache": `
					local follow = KEYS[1]
					local exist = redis.call("EXISTS", follow)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", follow, "followed"))
						if number > 0 then
  							redis.call("HINCRBY", follow, "followed", -1)
						end
					end

					local followed = KEYS[2]
					local exist = redis.call("EXISTS", followed)
					if exist == 1 then
						local number = tonumber(redis.call("HGET", followed, "follow"))
						if number > 0 then
  							redis.call("HINCRBY", followed, "follow", -1)
						end
					end
					return 0
	`,
		"CancelUserMedalFromCache": `
					local key = KEYS[1]
					local change = ARGV[1]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("HSET", key, change, 2)
					end
					return 0
	`,
		"AddAchievementScoreToCache": `
					local uuid = KEYS[1]
					local value = ARGV[1]
					local exist = redis.call("EXISTS", uuid)
					if exist == 1 then
						redis.call("HINCRBY", uuid, "score", value)
					end
					return 0
	`,
	}
)

func NewRedis(logger log.Logger) redis.Cmdable {
	l := log.NewHelper(log.With(logger, "lua", "redis"))

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          2,
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
