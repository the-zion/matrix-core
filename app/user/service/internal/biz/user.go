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
	GetProfileUpdate(ctx context.Context, uuid string) (*ProfileUpdate, error)
	SetProfile(ctx context.Context, profile *ProfileUpdate) error
	SetProfileUpdate(ctx context.Context, profile *ProfileUpdate, status int32) (*ProfileUpdate, error)
	SendProfileToMq(ctx context.Context, profile *ProfileUpdate) error
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

func (r *UserUseCase) GetProfileUpdate(ctx context.Context, uuid string) (*ProfileUpdate, error) {
	profile, err := r.repo.GetProfileUpdate(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetProfileUpdateFailed("get user profile update failed: %s", err.Error())
	}
	return profile, nil
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
