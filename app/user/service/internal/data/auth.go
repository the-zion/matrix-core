package data

import (
	"context"
	"fmt"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
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
		user.Phone = account
	case "Email":
		user.Email = account
	}
	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Select(mode).Create(user).Error; err != nil {
			return errors.Wrapf(err, fmt.Sprintf("fail to register a account: account(%v)", account))
		}

		if err := tx.Create(&Profile{UserId: int64(user.ID), Username: account[3:]}).Error; err != nil {
			return errors.Wrapf(err, fmt.Sprintf("fail to register a profile: user_id(%v)", user.ID))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id: int64(user.ID),
	}, nil
}
