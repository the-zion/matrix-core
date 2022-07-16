package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

type TalkRepo interface {
	GetTalk(ctx context.Context, id int32) (*Talk, error)
	GetTalkList(ctx context.Context, page int32) ([]*Talk, error)
	GetLastTalkDraft(ctx context.Context, uuid string) (*TalkDraft, error)

	CreateTalkDraft(ctx context.Context, uuid string) (int32, error)
	CreateTalkFolder(ctx context.Context, id int32, uuid string) error
	CreateTalk(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkStatistic(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkCache(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkSearch(ctx context.Context, id int32, uuid string) error

	SendTalk(ctx context.Context, id int32, uuid string) (*TalkDraft, error)
	SendReviewToMq(ctx context.Context, review *TalkReview) error
	SendTalkToMq(ctx context.Context, article *Talk, mode string) error

	DeleteTalkDraft(ctx context.Context, id int32, uuid string) error
	EditTalkCos(ctx context.Context, id int32, uuid string) error
	EditTalkSearch(ctx context.Context, id int32, uuid string) error
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

func (r *TalkUseCase) GetTalkList(ctx context.Context, page int32) ([]*Talk, error) {
	talkList, err := r.repo.GetTalkList(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetTalkListFailed("get talk list failed: %s", err.Error())
	}
	return talkList, nil
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
			Mode: "create",
		})
		if err != nil {
			return v1.ErrorCreateTalkFailed("send create review to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) SendTalkEdit(ctx context.Context, id int32, uuid string) error {
	talk, err := r.repo.GetTalk(ctx, id)
	if err != nil {
		return v1.ErrorGetTalkFailed("get talk failed: %s", err.Error())
	}

	if talk.Uuid != uuid {
		return v1.ErrorGetTalkFailed("send talk edit failed: no auth")
	}

	err = r.repo.SendReviewToMq(ctx, &TalkReview{
		Uuid: talk.Uuid,
		Id:   talk.TalkId,
		Mode: "edit",
	})
	if err != nil {
		return v1.ErrorGetTalkFailed("send edit review to mq failed: %s", err.Error())
	}
	return nil
}

func (r *TalkUseCase) CreateTalk(ctx context.Context, id, auth int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteTalkDraft(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateTalkFailed("delete talk draft failed: %s", err.Error())
		}

		err = r.repo.CreateTalk(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateTalkFailed("create talk failed: %s", err.Error())
		}

		err = r.repo.CreateTalkStatistic(ctx, id, auth, uuid)
		if err != nil {
			return v1.ErrorCreateTalkFailed("create talk statistic failed: %s", err.Error())
		}

		err = r.repo.SendTalkToMq(ctx, &Talk{
			TalkId: id,
			Uuid:   uuid,
			Auth:   auth,
		}, "create_talk_cache_and_search")
		if err != nil {
			return v1.ErrorCreateTalkFailed("create talk to mq failed: %s", err.Error())
		}

		return nil
	})
}

func (r *TalkUseCase) EditTalk(ctx context.Context, id int32, uuid string) error {
	err := r.repo.SendTalkToMq(ctx, &Talk{
		TalkId: id,
		Uuid:   uuid,
	}, "edit_talk_cos_and_search")
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk to mq failed: %s", err.Error())
	}
	return nil
}

func (r *TalkUseCase) CreateTalkCacheAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.CreateTalkCache(ctx, id, auth, uuid)
	if err != nil {
		return v1.ErrorCreateTalkFailed("create talk cache failed: %s", err.Error())
	}

	if auth == 2 {
		return nil
	}

	err = r.repo.CreateTalkSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorCreateTalkFailed("create talk search failed: %s", err.Error())
	}
	return nil
}

func (r *TalkUseCase) EditTalkCosAndSearch(ctx context.Context, id int32, uuid string) error {
	err := r.repo.EditTalkCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk cache failed: %s", err.Error())
	}

	err = r.repo.EditTalkSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk search failed: %s", err.Error())
	}
	return nil
}
