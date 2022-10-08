package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/message/service/internal/biz"
	"strconv"
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
	timeValue, err := r.data.redisCli.HGet(ctx, "mailbox_last_time", uuid).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, "fail to get mailbox last time: uuid(%s)", uuid)
	}

	if timeValue == "" {
		timeValue = "0"
	}

	num, err := strconv.ParseInt(timeValue, 10, 32)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: time(%s)", timeValue))
	}
	return &biz.MailBox{
		Time: int32(num),
	}, nil
}

func (r *messageRepo) GetMessageNotification(ctx context.Context, uuid string, follows []string) (*biz.Notification, error) {
	result, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMGet(context.Background(), "message_timeline", follows...)
		pipe.HMGet(context.Background(), "message_comment", uuid)
		pipe.HMGet(context.Background(), "message_sub_comment", uuid)
		pipe.HMGet(context.Background(), "message_system", uuid)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get message notification from cache, err(%s)", err.Error())
	}

	timeline := make(map[string]int32, 0)
	for index, times := range result[0].(*redis.SliceCmd).Val() {
		var t int64
		if times != nil {
			t, err = strconv.ParseInt(times.(string), 10, 32)
			if err != nil {
				return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: times(%v)", times))
			}
		}
		timeline[follows[index]] = int32(t)
	}

	var comment int32
	var c int64
	comments := result[1].(*redis.SliceCmd).Val()[0]
	if comments != nil {
		c, err = strconv.ParseInt(comments.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: comments(%v)", comments))
		}
	}
	comment = int32(c)

	var subComment int32
	var sc int64
	subComments := result[2].(*redis.SliceCmd).Val()[0]
	if subComments != nil {
		sc, err = strconv.ParseInt(subComments.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: subComments(%v)", subComments))
		}
	}
	subComment = int32(sc)

	var system int32
	var s int64
	systems := result[3].(*redis.SliceCmd).Val()[0]
	if systems != nil {
		c, err = strconv.ParseInt(systems.(string), 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: systems(%v)", systems))
		}
	}
	system = int32(s)

	return &biz.Notification{
		Timeline:           timeline,
		Comment:            comment,
		SubComment:         subComment,
		SystemNotification: system,
	}, nil
}

func (r *messageRepo) SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error {
	_, err := r.data.redisCli.HSet(ctx, "mailbox_last_time", uuid, time).Result()
	if err != nil {
		return errors.Wrapf(err, "fail to set mailbox last time: uuid(%s)", uuid)
	}
	return nil
}
