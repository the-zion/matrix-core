package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Redis struct {
	client redis.Cmdable
}

func (r *Redis) setNews(m map[string]string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	id, err := strconv.Atoi(m["id"])
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("cannot convert id %s from string to int: %s", m["id"], "redis-setNews"))
	}

	err = r.client.ZAddNX(ctx, "news", &redis.Z{
		Score:  float64(id),
		Member: m["id"],
	}).Err()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("cannot set news %s to redis: %s", m["id"], "redis-setNews"))
	}
	return nil
}

func newRedis(config map[string]interface{}) (*Redis, error) {
	addr := config["redis"].(map[string]interface{})["addr"].(string)
	password := config["redis"].(map[string]interface{})["password"].(string)

	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		DB:          1,
		DialTimeout: time.Second * 2,
		PoolSize:    10,
		Password:    password,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("redis connect: %s", "newRedis"))
	}
	return &Redis{
		client: client,
	}, nil
}
