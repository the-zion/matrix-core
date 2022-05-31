package data

import (
	"context"
	"fmt"
	"github.com/Cube-v2/matrix-core/app/user/service/internal/biz"
	v2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
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

func (r *authRepo) FindByAccount(ctx context.Context, account, mode string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Where(mode+" = ?", account).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v2.NotFound("account not found from db", fmt.Sprintf("account(%s), mode(%s) ", account, mode))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: account(%s), type(%s)", account, mode))
	}
	return &biz.User{
		Id: int64(user.Model.ID),
	}, nil
}

func (r *authRepo) VerifyCode(ctx context.Context, account, code, mode string) error {
	key := mode + "_" + account
	codeInCache, err := r.getCodeFromCache(ctx, key)
	if !v2.IsNotFound(err) {
		return err
	}
	if code != codeInCache {
		return errors.Errorf("code error")
	}
	r.removeUserCodeFromCache(ctx, key)
	return nil
}

func (r *authRepo) RegisterWithPhone(ctx context.Context, phone string) (*biz.User, error) {
	user := &User{Phone: phone}
	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("Phone").Create(user).Error; err != nil {
			return errors.Wrapf(err, fmt.Sprintf("fail to register a account: account(%v)", phone))
		}

		if err := tx.Create(&Profile{UserId: int64(user.ID), Username: phone}).Error; err != nil {
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

func (r *authRepo) RegisterWithEmail(ctx context.Context, email string) (*biz.User, error) {
	user := &User{Email: email}
	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("Email").Create(user).Error; err != nil {
			return errors.Wrapf(err, fmt.Sprintf("fail to register a account: account(%v)", email))
		}

		if err := tx.Create(&Profile{UserId: int64(user.ID), Username: email}).Error; err != nil {
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

//func (r *authRepo) UserRegister(ctx context.Context, account, mode string) (*biz.User, error) {
//	user := &User{}
//	switch mode {
//	case "Phone":
//		user.Phone = account
//	case "Email":
//		user.Email = account
//	}
//	err := r.data.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//		if err := tx.Select(mode).Create(user).Error; err != nil {
//			return errors.Wrapf(err, fmt.Sprintf("fail to register a account: account(%v)", account))
//		}
//
//		if err := tx.Create(&Profile{UserId: int64(user.ID), Username: account[3:]}).Error; err != nil {
//			return errors.Wrapf(err, fmt.Sprintf("fail to register a profile: user_id(%v)", user.ID))
//		}
//		return nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	return &biz.User{
//		Id: int64(user.ID),
//	}, nil
//}

func (r *authRepo) getCodeFromCache(ctx context.Context, key string) (string, error) {
	code, err := r.data.redisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", v2.NotFound("code not found from cache", fmt.Sprintf("key(%s)", key))
	}
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get code from cache: redis.Get(%v)", key))
	}
	return code, nil
}

func (r *authRepo) removeUserCodeFromCache(ctx context.Context, key string) {
	_, err := r.data.redisCli.Del(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to delete code from cache: redis.Del(key, %v), error(%v)", key, err)
	}
}
