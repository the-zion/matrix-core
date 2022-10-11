package data

import (
	"context"
	"encoding/json"
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

func (r *messageRepo) GetMessageSystemNotification(ctx context.Context, page int32, uuid string) ([]*biz.SystemNotification, error) {
	key := "system_notification_" + uuid
	review, err := r.getMessageSystemNotificationFromCache(ctx, page, key)
	if err != nil {
		return nil, err
	}

	size := len(review)
	if size != 0 {
		return review, nil
	}

	review, err = r.getMessageSystemNotificationFromDB(ctx, page, uuid)
	if err != nil {
		return nil, err
	}

	size = len(review)
	if size != 0 {
		go r.data.Recover(context.Background(), func(ctx context.Context) {
			r.setMessageSystemNotificationToCache(key, review)
		})()
	}
	return review, nil
}

func (r *messageRepo) getMessageSystemNotificationFromCache(ctx context.Context, page int32, key string) ([]*biz.SystemNotification, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.LRange(ctx, key, index*20, index*20+19).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get message system notification from cache: key(%s), page(%v)", key, page))
	}

	notificationList := make([]*biz.SystemNotification, 0)
	for _, item := range list {
		var notification = &biz.SystemNotification{}
		err = json.Unmarshal([]byte(item), notification)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: notification(%v)", item))
		}
		notificationList = append(notificationList, &biz.SystemNotification{
			Id:               notification.Id,
			ContentId:        notification.ContentId,
			CreatedAt:        notification.CreatedAt,
			NotificationType: notification.NotificationType,
			Title:            notification.Title,
			Uid:              notification.Uid,
			Uuid:             notification.Uuid,
			Label:            notification.Label,
			Result:           notification.Result,
			Section:          notification.Section,
			Text:             notification.Text,
		})
	}
	return notificationList, nil
}

func (r *messageRepo) getMessageSystemNotificationFromDB(ctx context.Context, page int32, uuid string) ([]*biz.SystemNotification, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*SystemNotification, 0)
	err := r.data.db.WithContext(ctx).Where("uuid", uuid).Order("id desc").Offset(index * 20).Limit(20).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get message system notification from db: page(%v), uuid(%s)", page, uuid))
	}

	notification := make([]*biz.SystemNotification, 0)
	for _, item := range list {
		notification = append(notification, &biz.SystemNotification{
			Id:               int32(item.ID),
			ContentId:        item.ContentId,
			CreatedAt:        item.CreatedAt.Format("2006-01-02"),
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
	return notification, nil
}

func (r *messageRepo) setMessageSystemNotificationToCache(key string, notification []*biz.SystemNotification) {
	ctx := context.Background()
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		list := make([]interface{}, 0)
		for _, item := range notification {
			m, err := json.Marshal(item)
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("fail to marshal avatar review: notification(%v)", item))
			}
			list = append(list, m)
		}
		pipe.RPush(ctx, key, list...)
		pipe.Expire(ctx, key, time.Hour*8)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set message system notification to cache: notifications(%v), err(%v)", notification, err)
	}
}

func (r *messageRepo) SetMailBoxLastTime(ctx context.Context, uuid string, time int32) error {
	_, err := r.data.redisCli.HSet(ctx, "mailbox_last_time", uuid, time).Result()
	if err != nil {
		return errors.Wrapf(err, "fail to set mailbox last time: uuid(%s)", uuid)
	}
	return nil
}

func (r *messageRepo) RemoveMailBoxCommentCount(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HSet(ctx, "message_comment", uuid, 0).Result()
	if err != nil {
		return errors.Wrapf(err, "fail to remove mailbox comment count: uuid(%s)", uuid)
	}
	return nil
}

func (r *messageRepo) RemoveMailBoxSubCommentCount(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HSet(ctx, "message_sub_comment", uuid, 0).Result()
	if err != nil {
		return errors.Wrapf(err, "fail to remove mailbox sub comment count: uuid(%s)", uuid)
	}
	return nil
}

func (r *messageRepo) RemoveMailBoxSystemNotificationCount(ctx context.Context, uuid string) error {
	_, err := r.data.redisCli.HSet(ctx, "message_system", uuid, 0).Result()
	if err != nil {
		return errors.Wrapf(err, "fail to remove mailbox system notification count: uuid(%s)", uuid)
	}
	return nil
}

func (r *messageRepo) AddMailBoxSystemNotification(ctx context.Context, contentId int32, notificationType string, title string, uuid string, label string, result int32, section string, text string, uid string) (*biz.SystemNotification, error) {
	ar := &SystemNotification{
		ContentId:        contentId,
		NotificationType: notificationType,
		Title:            title,
		Uid:              uid,
		Uuid:             uuid,
		Label:            label,
		Result:           result,
		Section:          section,
		Text:             text,
	}
	err := r.data.db.WithContext(ctx).Select("ContentId", "NotificationType", "Title", "Uuid", "Uid", "Label", "Result", "Section", "Text").Create(ar).Error
	if err != nil {
		return nil, errors.Wrapf(err, "fail to add mail box notification: notification(%v)", ar)
	}
	return &biz.SystemNotification{
		Id:               int32(ar.ID),
		CreatedAt:        time.Now().Format("2006-01-02"),
		ContentId:        contentId,
		NotificationType: notificationType,
		Title:            title,
		Uuid:             uuid,
		Label:            label,
		Result:           result,
		Section:          section,
		Text:             text,
	}, nil
}

func (r *messageRepo) AddMailBoxSystemNotificationToCache(ctx context.Context, notification *biz.SystemNotification) error {
	marshal, err := json.Marshal(notification)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set system notification to json: json.Marshal(%v)", notification))
	}
	var script = redis.NewScript(`
					local key = KEYS[1]
					local value = ARGV[1]
					local uuid = ARGV[2]
					local exist = redis.call("EXISTS", key)
					if exist == 1 then
						redis.call("LPUSH", key, value)
					end
					redis.call("HINCRBY", "message_system", uuid, 1)
					return 0
	`)
	keys := []string{"system_notification_" + notification.Uuid}
	values := []interface{}{marshal, notification.Uuid}
	_, err = script.Run(ctx, r.data.redisCli, keys, values...).Result()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set system notification to json: notification(%v)", notification))
	}
	return nil
}
