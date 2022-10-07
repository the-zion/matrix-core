package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	messageV1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.MessageRepo = (*messageRepo)(nil)

type messageRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewMessageRepo(data *Data, logger log.Logger) biz.MessageRepo {
	return &messageRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "bff/data/message")),
		sg:   &singleflight.Group{},
	}
}

func (r *messageRepo) GetMailBoxLastTime(ctx context.Context, uuid string) (*biz.MailBox, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_mailbox_last_time_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.mc.GetMailBoxLastTime(ctx, &messageV1.GetMailBoxLastTimeReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.MailBox{
			Time: reply.Time,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.MailBox), nil
}

func (r *messageRepo) SetMailBoxLastTime(ctx context.Context, uuid string) error {
	_, err := r.data.mc.SetMailBoxLastTime(ctx, &messageV1.SetMailBoxLastTimeReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}
