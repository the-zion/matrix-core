package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/comment/service/internal/biz"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var _ biz.CommentRepo = (*commentRepo)(nil)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "comment/data/comment")),
	}
}

func (r *commentRepo) GetLastCommentDraft(ctx context.Context, uuid string) (*biz.CommentDraft, error) {
	draft := &CommentDraft{}
	err := r.data.db.WithContext(ctx).Where("uuid = ? and status = ?", uuid, 1).Last(draft).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("comment draft not found from db", fmt.Sprintf("uuid(%s)", uuid))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: uuid(%s)", uuid))
	}
	return &biz.CommentDraft{
		Id: int32(draft.ID),
	}, nil
}

func (r *commentRepo) CreateCommentDraft(ctx context.Context, uuid string) (int32, error) {
	draft := &CommentDraft{
		Uuid: uuid,
	}
	err := r.data.DB(ctx).Select("Uuid").Create(draft).Error
	if err != nil {
		return 0, errors.Wrapf(err, fmt.Sprintf("fail to create a comment draft: uuid(%s)", uuid))
	}
	return int32(draft.ID), nil
}

func (r *commentRepo) CreateCommentFolder(ctx context.Context, id int32, uuid string) error {
	name := "comment/" + uuid + "/" + strconv.Itoa(int(id)) + "/"
	_, err := r.data.cosCli.Object.Put(ctx, name, strings.NewReader(""), nil)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to create a comment folder: id(%v)", id))
	}
	return nil
}

func (r *commentRepo) SendComment(ctx context.Context, id int32, uuid string) (*biz.CommentDraft, error) {
	cd := &CommentDraft{
		Status: 2,
	}
	err := r.data.DB(ctx).Model(&CommentDraft{}).Where("id = ? and uuid = ? and status = ?", id, uuid, 1).Updates(cd).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to mark draft to 2: uuid(%s), id(%v)", uuid, id))
	}
	return &biz.CommentDraft{
		Uuid: uuid,
		Id:   id,
	}, nil
}

func (r *commentRepo) SetRecord(ctx context.Context, id int32, uuid, ip string) error {
	record := &Record{
		CommonId: id,
		Uuid:     uuid,
		Ip:       ip,
	}
	err := r.data.DB(ctx).Create(record).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to add an record to db: id(%v), uuid(%s), ip(%s)", id, uuid, ip))
	}
	return nil
}

func (r *commentRepo) SendReviewToMq(ctx context.Context, review *biz.CommentReview) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}
	msg := &primitive.Message{
		Topic: "review",
		Body:  data,
	}
	msg.WithKeys([]string{review.Uuid})
	_, err = r.data.reviewMqPro.producer.SendSync(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send review to mq: %v", err))
	}
	return nil
}
