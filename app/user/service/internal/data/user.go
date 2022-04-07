package data

import (
	"context"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;size:500"`
	Phone        string `gorm:"uniqueIndex;size:500"`
	Username     string `gorm:"uniqueIndex;size:500"`
	Wechat       string `gorm:"uniqueIndex"`
	Github       string `gorm:"uniqueIndex"`
	PasswordHash string `gorm:"size:500"`
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/server-service")),
	}
}

func (r *userRepo) FindByUserPhone(ctx context.Context, phone string) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) FindByUserEmail(ctx context.Context, email string) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) VerifyPassword(ctx context.Context, u *biz.User, password string) error {
	return nil
}

func (r *userRepo) VerifyCode(ctx context.Context, u *biz.User, code string) error {
	return nil
}

func (r *userRepo) UserRegister(ctx context.Context, u *biz.User) error {
	return nil
}
