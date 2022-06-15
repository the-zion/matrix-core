package biz

import (
	"context"
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
}

type UserUseCase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUseCase(repo UserRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/userUseCase")),
	}
}

func (r *UserUseCase) GetUserProfile(ctx context.Context, uuid string) (*Profile, error) {
	profile, err := r.repo.GetProfile(ctx, uuid)
	if err != nil {
		return nil, v1.ErrorGetProfileFailed("get user profile failed: %s", err.Error())
	}
	return &Profile{
		Uuid:      profile.Uuid,
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		School:    profile.School,
		Company:   profile.Company,
		Homepage:  profile.Homepage,
		Introduce: profile.Introduce,
	}, nil
}

func (r *UserUseCase) GetUser(ctx context.Context, id int64) (*User, error) {
	user, err := r.repo.GetUser(ctx, id)
	if err != nil {
		return nil, v1.ErrorGetUserFailed("get user failed: %s", err.Error())
	}
	return user, nil
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
