package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/achievement/service/v1"
	"golang.org/x/sync/errgroup"
)

type AchievementRepo interface {
	GetAchievementList(ctx context.Context, uuids []string) ([]*Achievement, error)
	GetUserAchievement(ctx context.Context, uuid string) (*Achievement, error)
	GetUserMedal(ctx context.Context, uuid string) (*Medal, error)
	GetUserActive(ctx context.Context, uuid string) (*Active, error)
	SetAchievementAgree(ctx context.Context, uuid string) error
	SetActiveAgree(ctx context.Context, userUuid string) error
	SetAchievementView(ctx context.Context, uuid string) error
	SetAchievementCollect(ctx context.Context, uuid string) error
	SetAchievementFollow(ctx context.Context, uuid string) error
	SetAchievementFollowed(ctx context.Context, uuid string) error
	SetAchievementAgreeToCache(ctx context.Context, uuid string) error
	SetAchievementViewToCache(ctx context.Context, uuid string) error
	SetAchievementCollectToCache(ctx context.Context, uuid string) error
	SetAchievementFollowToCache(ctx context.Context, follow, followed string) error
	SetUserMedalToCache(ctx context.Context, medal, uuid string) error
	SetUserMedal(ctx context.Context, medal, uuid string) error
	CancelAchievementAgree(ctx context.Context, uuid string) error
	CancelActiveAgree(ctx context.Context, userUuid string) error
	CancelAchievementAgreeFromCache(ctx context.Context, uuid string) error
	CancelAchievementCollect(ctx context.Context, uuid string) error
	CancelAchievementCollectFromCache(ctx context.Context, uuid string) error
	CancelAchievementFollow(ctx context.Context, uuid string) error
	CancelAchievementFollowed(ctx context.Context, uuid string) error
	CancelAchievementFollowFromCache(ctx context.Context, follow, followed string) error
	CancelUserMedalFromCache(ctx context.Context, medal, uuid string) error
	CancelUserMedal(ctx context.Context, medal, uuid string) error
	AddAchievementScore(ctx context.Context, uuid string, score int32) error
	AddAchievementScoreToCache(ctx context.Context, uuid string, score int32) error
	SendMedalToMq(ctx context.Context, medal, uuid, mode string) error
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

func (r *AchievementUseCase) SetAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetAchievementAgree(ctx, uuid)
		if err != nil {
			return v1.ErrorSetAchievementAgreeFailed("set achievement agree failed", err.Error())
		}
		err = r.repo.SetActiveAgree(ctx, userUuid)
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

func (r *AchievementUseCase) CancelAchievementAgree(ctx context.Context, uuid, userUuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelAchievementAgree(ctx, uuid)
		if err != nil {
			return v1.ErrorCancelAchievementAgreeFailed("cancel achievement agree failed", err.Error())
		}
		err = r.repo.CancelActiveAgree(ctx, userUuid)
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

func (r *AchievementUseCase) AddAchievementScore(ctx context.Context, uuid string, score int32) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.AddAchievementScore(ctx, uuid, score)
		if err != nil {
			return v1.ErrorAddAchievementScoreFailed("add achievement score failed", err.Error())
		}
		err = r.repo.AddAchievementScoreToCache(ctx, uuid, score)
		if err != nil {
			return v1.ErrorAddAchievementScoreFailed("add achievement score failed", err.Error())
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

func (r *AchievementUseCase) GetUserMedal(ctx context.Context, uuid string) (*Medal, error) {
	medal, err := r.repo.GetUserMedal(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetMedalFailed("get medal failed: %s", err.Error())
	}
	return medal, nil
}

func (r *AchievementUseCase) GetUserActive(ctx context.Context, uuid string) (*Active, error) {
	active, err := r.repo.GetUserActive(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetActiveFailed("get active failed: %s", err.Error())
	}
	return active, nil
}

func (r *AchievementUseCase) SetUserMedal(ctx context.Context, medal, uuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserMedalToCache(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorSetMedalFailed("set user medal to cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendMedalToMq(ctx, medal, uuid, "set_user_medal_db_and_cache")
		if err != nil {
			return v1.ErrorSetMedalFailed("set user medal to mq failed: %s", err.Error())
		}
		return nil
	}))
	return g.Wait()
}

func (r *AchievementUseCase) SetUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserMedal(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorSetMedalFailed("set medal failed", err.Error())
		}
		err = r.repo.SetUserMedalToCache(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorSetMedalFailed("set medal to cache failed", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) CancelUserMedalSet(ctx context.Context, medal, uuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserMedalFromCache(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorCancelMedalSetFailed("cancel user medal from cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendMedalToMq(ctx, medal, uuid, "cancel_user_medal_db_and_cache")
		if err != nil {
			return v1.ErrorCancelMedalSetFailed("cancel user medal to mq failed: %s", err.Error())
		}
		return nil
	}))
	return g.Wait()
}

func (r *AchievementUseCase) CancelUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserMedal(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorCancelMedalSetFailed("cancel user medal set failed", err.Error())
		}
		err = r.repo.CancelUserMedalFromCache(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorCancelMedalSetFailed("cancel user medal set failed: %s", err.Error())
		}
		return nil
	})
}

func (r *AchievementUseCase) AccessUserMedal(ctx context.Context, medal, uuid string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserMedalToCache(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorAccessMedalFailed("access user medal to cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SendMedalToMq(ctx, medal, uuid, "access_user_medal_db_and_cache")
		if err != nil {
			return v1.ErrorAccessMedalFailed("access user medal to mq failed: %s", err.Error())
		}
		return nil
	}))
	return g.Wait()
}

func (r *AchievementUseCase) AccessUserMedalDbAndCache(ctx context.Context, medal, uuid string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserMedal(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorCancelMedalSetFailed("cancel user medal set failed", err.Error())
		}
		err = r.repo.SetUserMedalToCache(ctx, medal, uuid)
		if err != nil {
			return v1.ErrorAccessMedalFailed("access user medal to cache failed: %s", err.Error())
		}
		return nil
	})
}
