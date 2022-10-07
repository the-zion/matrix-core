package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"strconv"
	"time"
)

type messageRepo struct {
	data *Data
	log  *log.Helper
}

func NewMessageRepo(data *Data, logger log.Logger) biz.MessageRepo {
	return &messageRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "message/data/message")),
	}
}

func (r *messageRepo) GetMailBoxLastTime(ctx context.Context, uuid string) (*biz.MailBox, error) {
	timeValue := "0"
	timeValue, err := r.data.redisCli.HGet(ctx, "mailbox_last_time", uuid).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		r.log.Errorf("fail to get mailbox last time: uuid(%s)", uuid)
	}
	num, err := strconv.ParseInt(timeValue, 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: time(%s)", timeValue))
	}
	return &biz.MailBox{
		Time: int32(num),
	}, nil
}

func (r *messageRepo) SetMailBoxLastTime(ctx context.Context, uuid string) error {
	err := r.data.redisCli.HSet(ctx, "mailbox_last_time", uuid, float64(time.Now().Unix()))
	if err != nil {
		r.log.Errorf("fail to set mailbox last time: uuid(%s)", uuid)
	}
	return nil
}
