package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	Id       int64
	Username string
	Password string
}

type UserRepo interface {
	FindByUserPhone(ctx context.Context, phone string) (*User, error)
	FindByUserEmail(ctx context.Context, email string) (*User, error)
	VerifyPassword(ctx context.Context, u *User, password string) error
}

type UserUseCase struct {
	repo   UserRepo
	log    *log.Helper
	authUc *AuthUseCase
}

func NewUserUseCase(repo UserRepo, logger log.Logger, authUc *AuthUseCase) *UserUseCase {
	l := log.NewHelper(log.With(logger, "module", "usecase/user"))
	return &UserUseCase{
		repo:   repo,
		log:    l,
		authUc: authUc,
	}
}
