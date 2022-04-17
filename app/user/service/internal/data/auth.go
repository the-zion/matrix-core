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

func NewAuthRepo(data *Data, logger log.Logger) biz.AuthRepo {
	return &authRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/auth")),
	}
}

func (r *authRepo) UserRegister(ctx context.Context, account, mode string) (*biz.User, error) {
	user := &User{}
	switch mode {
	case "Phone":
		account = account[3:]
		user.Phone = account
	case "Email":
		user.Email = account
	}
	if err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Select(mode).Create(user).Error; err != nil {
			r.log.Errorf("fail to register a account: account(%v) error(%v)", account, err.Error())
			return err
		}

		if err := tx.Create(&Profile{UserId: int64(user.ID), Username: account}).Error; err != nil {
			r.log.Errorf("fail to register a profile: user_id(%v) error(%v)", user.ID, err.Error())
			return err
		}
		return nil
	}); err != nil {
		return nil, biz.ErrUserRegisterFailed
	}
	return &biz.User{
		Id: int64(user.ID),
	}, nil
}
