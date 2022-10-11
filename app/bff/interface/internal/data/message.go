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

func (r *messageRepo) GetMessageNotification(ctx context.Context, uuid string, follows []string) (*biz.Notification, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_message_notification_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.mc.GetMessageNotification(ctx, &messageV1.GetMessageNotificationReq{
			Uuid:    uuid,
			Follows: follows,
		})
		if err != nil {
			return nil, err
		}
		return &biz.Notification{
			Timeline:           reply.Timeline,
			Comment:            reply.Comment,
			SubComment:         reply.SubComment,
			SystemNotification: reply.System,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.Notification), nil
}

func (r *messageRepo) GetMessageSystemNotification(ctx context.Context, page int32, uuid string) ([]*biz.SystemNotification, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_message_system_notification_%v_%v", page, uuid), func() (interface{}, error) {
		reply := make([]*biz.SystemNotification, 0)
		notificationList, err := r.data.mc.GetMessageSystemNotification(ctx, &messageV1.GetMessageSystemNotificationReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		for _, item := range notificationList.List {
			reply = append(reply, &biz.SystemNotification{
				Id:               item.Id,
				ContentId:        item.ContentId,
				CreatedAt:        item.CreatedAt,
				NotificationType: item.NotificationType,
				Title:            item.Title,
				Uid:              item.Uid,
				Uuid:             item.Uuid,
				Label:            item.Label,
				Result:           item.Result,
				Section:          item.Section,
				Text:             item.Text,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.SystemNotification), nil
}

func (r *messageRepo) SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error {
	_, err := r.data.mc.SetMailBoxLastTime(ctx, &messageV1.SetMailBoxLastTimeReq{
		Uuid: uuid,
		Time: time,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepo) RemoveMailBoxCommentCount(ctx context.Context, uuid string) error {
	_, err := r.data.mc.RemoveMailBoxCommentCount(ctx, &messageV1.RemoveMailBoxCommentCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepo) RemoveMailBoxSubCommentCount(ctx context.Context, uuid string) error {
	_, err := r.data.mc.RemoveMailBoxSubCommentCount(ctx, &messageV1.RemoveMailBoxSubCommentCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *messageRepo) RemoveMailBoxSystemNotificationCount(ctx context.Context, uuid string) error {
	_, err := r.data.mc.RemoveMailBoxSystemNotificationCount(ctx, &messageV1.RemoveMailBoxSystemNotificationCountReq{
		Uuid: uuid,
	})
	if err != nil {
		return err
	}
	return nil
}
