package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"golang.org/x/sync/errgroup"
)

type UserRepo interface {
	GetAccount(ctx context.Context, uuid string) (*User, error)
	GetProfile(ctx context.Context, uuid string) (*Profile, error)
	GetProfileList(ctx context.Context, uuids []string) ([]*Profile, error)
	GetUserFollow(ctx context.Context, uuid, userUuid string) (bool, error)
	GetUserFollows(ctx context.Context, uuid string) ([]string, error)
	GetProfileUpdate(ctx context.Context, uuid string) (*ProfileUpdate, error)
	GetFollowList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowListCount(ctx context.Context, uuid string) (int32, error)
	GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowedListCount(ctx context.Context, uuid string) (int32, error)
	GetUserSearch(ctx context.Context, page int32, search string) ([]*UserSearch, int32, error)
	EditUserSearch(ctx context.Context, uuid string, profile *ProfileUpdate) error
	SetProfile(ctx context.Context, profile *ProfileUpdate) error
	SetProfileUpdate(ctx context.Context, profile *ProfileUpdate, status int32) (*ProfileUpdate, error)
	SetUserFollow(ctx context.Context, uuid, userId string) error
	SetUserFollowToCache(ctx context.Context, uuid, userId string) error
	SetFollowToMq(ctx context.Context, follow *Follow, mode string) error
	CancelUserFollow(ctx context.Context, uuid, userId string) error
	CancelUserFollowFromCache(ctx context.Context, uuid, userId string) error
	SendProfileToMq(ctx context.Context, profile *ProfileUpdate) error
	SendUserStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error
	ModifyProfileUpdateStatus(ctx context.Context, uuid, update string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
	re   Recovery
	tm   Transaction
}

func NewUserUseCase(repo UserRepo, re Recovery, tm Transaction, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/userUseCase")),
		tm:   tm,
		re:   re,
	}
}

func (r *UserUseCase) GetAccount(ctx context.Context, uuid string) (*User, error) {
	account, err := r.repo.GetAccount(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetAccountFailed("get user account failed: %s", err.Error())
	}
	if account.Password != "" {
		account.Password = "********"
	}
	return account, nil
}

func (r *UserUseCase) GetProfile(ctx context.Context, uuid string) (*Profile, error) {
	profile, err := r.repo.GetProfile(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetProfileFailed("get user profile failed: %s", err.Error())
	}
	return profile, nil
}

func (r *UserUseCase) GetProfileList(ctx context.Context, uuids []string) ([]*Profile, error) {
	profileList, err := r.repo.GetProfileList(ctx, uuids)
	if err != nil {
		return nil, v1.ErrorGetProfileFailed("get user profile list failed: %s", err.Error())
	}
	return profileList, nil
}

func (r *UserUseCase) GetProfileUpdate(ctx context.Context, uuid string) (*ProfileUpdate, error) {
	profile, err := r.repo.GetProfileUpdate(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetProfileUpdateFailed("get user profile update failed: %s", err.Error())
	}
	return profile, nil
}

func (r *UserUseCase) GetUserFollow(ctx context.Context, uuid, userUuid string) (bool, error) {
	follow, err := r.repo.GetUserFollow(ctx, uuid, userUuid)
	if err != nil {
		return false, v1.ErrorGetUserFollowFailed("get user follow failed: %s", err.Error())
	}
	return follow, nil
}

func (r *UserUseCase) GetUserFollows(ctx context.Context, uuid string) (map[string]bool, error) {
	follows, err := r.repo.GetUserFollows(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetUserFollowFailed("get user follows failed: %s", err.Error())
	}
	followsMap := make(map[string]bool, 0)
	for _, item := range follows {
		followsMap[item] = true
	}
	return followsMap, nil
}

func (r *UserUseCase) GetUserSearch(ctx context.Context, page int32, search string) ([]*UserSearch, int32, error) {
	userList, total, err := r.repo.GetUserSearch(ctx, page, search)
	if err != nil {
		return nil, 0, v1.ErrorGetUserSearchFailed("get user search failed: %s", err.Error())
	}
	return userList, total, nil
}

func (r *UserUseCase) GetFollowList(ctx context.Context, page int32, uuid string) ([]*Follow, error) {
	follow, err := r.repo.GetFollowList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetFollowListFailed("get follow list failed: %s", err.Error())
	}
	return follow, nil
}

func (r *UserUseCase) GetFollowListCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetFollowListCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetFollowListCountFailed("get follow list count failed: %s", err.Error())
	}
	return count, nil
}

func (r *UserUseCase) GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error) {
	follow, err := r.repo.GetFollowedList(ctx, page, uuid)
	if err != nil {
		return nil, v1.ErrorGetFollowListFailed("get followed list failed: %s", err.Error())
	}
	return follow, nil
}

func (r *UserUseCase) GetFollowedListCount(ctx context.Context, uuid string) (int32, error) {
	count, err := r.repo.GetFollowedListCount(ctx, uuid)
	if err != nil {
		return 0, v1.ErrorGetFollowListCountFailed("get followed list count failed: %s", err.Error())
	}
	return count, nil
}

func (r *UserUseCase) SetProfileUpdate(ctx context.Context, profile *ProfileUpdate) error {
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		p, err := r.repo.SetProfileUpdate(ctx, profile, 2)
		if err != nil {
			return err
		}
		err = r.repo.SendProfileToMq(ctx, p)
		if err != nil {
			return err
		}
		return nil
	})
	if kerrors.IsConflict(err) {
		return v1.ErrorUserNameConflict("user profile update failed: %s", err.Error())
	}
	if err != nil {
		return v1.ErrorSetProfileUpdateFailed("user profile update failed: %s", err.Error())
	}
	return nil
}

func (r *UserUseCase) SetUserFollow(ctx context.Context, uuid, userId string) error {
	if uuid == userId {
		return v1.ErrorSetFollowFailed("uuid and userId are the same: %s")
	}
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserFollowToCache(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow to cache failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetFollowToMq(ctx, &Follow{
			Follow:   uuid,
			Followed: userId,
		}, "set_follow_db_and_cache")
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow failed: %s", err.Error())
		}
		return nil
	}))

	err := g.Wait()
	if err != nil {
		return nil
	}
	return nil
}

func (r *UserUseCase) SetFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.SetUserFollow(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow failed: %s", err.Error())
		}

		err = r.repo.SetUserFollowToCache(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow to cache failed: %s", err.Error())
		}

		err = r.repo.SendUserStatisticToMq(ctx, uuid, userId, "follow")
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow failed: %s", err.Error())
		}

		return nil
	})
}

func (r *UserUseCase) CancelUserFollow(ctx context.Context, uuid, userId string) error {
	g, _ := errgroup.WithContext(ctx)
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserFollowFromCache(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorCancelFollowFailed("cancel follow failed: %s", err.Error())
		}
		return nil
	}))
	g.Go(r.re.GroupRecover(ctx, func(ctx context.Context) error {
		err := r.repo.SetFollowToMq(ctx, &Follow{
			Follow:   uuid,
			Followed: userId,
		}, "cancel_follow_db_and_cache")
		if err != nil {
			return v1.ErrorSetFollowFailed("cancel follow failed: %s", err.Error())
		}
		return nil
	}))

	err := g.Wait()
	if err != nil {
		return nil
	}
	return nil
}

func (r *UserUseCase) CancelFollowDbAndCache(ctx context.Context, uuid, userId string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserFollow(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorCancelFollowFailed("cancel follow failed: %s", err.Error())
		}

		err = r.repo.CancelUserFollowFromCache(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorCancelFollowFailed("cancel follow failed: %s", err.Error())
		}

		err = r.repo.SendUserStatisticToMq(ctx, uuid, userId, "follow_cancel")
		if err != nil {
			return v1.ErrorSetFollowFailed("cancel follow failed: %s", err.Error())
		}

		return nil
	})
}

func (r *UserUseCase) ProfileReviewPass(ctx context.Context, uuid, update string) error {
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.ModifyProfileUpdateStatus(ctx, uuid, update)
		if err != nil {
			return err
		}
		profile, err := r.repo.GetProfileUpdate(ctx, uuid)
		if err != nil {
			return err
		}
		err = r.repo.SetProfile(ctx, profile)
		if err != nil {
			return err
		}
		err = r.repo.EditUserSearch(ctx, uuid, profile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return v1.ErrorProfileReviewModifyFailed("profile modify failed: %s", err.Error())
	}
	return nil
}

func (r *UserUseCase) ProfileReviewNotPass(ctx context.Context, uuid string) error {
	err := r.tm.ExecTx(ctx, func(ctx context.Context) error {
		profile, err := r.repo.GetProfile(ctx, uuid)
		if err != nil {
			return err
		}
		_, err = r.repo.SetProfileUpdate(ctx, &ProfileUpdate{
			Profile: *profile,
		}, 1)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return v1.ErrorProfileUpdateModifyFailed("profile update modify failed: %s", err.Error())
	}
	return nil
}
