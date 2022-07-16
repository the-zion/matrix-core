package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/creation/service/internal/biz"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var _ biz.TalkRepo = (*talkRepo)(nil)

type talkRepo struct {
	data *Data
	log  *log.Helper
}

func NewTalkRepo(data *Data, logger log.Logger) biz.TalkRepo {
	return &talkRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "creation/data/talk")),
	}
}

func (r *talkRepo) GetTalk(ctx context.Context, id int32) (*biz.Talk, error) {
	talk := &Talk{}
	err := r.data.db.WithContext(ctx).Where("talk_id = ?", id).First(talk).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from db: id(%v)", id))
	}
	return &biz.Talk{
		TalkId: id,
		Uuid:   talk.Uuid,
	}, nil
}

func (r *talkRepo) GetTalkList(ctx context.Context, page int32) ([]*biz.Talk, error) {
	talk, err := r.getTalkFromCache(ctx, page)
	if err != nil {
		return nil, err
	}

	size := len(talk)
	if size != 0 {
		return talk, nil
	}

	talk, err = r.getTalkFromDB(ctx, page)
	if err != nil {
		return nil, err
	}

	size = len(talk)
	if size != 0 {
		go r.setTalkToCache("talk", talk)
	}
	return talk, nil
}

func (r *talkRepo) getTalkFromCache(ctx context.Context, page int32) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int64(page - 1)
	list, err := r.data.redisCli.ZRevRange(ctx, "talk", index*10, index+9).Result()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from cache: key(%s), page(%v)", "talk", page))
	}

	talk := make([]*biz.Talk, 0)
	for _, item := range list {
		member := strings.Split(item, "%")
		id, err := strconv.ParseInt(member[0], 10, 32)
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to covert string to int64: id(%s)", member[0]))
		}
		talk = append(talk, &biz.Talk{
			TalkId: int32(id),
			Uuid:   member[1],
		})
	}
	return talk, nil
}

func (r *talkRepo) getTalkFromDB(ctx context.Context, page int32) ([]*biz.Talk, error) {
	if page < 1 {
		page = 1
	}
	index := int(page - 1)
	list := make([]*Talk, 0)
	err := r.data.db.WithContext(ctx).Where("auth", 1).Order("talk_id desc").Offset(index * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get talk from db: page(%v)", page))
	}

	talk := make([]*biz.Talk, 0)
	for _, item := range list {
		talk = append(talk, &biz.Talk{
			TalkId: item.TalkId,
			Uuid:   item.Uuid,
		})
	}
	return talk, nil
}

func (r *talkRepo) setTalkToCache(key string, talk []*biz.Talk) {
	_, err := r.data.redisCli.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		z := make([]*redis.Z, 0)
		for _, item := range talk {
			z = append(z, &redis.Z{
				Score:  float64(item.TalkId),
				Member: strconv.Itoa(int(item.TalkId)) + "%" + item.Uuid,
			})
		}
		pipe.ZAddNX(context.Background(), key, z...)
		return nil
	})
	if err != nil {
		r.log.Errorf("fail to set talk to cache: talk(%v)", talk)
	}
}

func (r *talkRepo) GetLastTalkDraft(ctx context.Context, uuid string) (*biz.TalkDraft, error) {
	draft := &TalkDraft{}
	err := r.data.db.WithContext(ctx).Where("uuid = ?", uuid).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("talk draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
	}
	return &biz.TalkDraft{
		Id:     int32(draft.ID),
		Status: draft.Status,
	}, nil
}

func (r *talkRepo) CreateTalkDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &TalkDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create a talk draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *talkRepo) CreateTalkFolder(ctx context.Context, id int32, uuid string) error {
	name := "talk/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a talk folder: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SendTalk(ctx context.Context, id int32, uuid string) (*biz.TalkDraft, error) {
	td := &TalkDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&TalkDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 3).Updates(td).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 3: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.TalkDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *talkRepo) SendReviewToMq(ctx context.Context, review *biz.TalkReview) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk_review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.talkReviewMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}

func (r *talkRepo) DeleteTalkDraft(ctx context.Context, id int32, uuid string) error {
	td := &TalkDraft{}
	td.ID = uint(id)
	err := r.data.DB(ctx).Where("uuid = ?", uuid).Delete(td).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to delete talk draft: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) CreateTalk(ctx context.Context, id, auth int32, uuid string) error {
	talk := &Talk{
		TalkId: id,
		Uuid:   uuid,
		Auth:   auth,
	}
	err := r.data.DB(ctx).Select("TalkId", "Uuid", "Auth").Create(talk).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a talk: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) CreateTalkStatistic(ctx context.Context, id, auth int32, uuid string) error {
	ts := &TalkStatistic{
		TalkId: id,
		Uuid:   uuid,
		Auth:   auth,
	}
	err := r.data.DB(ctx).Select("TalkId", "Uuid", "Auth").Create(ts).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a talk statistic: id(%v)", id))
	}
	return nil
}

func (r *talkRepo) SendTalkToMq(ctx context.Context, talk *biz.Talk, mode string) error {
	talkMap := map[string]interface{}{}
	talkMap["uuid"] = talk.Uuid
	talkMap["id"] = talk.TalkId
	talkMap["auth"] = talk.Auth
	talkMap["mode"] = mode

	data, err := json.Marshal(talkMap)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "talk",
		Body:  data,
	}
	msg.WithKeys([]string{talk.Uuid})
	_, err = r.data.talkMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send talk to mq: %v", talk))
	}
	return nil
}

func (r *talkRepo) CreateTalkCache(ctx context.Context, id, auth int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	_, err := r.data.redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSetNX(ctx, "talk_"+ids, "uuid", uuid)
		pipe.HSetNX(ctx, "talk_"+ids, "agree", 0)
		pipe.HSetNX(ctx, "talk_"+ids, "collect", 0)
		pipe.HSetNX(ctx, "talk_"+ids, "view", 0)
		pipe.HSetNX(ctx, "talk_"+ids, "comment", 0)

		if auth == 2 {
			return nil
		}

		pipe.ZAddNX(ctx, "talk", &redis.Z{
			Score:  float64(id),
			Member: ids + "%" + uuid,
		})
		pipe.ZAddNX(ctx, "talk_hot", &redis.Z{
			Score:  0,
			Member: ids + "%" + uuid,
		})
		pipe.ZAddNX(ctx, "leaderboard", &redis.Z{
			Score:  0,
			Member: ids + "%" + uuid + "%talk",
		})
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create talk cache: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) EditTalkCos(ctx context.Context, id int32, uuid string) error {
	err := r.EditTalkCosContent(ctx, id, uuid)
	if err != nil {
		return err
	}

	err = r.EditTalkCosIntroduce(ctx, id, uuid)
	if err != nil {
		return err
	}
	return nil
}

func (r *talkRepo) EditTalkCosContent(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "talk/" + uuid + "/" + ids + "/content-edit"
	dest := "talk/" + uuid + "/" + ids + "/content"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy talk from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) EditTalkCosIntroduce(ctx context.Context, id int32, uuid string) error {
	ids := strconv.Itoa(int(id))
	name := "talk/" + uuid + "/" + ids + "/introduce-edit"
	dest := "talk/" + uuid + "/" + ids + "/introduce"
	sourceURL := fmt.Sprintf("%s/%s", r.data.cosCli.BaseURL.BucketURL.Host, name)
	_, _, err := r.data.cosCli.Object.Copy(ctx, dest, sourceURL, nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to copy talk from edit to content: uuid(%s), id(%v)", uuid, id))
	}
	return nil
}

func (r *talkRepo) CreateTalkSearch(ctx context.Context, id int32, uuid string) error {
	return nil
}

func (r *talkRepo) EditTalkSearch(ctx context.Context, id int32, uuid string) error {
	return nil
}
