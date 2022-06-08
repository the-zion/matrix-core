package data

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	v2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/pkg/util"
	"gorm.io/gorm"
	"strings"
	"time"
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

func (r *authRepo) FindUserByPhone(ctx context.Context, phone string) (string, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Where("phone = ?", phone).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", v2.NotFound("phone not found from db", fmt.Sprintf("phone(%s)", phone))
	}
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("db query system error: phone(%s)", phone))
	}
	return user.Uuid, nil
}

func (r *authRepo) CreateUserWithPhone(ctx context.Context, phone string) (string, error) {
	uuid, err := util.UUIdV4()
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to create uuid: phone(%s)", phone))
	}

	user := &User{
		Uuid:  uuid,
		Phone: phone,
	}
	err = r.data.DB(ctx).Select("Phone").Create(user).Error
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to create a user: phone(%s)", phone))
	}

	return uuid, nil
}

func (r *authRepo) CreateUserProfile(ctx context.Context, account, uuid string) error {
	err := r.data.DB(ctx).Create(&Profile{Uuid: uuid, Username: account}).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to register a profile: uuid(%s)", uuid))
	}
	return nil
}

func (r *authRepo) SendPhoneCode(ctx context.Context, template, phone string) error {
	code := util.RandomNumber()
	err := r.setCodeToCache(ctx, "phone_"+phone, code)
	if err != nil {
		return err
	}

	message := strings.Join([]string{phone, code, template, "phone"}, ";")
	msg := &primitive.Message{
		Topic: "code",
		Body:  []byte(message),
	}
	msg.WithTag("phone")
	err = r.data.mqPro.SendOneWay(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send code to producer: %s", message))
	}

	return nil
}

func (r *authRepo) SendEmailCode(ctx context.Context, template, email string) error {
	code := util.RandomNumber()
	err := r.setCodeToCache(ctx, "email_"+email, code)
	if err != nil {
		return err
	}

	message := strings.Join([]string{email, code, template, "email"}, ";")
	msg := &primitive.Message{
		Topic: "code",
		Body:  []byte(message),
	}
	msg.WithTag("email")
	err = r.data.mqPro.SendOneWay(ctx, msg)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send code to producer: %s", message))
	}

	return nil
}

func (r *authRepo) VerifyPhoneCode(ctx context.Context, phone, code string) error {
	key := "phone_" + phone
	return r.verifyCode(ctx, key, code)
}

func (r *authRepo) verifyCode(ctx context.Context, key, code string) error {
	codeInCache, err := r.getCodeFromCache(ctx, key)
	if !v2.IsNotFound(err) {
		return err
	}
	if code != codeInCache {
		return errors.Errorf("code error")
	}
	r.removeCodeFromCache(ctx, key)
	return nil
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

func (r *authRepo) setCodeToCache(ctx context.Context, key, code string) error {
	err := r.data.redisCli.Set(ctx, key, code, time.Minute*2).Err()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set code to cache: redis.Set(%v), code(%s)", key, code))
	}
	return nil
}

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

func (r *authRepo) removeCodeFromCache(ctx context.Context, key string) {
	_, err := r.data.redisCli.Del(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to delete code from cache: redis.Del(key, %v), error(%v)", key, err)
	}
}
