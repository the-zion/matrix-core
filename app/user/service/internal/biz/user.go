package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

type UserRepo interface {
	GetAccount(ctx context.Context, uuid string) (*User, error)
	GetProfile(ctx context.Context, uuid string) (*Profile, error)
	GetProfileList(ctx context.Context, uuids []string) ([]*Profile, error)
	GetUserFollow(ctx context.Context, uuid, userUuid string) (bool, error)
	GetUserFollows(ctx context.Context, userId string, uuids []string) ([]*Follows, error)
	GetProfileUpdate(ctx context.Context, uuid string) (*ProfileUpdate, error)
	GetFollowList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowListCount(ctx context.Context, uuid string) (int32, error)
	GetFollowedList(ctx context.Context, page int32, uuid string) ([]*Follow, error)
	GetFollowedListCount(ctx context.Context, uuid string) (int32, error)
	SetProfile(ctx context.Context, profile *ProfileUpdate) error
	SetProfileUpdate(ctx context.Context, profile *ProfileUpdate, status int32) (*ProfileUpdate, error)
	SetUserFollow(ctx context.Context, uuid, userId string) error
	CancelUserFollow(ctx context.Context, uuid, userId string) error
	SendProfileToMq(ctx context.Context, profile *ProfileUpdate) error
	SendUserStatisticToMq(ctx context.Context, uuid, userUuid, mode string) error
	ModifyProfileUpdateStatus(ctx context.Context, uuid, update string) error
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
	tm   Transaction
}

func NewUserUseCase(repo UserRepo, tm Transaction, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/userUseCase")),
		tm:   tm,
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

func (r *UserUseCase) GetUserFollows(ctx context.Context, userId string, uuids []string) ([]*Follows, error) {
	follows, err := r.repo.GetUserFollows(ctx, userId, uuids)
	if err != nil {
		return nil, v1.ErrorGetUserFollowFailed("get user follows failed: %s", err.Error())
	}
	return follows, nil
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
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		if uuid == userId {
			return v1.ErrorSetFollowFailed("uuid and userId are the same: %s")
		}

		err := r.repo.SetUserFollow(ctx, uuid, userId)
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow failed: %s", err.Error())
		}

		err = r.repo.SendUserStatisticToMq(ctx, uuid, userId, "follow")
		if err != nil {
			return v1.ErrorSetFollowFailed("set follow failed: %s", err.Error())
		}

		return nil
	})
}

func (r *UserUseCase) CancelUserFollow(ctx context.Context, uuid, userId string) error {
	return r.tm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.repo.CancelUserFollow(ctx, uuid, userId)
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
