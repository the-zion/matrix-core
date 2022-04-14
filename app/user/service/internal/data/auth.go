package data

import (
	"context"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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
	user.Username = account
	switch mode {
	case "Phone":
		user.Phone = account
	case "Email":
		user.Email = account
	}
	result := r.data.db.WithContext(ctx).Select(mode, "Username").Create(user)
	if result.Error != nil {
		r.log.Errorf("fail to register a account: account(%v) error(%v)", account, result.Error.Error())
		return nil, biz.ErrUserRegisterFailed
	}
	return &biz.User{
		Id: int64(user.Model.ID),
	}, nil
}
