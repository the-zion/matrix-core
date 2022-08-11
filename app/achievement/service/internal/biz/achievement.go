package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
)

type AchievementRepo interface {
	GetAchievementList(ctx context.Context, uuids []string) ([]*Achievement, error)
	GetUserAchievement(ctx context.Context, uuid string) (*Achievement, error)
	SetAchievementAgree(ctx context.Context, uuid string) error
	SetAchievementView(ctx context.Context, uuid string) error
	SetAchievementCollect(ctx context.Context, uuid string) error
	SetAchievementFollow(ctx context.Context, uuid string) error
	SetAchievementFollowed(ctx context.Context, uuid string) error
	SetAchievementAgreeToCache(ctx context.Context, uuid string) error
	SetAchievementViewToCache(ctx context.Context, uuid string) error
	SetAchievementCollectToCache(ctx context.Context, uuid string) error
	SetAchievementFollowToCache(ctx context.Context, follow, followed string) error
	CancelAchievementAgree(ctx context.Context, uuid string) error
	CancelAchievementAgreeFromCache(ctx context.Context, uuid string) error
	CancelAchievementCollect(ctx context.Context, uuid string) error
	CancelAchievementCollectFromCache(ctx context.Context, uuid string) error
	CancelAchievementFollow(ctx context.Context, uuid string) error
	CancelAchievementFollowed(ctx context.Context, uuid string) error
	CancelAchievementFollowFromCache(ctx context.Context, follow, followed string) error
}

type AchievementUseCase struct {
	repo AchievementRepo
	tm   Transaction
	re   Recovery
	log  *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, re Recovery, tm Transaction, logger log.Logger) *AchievementUseCase {
	return &AchievementUseCase{
		repo: repo,
		tm:   tm,
		re:   re,
		log:  log.NewHelper(log.With(logger, "module", "achievement/biz/AchievementUseCase")),
	}
}

func (r *AchievementUseCase) SetAchievementAgree(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementAgree(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementAgreeFailed("set achievement agree failed", err.Error())
		}
		err = r.repo.SetAchievementAgreeToCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementAgreeFailed("set achievement agree to cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) CancelAchievementAgree(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelAchievementAgree(ctx, uuid)
		if err != nil {
			return v1.ErrorCancelAchievementAgreeFailed("cancel achievement agree failed", err.Error())
		}
		err = r.repo.CancelAchievementAgreeFromCache(ctx, uuid)
		if err != nil {
			return v1.ErrorCancelAchievementAgreeFailed("cancel achievement agree from cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) SetAchievementView(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementView(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementViewFailed("set achievement view failed", err.Error())
		}
		err = r.repo.SetAchievementViewToCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementViewFailed("set achievement view to cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) SetAchievementCollect(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementCollect(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementCollectFailed("set achievement collect failed", err.Error())
		}
		err = r.repo.SetAchievementCollectToCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementCollectFailed("set achievement collect to cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) CancelAchievementCollect(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelAchievementCollect(ctx, uuid)
		if err != nil {
			return v1.ErrorCancelAchievementCollectFailed("cancel achievement collect failed", err.Error())
		}
		err = r.repo.CancelAchievementCollectFromCache(ctx, uuid)
		if err != nil {
			return v1.ErrorCancelAchievementCollectFailed("cancel achievement collect from cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) SetAchievementFollow(ctx context.Context, follow, followed string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementFollow(ctx, followed)
		if err != nil {
			return v1.ErrorSetAchievementFollowFailed("set achievement follow failed", err.Error())
		}

		err = r.repo.SetAchievementFollowed(ctx, follow)
		if err != nil {
			return v1.ErrorSetAchievementFollowFailed("set achievement follow failed", err.Error())
		}

		err = r.repo.SetAchievementFollowToCache(ctx, follow, followed)
		if err != nil {
			return v1.ErrorSetAchievementFollowFailed("set achievement follow to cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) CancelAchievementFollow(ctx context.Context, follow, followed string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelAchievementFollow(ctx, followed)
		if err != nil {
			return v1.ErrorCancelAchievementFollowFailed("cancel achievement follow failed", err.Error())
		}

		err = r.repo.CancelAchievementFollowed(ctx, follow)
		if err != nil {
			return v1.ErrorCancelAchievementFollowFailed("cancel achievement follow failed", err.Error())
		}

		err = r.repo.CancelAchievementFollowFromCache(ctx, follow, followed)
		if err != nil {
			return v1.ErrorCancelAchievementFollowFailed("cancel achievement follow from cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) GetAchievementList(ctx context.Context, uuids []string) ([]*Achievement, error) {
	achievementList, err := r.repo.GetAchievementList(ctx, uuids)
	if err != nil {
		return nil, v1.ErrorGetAchievementListFailed("get achievement list failed: %s", err.Error())
	}
	return achievementList, nil
}

func (r *AchievementUseCase) GetUserAchievement(ctx context.Context, uuid string) (*Achievement, error) {
	achievement, err := r.repo.GetUserAchievement(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetAchievementFailed("get achievement failed: %s", err.Error())
	}
	return achievement, nil
}
