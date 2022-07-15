package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type TalkRepo interface {
	GetLastTalkDraft(ctx context.Context, uuid string) (*TalkDraft, error)
	CreateTalkDraft(ctx context.Context, uuid string) (int32, error)
	CreateTalkFolder(ctx context.Context, id int32, uuid string) error
	SendTalk(ctx context.Context, id int32, uuid string) (*TalkDraft, error)
	SendReviewToMq(ctx context.Context, review *TalkReview) error
}

type TalkUseCase struct {
	repo TalkRepo
	tm   Transaction
	log  *log.Helper
}

func NewTalkUseCase(repo TalkRepo, tm Transaction, logger log.Logger) *TalkUseCase {
	return &TalkUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "creation/biz/talkUseCase")),
	}
}

func (r *TalkUseCase) GetLastTalkDraft(ctx context.Context, uuid string) (*TalkDraft, error) {
	draft, err := r.repo.GetLastTalkDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetArticleDraftFailed("get last draft failed: %s", err.Error())
	}
	return draft, nil
}

func (r *TalkUseCase) CreateTalkDraft(ctx context.Context, uuid string) (int32, error) {
	var id int32
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		id, err = r.repo.CreateTalkDraft(ctx, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create talk draft failed: %s", err.Error())
		}

		err = r.repo.CreateTalkFolder(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateDraftFailed("create talk folder failed: %s", err.Error())
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TalkUseCase) SendTalk(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendTalk(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateTalkFailed("send talk failed: %s", err.Error())
		}
		err = r.repo.SendReviewToMq(ctx, &TalkReview{
			Uuid: draft.Uuid,
			Id:   draft.Id,
			Mode: "create_talk",
		})
		if err != nil {
			return v1.ErrorCreateTalkFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}
