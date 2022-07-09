package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
)

type AchievementRepo interface {
	SetAchievementAgree(ctx context.Context, uuid string) error
	SetAchievementView(ctx context.Context, uuid string) error
	SetAchievementCollect(ctx context.Context, uuid string) error
	SetAchievementAgreeToCache(ctx context.Context, uuid string) error
	SetAchievementViewToCache(ctx context.Context, uuid string) error
	SetAchievementCollectToCache(ctx context.Context, uuid string) error
	CancelAchievementAgree(ctx context.Context, uuid string) error
	CancelAchievementAgreeFromCache(ctx context.Context, uuid string) error
	CancelAchievementCollect(ctx context.Context, uuid string) error
	CancelAchievementCollectFromCache(ctx context.Context, uuid string) error
}

type AchievementUseCase struct {
	repo AchievementRepo
	tm   Transaction
	log  *log.Helper
}

func NewAchievementUseCase(repo AchievementRepo, tm Transaction, logger log.Logger) *AchievementUseCase {
	return &AchievementUseCase{
		repo: repo,
		tm:   tm,
		log:  log.NewHelper(log.With(logger, "module", "achievement/biz/AchievementUseCase")),
	}
}

func (r *AchievementUseCase) SetAchievementAgree(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementAgree(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementAgreeFailed("set achievement agree failed", err.Error())
		}
		err = r.repo.SetAchievementAgreeToCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementAgreeFailed("set achievement agree to cache failed", err.Error())
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
			return v1.ErrorSendAchievementViewFailed("set achievement view failed", err.Error())
		}
		err = r.repo.SetAchievementViewToCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementViewFailed("set achievement view to cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) SetAchievementCollect(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementCollect(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementCollectFailed("set achievement collect failed", err.Error())
		}
		err = r.repo.SetAchievementCollectToCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementCollectFailed("set achievement collect to cache failed", err.Error())
		}
		return nil
	})
}
func (r *AchievementUseCase) CancelAchievementCollect(ctx context.Context, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelAchievementCollect(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementCollectFailed("cancel achievement collect failed", err.Error())
		}
		err = r.repo.CancelAchievementCollectFromCache(ctx, uuid)
		if err != nil {
			return v1.ErrorSendAchievementCollectFailed("cancel achievement collect from cache failed", err.Error())
		}
		return nil
	})
}
