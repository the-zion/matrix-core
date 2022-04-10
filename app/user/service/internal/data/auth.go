package data

import (
	"context"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var _ biz.AuthRepo = (*authRepo)(nil)

type authRepo struct {
	data *Data
	log  *log.Helper
}

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;size:200"`
	Phone    string `gorm:"uniqueIndex;size:200"`
	Username string `gorm:"uniqueIndex;size:200"`
	Wechat   string `gorm:"uniqueIndex;size:500"`
	Github   string `gorm:"uniqueIndex;size:500"`
	Password string `gorm:"size:500"`
}

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/auth")),
	}
}

func (r *authRepo) FindByUserPhone(ctx context.Context, phone string) (*biz.Auth, error) {
	user := &User{}
	result := r.data.db.WithContext(ctx).Where("phone = ?", phone).First(user)
	if result.Error != nil {
		return nil, biz.ErrUserNotFound
	}
	return &biz.Auth{
		Id: int64(user.Model.ID),
	}, nil
}

func (r *authRepo) FindByUserEmail(ctx context.Context, email string) (*biz.Auth, error) {
	user := &User{}
	result := r.data.db.WithContext(ctx).Where("email = ?", email).First(user)
	if result.Error != nil {
		return nil, biz.ErrUserNotFound
	}
	return &biz.Auth{
		Id: int64(user.Model.ID),
	}, nil
}

func (r *authRepo) VerifyPassword(ctx context.Context, id int64, password string) error {
	user := &User{}
	result := r.data.db.WithContext(ctx).Where("id = ? AND password = ?", id, password).First(&user)
	if result.Error != nil {
		return biz.ErrPasswordError
	}
	return nil
}

func (r *authRepo) VerifyCode(ctx context.Context, key, code string) error {
	return nil
}

func (r *authRepo) UserRegister(ctx context.Context, u *biz.Auth, mode string) (*biz.Auth, error) {
	user := &User{
		Email: u.Email,
		Phone: u.Phone,
	}
	result := r.data.db.WithContext(ctx).Select(mode).Create(user)
	if result.Error != nil {
		return nil, biz.ErrUserRegisterFailed
	}
	return &biz.Auth{
		Id: int64(user.Model.ID),
	}, nil
}
