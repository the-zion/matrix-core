package biz

import (
	"context"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

type User struct {
	Uuid     string
	Password string
	Phone    string
	Email    string
	Wechat   string
	Github   string
}

type UserRepo interface {
	GetUser(ctx context.Context, id int64) (*User, error)
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

func (r *UserUseCase) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := r.repo.GetUser(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("get user failed: %s", err.Error())
	}
	return user, nil
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
		profile.Uuid = uuid
		err = r.repo.SetProfile(ctx, profile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return v1.ErrorProfileReviewModify("profile modify failed: %s", err.Error())
	}
	return nil
}

//func (r *UserUseCase) SetUserPhone(ctx context.Context, id int64, phone, password, code string) error {
//	err := r.repo.VerifyCode(ctx, phone, code, "phone")
//	if err != nil {
//		return v1.ErrorVerifyCodeFailed("set phone failed: %s", err.Error())
//	}
//
//	err = r.repo.VerifyPassword(ctx, id, password)
//	if err != nil {
//		return v1.ErrorVerifyPasswordFailed("set phone failed: %s", err.Error())
//	}
//
//	err = r.repo.SetUserPhone(ctx, id, phone)
//	if err != nil {
//		return v1.ErrorSetPhoneFailed("set phone failed: %s", err.Error())
//	}
//	return nil
//}
//
//func (r *UserUseCase) SetUserEmail(ctx context.Context, id int64, email, password, code string) error {
//	err := r.repo.VerifyCode(ctx, email, code, "email")
//	if err != nil {
//		return v1.ErrorVerifyCodeFailed("set email failed: %s", err.Error())
//	}
//
//	err = r.repo.VerifyPassword(ctx, id, password)
//	if err != nil {
//		return v1.ErrorVerifyPasswordFailed("set phone failed: %s", err.Error())
//	}
//
//	err = r.repo.SetUserEmail(ctx, id, email)
//	if err != nil {
//		return v1.ErrorSetEmailFailed("set email failed: %s", err.Error())
//	}
//	return nil
//}
