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
	GetTalkListHot(ctx context.Context, page int32) ([]*TalkStatistic, error)
	GetUserTalkList(ctx context.Context, page int32, uuid string) ([]*Talk, error)
	GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*Talk, error)
	GetTalkCount(ctx context.Context, uuid string) (int32, error)
	GetTalkCountVisitor(ctx context.Context, uuid string) (int32, error)
	GetTalkListStatistic(ctx context.Context, ids []int32) ([]*TalkStatistic, error)
	GetTalkStatistic(ctx context.Context, id int32) (*TalkStatistic, error)
	GetLastTalkDraft(ctx context.Context, uuid string) (*TalkDraft, error)
	GetTalkAgreeJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetTalkCollectJudge(ctx context.Context, id int32, uuid string) (bool, error)
	GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*TalkSearch, int32, error)

	CreateTalkDraft(ctx context.Context, uuid string) (int32, error)
	CreateTalkFolder(ctx context.Context, id int32, uuid string) error
	CreateTalk(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkStatistic(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkCache(ctx context.Context, id, auth int32, uuid string) error
	CreateTalkSearch(ctx context.Context, id int32, uuid string) error
	AddTalkComment(ctx context.Context, id int32) error
	AddTalkCommentToCache(ctx context.Context, id int32, uuid string) error
	ReduceTalkComment(ctx context.Context, id int32) error
	ReduceTalkCommentToCache(ctx context.Context, id int32, uuid string) error

	SetTalkAgree(ctx context.Context, id int32, uuid string) error
	SetTalkView(ctx context.Context, id int32, uuid string) error
	SetTalkViewToCache(ctx context.Context, id int32, uuid string) error
	SetTalkAgreeToCache(ctx context.Context, id int32, uuid, userUuid string) error
	SetTalkCollect(ctx context.Context, id int32, uuid string) error
	SetTalkUserCollect(ctx context.Context, id, collectionsId int32, userUuid string) error
	SetTalkCollectToCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelTalkAgree(ctx context.Context, id int32, uuid string) error
	CancelTalkAgreeFromCache(ctx context.Context, id int32, uuid, userUuid string) error
	CancelTalkUserCollect(ctx context.Context, id int32, userUuid string) error
	CancelTalkCollect(ctx context.Context, id int32, uuid string) error
	CancelTalkCollectFromCache(ctx context.Context, id int32, uuid, userUuid string) error

	SendTalk(ctx context.Context, id int32, uuid string) (*TalkDraft, error)
	SendReviewToMq(ctx context.Context, review *TalkReview) error
	SendTalkToMq(ctx context.Context, talk *Talk, mode string) error
	SendTalkStatisticToMq(ctx context.Context, uuid, mode string) error

	DeleteTalk(ctx context.Context, id int32, uuid string) error
	DeleteTalkDraft(ctx context.Context, id int32, uuid string) error
	DeleteTalkStatistic(ctx context.Context, id int32, uuid string) error
	DeleteTalkCache(ctx context.Context, id int32, uuid string) error
	DeleteTalkSearch(ctx context.Context, id int32, uuid string) error

	FreezeTalkCos(ctx context.Context, id int32, uuid string) error

	UpdateTalkCache(ctx context.Context, id, auth int32, uuid string) error
	EditTalkCos(ctx context.Context, id int32, uuid string) error
	EditTalkSearch(ctx context.Context, id int32, uuid string) error
}

type TalkUseCase struct {
	repo         TalkRepo
	creationRepo CreationRepo
	tm           Transaction
	log          *log.Helper
}

func NewTalkUseCase(repo TalkRepo, creationRepo CreationRepo, tm Transaction, logger log.Logger) *TalkUseCase {
	return &TalkUseCase{
		repo:         repo,
		creationRepo: creationRepo,
		tm:           tm,
		log:          log.NewHelper(log.With(logger, "module", "creation/biz/talkUseCase")),
	}
}

func (r *TalkUseCase) GetLastTalkDraft(ctx context.Context, uuid string) (*TalkDraft, error) {
	draft, err := r.repo.GetLastTalkDraft(ctx, uuid)
	if kerrors.IsNotFound(err) {
		return nil, v1.ErrorRecordNotFound("last draft not found: %s", err.Error())
	}
	if err != nil {
		return nil, v1.ErrorGetTalkDraftFailed("get last draft failed: %s", err.Error())
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

func (r *TalkUseCase) GetTalkListHot(ctx context.Context, page int32) ([]*TalkStatistic, error) {
	talkList, err := r.repo.GetTalkListHot(ctx, page)
	if err != nil {
		return nil, v1.ErrorGetTalkListFailed("get talk hot list failed: %s", err.Error())
	}
	return talkList, nil
}

func (r *TalkUseCase) GetUserTalkList(ctx context.Context, page int32, uuid string) ([]*Talk, error) {
	talkList, err := r.repo.GetUserTalkList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetTalkListFailed("get user talk list failed: %s", err.Error())
	}
	return talkList, nil
}

func (r *TalkUseCase) GetUserTalkListVisitor(ctx context.Context, page int32, uuid string) ([]*Talk, error) {
	talkList, err := r.repo.GetUserTalkListVisitor(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetTalkListFailed("get user talk list visitor failed: %s", err.Error())
	}
	return talkList, nil
}

func (r *TalkUseCase) GetTalkCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetTalkCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get talk count failed: %s", err.Error())
	}
	return count, nil
}

func (r *TalkUseCase) GetTalkCountVisitor(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetTalkCountVisitor(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetCountFailed("get talk count visitor failed: %s", err.Error())
	}
	return count, nil
}

func (r *TalkUseCase) GetTalkListStatistic(ctx context.Context, ids []int32) ([]*TalkStatistic, error) {
	statisticList, err := r.repo.GetTalkListStatistic(ctx, ids)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get talk list statistic failed: %s", err.Error())
	}
	return statisticList, nil
}

func (r *TalkUseCase) GetTalkStatistic(ctx context.Context, id int32) (*TalkStatistic, error) {
	statistic, err := r.repo.GetTalkStatistic(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetStatisticFailed("get talk statistic failed: %s", err.Error())
	}
	return statistic, nil
}

func (r *TalkUseCase) GetTalkSearch(ctx context.Context, page int32, search, time string) ([]*TalkSearch, int32, error) {
	talkList, total, err := r.repo.GetTalkSearch(ctx, page, search, time)
	if err != nil {
		return nil, 0, v1.ErrorGetTalkSearchFailed("get talk search failed: %s", err.Error())
	}
	return talkList, total, nil
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

func (r *TalkUseCase) SendTalk(ctx context.Context, id int32, uuid, ip string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		draft, err := r.repo.SendTalk(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCreateTalkFailed("send talk failed: %s", err.Error())
		}

		err = r.creationRepo.SetRecord(ctx, id, 3, uuid, "create", ip)
		if err != nil {
			return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
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

func (r *TalkUseCase) SendTalkEdit(ctx context.Context, id int32, uuid, ip string) error {
	talk, err := r.repo.GetTalk(ctx, id)
	if err != nil {
		return v1.ErrorGetTalkFailed("get talk failed: %s", err.Error())
	}

	if talk.Uuid != uuid {
		return v1.ErrorGetTalkFailed("send talk edit failed: no auth")
	}

	err = r.creationRepo.SetRecord(ctx, id, 3, uuid, "edit", ip)
	if err != nil {
		return v1.ErrorSetRecordFailed("set record failed: %s", err.Error())
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

func (r *TalkUseCase) EditTalk(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.SendTalkToMq(ctx, &Talk{
		TalkId: id,
		Auth:   auth,
		Uuid:   uuid,
	}, "edit_talk_cos_and_search")
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk to mq failed: %s", err.Error())
	}
	return nil
}

func (r *TalkUseCase) DeleteTalk(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		var err error
		err = r.repo.DeleteTalk(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteTalkFailed("delete talk failed: %s", err.Error())
		}

		err = r.repo.DeleteTalkStatistic(ctx, id, uuid)
		if err != nil {
			return v1.ErrorDeleteTalkFailed("delete talk statistic failed: %s", err.Error())
		}

		err = r.repo.SendTalkToMq(ctx, &Talk{
			TalkId: id,
			Uuid:   uuid,
		}, "delete_talk_cache_and_search")
		if err != nil {
			return v1.ErrorDeleteTalkFailed("delete talk to mq failed: %s", err.Error())
		}
		return nil
	})
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

func (r *TalkUseCase) EditTalkCosAndSearch(ctx context.Context, id, auth int32, uuid string) error {
	err := r.repo.UpdateTalkCache(ctx, id, auth, uuid)
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk cache failed: %s", err.Error())
	}

	err = r.repo.EditTalkCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk cache failed: %s", err.Error())
	}

	err = r.repo.EditTalkSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorEditTalkFailed("edit talk search failed: %s", err.Error())
	}
	return nil
}

func (r *TalkUseCase) DeleteTalkCacheAndSearch(ctx context.Context, id int32, uuid string) error {
	err := r.repo.DeleteTalkCache(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteTalkFailed("delete talk cache failed: %s", err.Error())
	}

	err = r.repo.FreezeTalkCos(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteTalkFailed("freeze talk cos failed: %s", err.Error())
	}

	err = r.repo.DeleteTalkSearch(ctx, id, uuid)
	if err != nil {
		return v1.ErrorDeleteTalkFailed("delete talk search failed: %s", err.Error())
	}
	return nil
}

func (r *TalkUseCase) SetTalkAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetTalkAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set talk agree failed: %s", err.Error())
		}
		err = r.repo.SetTalkAgreeToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetAgreeFailed("set talk agree to cache failed: %s", err.Error())
		}
		err = r.repo.SendTalkStatisticToMq(ctx, uuid, "agree")
		if err != nil {
			return v1.ErrorSetAgreeFailed("set talk agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) CancelTalkAgree(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelTalkAgree(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel talk agree failed: %s", err.Error())
		}
		err = r.repo.CancelTalkAgreeFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelAgreeFailed("cancel talk agree from cache failed: %s", err.Error())
		}
		err = r.repo.SendTalkStatisticToMq(ctx, uuid, "agree_cancel")
		if err != nil {
			return v1.ErrorCancelAgreeFailed("set talk agree to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) SetTalkView(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetTalkView(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetViewFailed("set talk view failed: %s", err.Error())
		}
		err = r.repo.SetTalkViewToCache(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetViewFailed("set talk view to cache failed: %s", err.Error())
		}
		err = r.repo.SendTalkStatisticToMq(ctx, uuid, "view")
		if err != nil {
			return v1.ErrorSetViewFailed("set talk view to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) SetTalkCollect(ctx context.Context, id, collectionsId int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetTalkUserCollect(ctx, id, collectionsId, userUuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set talk user collect failed: %s", err.Error())
		}
		err = r.repo.SetTalkCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set talk collect failed: %s", err.Error())
		}
		err = r.repo.SetTalkCollectToCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorSetCollectFailed("set talk collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendTalkStatisticToMq(ctx, uuid, "collect")
		if err != nil {
			return v1.ErrorSetCollectFailed("set talk collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) CancelTalkCollect(ctx context.Context, id int32, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelTalkUserCollect(ctx, id, userUuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel talk user collect failed: %s", err.Error())
		}
		err = r.repo.CancelTalkCollect(ctx, id, uuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel talk collect failed: %s", err.Error())
		}
		err = r.repo.CancelTalkCollectFromCache(ctx, id, uuid, userUuid)
		if err != nil {
			return v1.ErrorCancelCollectFailed("cancel talk collect to cache failed: %s", err.Error())
		}
		err = r.repo.SendTalkStatisticToMq(ctx, uuid, "collect_cancel")
		if err != nil {
			return v1.ErrorCancelCollectFailed("set talk collect to mq failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) TalkStatisticJudge(ctx context.Context, id int32, uuid string) (*TalkStatisticJudge, error) {
	agree, err := r.repo.GetTalkAgreeJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetStatisticJudgeFailed("get talk statistic judge failed: %s", err.Error())
	}
	collect, err := r.repo.GetTalkCollectJudge(ctx, id, uuid)
	if err != nil {
		return nil, v1.ErrorGetStatisticJudgeFailed("get talk statistic judge failed: %s", err.Error())
	}
	return &TalkStatisticJudge{
		Agree:   agree,
		Collect: collect,
	}, nil
}

func (r *TalkUseCase) AddTalkComment(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.AddTalkComment(ctx, id)
		if err != nil {
			return v1.ErrorAddCommentFailed("add talk comment failed: %s", err.Error())
		}

		err = r.repo.AddTalkCommentToCache(ctx, id, uuid)
		if err != nil {
			return v1.ErrorAddCommentFailed("add talk comment failed: %s", err.Error())
		}
		return nil
	})
}

func (r *TalkUseCase) ReduceTalkComment(ctx context.Context, id int32, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.ReduceTalkComment(ctx, id)
		if err != nil {
			return v1.ErrorReduceCommentFailed("reduce talk comment failed: %s", err.Error())
		}

		err = r.repo.ReduceTalkCommentToCache(ctx, id, uuid)
		if err != nil {
			return v1.ErrorReduceCommentFailed("reduce talk comment failed: %s", err.Error())
		}
		return nil
	})
}
