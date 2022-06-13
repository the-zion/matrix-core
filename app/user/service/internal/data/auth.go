package data

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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

func (r *authRepo) FindUserByPhone(ctx context.Context, phone string) (*biz.User, error) {
	user := &User{}
	err := r.data.DB(ctx).WithContext(ctx).Where("phone = ?", phone).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("phone not found from db", fmt.Sprintf("phone(%s)", phone))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: phone(%s)", phone))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Github:   user.Github,
	}, nil
}

func (r *authRepo) FindUserByEmail(ctx context.Context, email string) (*biz.User, error) {
	user := &User{}
	err := r.data.DB(ctx).WithContext(ctx).Where("email = ?", email).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, kerrors.NotFound("email not found from db", fmt.Sprintf("email(%s)", email))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: email(%s)", email))
	}
	return &biz.User{
		Uuid:     user.Uuid,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
		Wechat:   user.Wechat,
		Github:   user.Github,
	}, nil
}

func (r *authRepo) CreateUserWithPhone(ctx context.Context, phone string) (*biz.User, error) {
	uuid, err := util.UUIdV4()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to create uuid: uuid(%s)", uuid))
	}

	user := &User{
		Uuid:  uuid,
		Phone: phone,
	}
	err = r.data.DB(ctx).Select("Phone", "Uuid").Create(user).Error
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: phone(%s)", phone))
	}

	return &biz.User{
		Uuid:  uuid,
		Phone: phone,
	}, nil
}

func (r *authRepo) CreateUserWithEmail(ctx context.Context, email, password string) (*biz.User, error) {
	uuid, err := util.UUIdV4()
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to create uuid: uuid(%s)", uuid))
	}

	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to hash password: password(%s)", password))
	}

	user := &User{
		Uuid:     uuid,
		Email:    email,
		Password: hashPassword,
	}
	err = r.data.DB(ctx).Select("Email", "Uuid", "Password").Create(user).Error
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "Duplicate") {
			return nil, kerrors.Conflict("email conflict", fmt.Sprintf("email(%s)", email))
		} else {
			return nil, errors.Wrapf(err, fmt.Sprintf("fail to create a user: email(%s)", email))
		}
	}

	return &biz.User{
		Uuid:     uuid,
		Email:    email,
		Password: hashPassword,
	}, nil
}

func (r *authRepo) CreateUserProfile(ctx context.Context, account, uuid string) error {
	err := r.data.DB(ctx).Select("Uuid", "Username").Create(&Profile{Uuid: uuid, Username: account}).Error
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

func (r *authRepo) SendCode(msgs ...*primitive.MessageExt) {
	for _, i := range msgs {
		body := strings.Split(string(i.Body), ";")
		if body[3] == "phone" {
			request := r.data.phoneCodeCli.request
			client := r.data.phoneCodeCli.client
			request.TemplateId = common.StringPtr(util.GetPhoneTemplate(body[2]))
			request.TemplateParamSet = common.StringPtrs([]string{body[1]})
			request.PhoneNumberSet = common.StringPtrs([]string{body[0]})
			_, err := client.SendSms(request)
			if err != nil {
				r.log.Errorf("fail to send phone code: code(%s) error: %v", body[1], err.Error())
			}
		}

		if body[3] == "email" {
			m := r.data.goMailCli.message
			d := r.data.goMailCli.dialer
			m.SetHeader("To", body[0])
			m.SetHeader("Subject", "matrix 魔方技术")
			m.SetBody("text/html", util.GetEmailTemplate(body[2], body[1]))
			err := d.DialAndSend(m)
			if err != nil {
				r.log.Errorf("fail to send email code: code(%s) error: %v", body[1], err.Error())
			}
		}
	}
}

func (r *authRepo) VerifyPhoneCode(ctx context.Context, phone, code string) error {
	key := "phone_" + phone
	return r.verifyCode(ctx, key, code)
}

func (r *authRepo) VerifyEmailCode(ctx context.Context, email, code string) error {
	key := "email_" + email
	return r.verifyCode(ctx, key, code)
}

func (r *authRepo) verifyCode(ctx context.Context, key, code string) error {
	codeInCache, err := r.getCodeFromCache(ctx, key)
	if !kerrors.IsNotFound(err) {
		return err
	}
	if code != codeInCache {
		return errors.Errorf("code error: code(%s)", code)
	}
	r.removeCodeFromCache(ctx, key)
	return nil
}

func (r *authRepo) VerifyPassword(ctx context.Context, account, password, mode string) (*biz.User, error) {
	var err error
	var user *biz.User
	if mode == "phone" {
		user, err = r.FindUserByPhone(ctx, account)
	} else {
		user, err = r.FindUserByEmail(ctx, account)
	}
	if err != nil {
		return nil, err
	}

	pass := util.CheckPasswordHash(password, user.Password)
	if !pass {
		return nil, errors.Errorf("password error: password(%s)", password)
	}
	return user, nil
}

func (r *authRepo) PasswordResetByPhone(ctx context.Context, phone, password string) error {
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to hash password: password(%s)", password))
	}

	err = r.data.DB(ctx).WithContext(ctx).Model(&User{}).Where("phone = ?", phone).Update("password", hashPassword).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reset password: password(%s)", password))
	}
	return nil
}

func (r *authRepo) PasswordResetByEmail(ctx context.Context, email, password string) error {
	hashPassword, err := util.HashPassword(password)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to hash password: password(%s)", password))
	}

	err = r.data.DB(ctx).WithContext(ctx).Model(&User{}).Where("email = ?", email).Update("password", hashPassword).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to reset password: password(%s)", password))
	}
	return nil
}

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
		return "", kerrors.NotFound("code not found from cache", fmt.Sprintf("key(%s)", key))
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
