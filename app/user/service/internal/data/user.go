package data

import (
	"context"
	"encoding/json"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/Cube-v2/cube-core/app/user/service/internal/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var _ biz.UserRepo = (*userRepo)(nil)

var userCacheKey = func(username string) string {
	return "user_" + username
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;size:200"`
	Phone    string `gorm:"uniqueIndex;size:200"`
	Wechat   string `gorm:"uniqueIndex;size:500"`
	Github   string `gorm:"uniqueIndex;size:500"`
	Password string `gorm:"size:500"`
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "user/data/user")),
	}
}

func (r *userRepo) FindByAccount(ctx context.Context, account, mode string) (*biz.User, error) {
	user := &User{}
	if err := r.data.db.WithContext(ctx).Where(mode+" = ?", account).First(user).Error; err != nil {
		r.log.Errorf("fail to get user from db: account(%v) error(%v)", account, err.Error())
		return nil, biz.ErrUserNotFound
	}
	return &biz.User{
		Id: int64(user.Model.ID),
	}, nil
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	key := userCacheKey(strconv.FormatInt(id, 10))
	target, err := r.getUserFromCache(ctx, key)
	if err != nil {
		user := &User{}
		if err = r.data.db.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
			r.log.Errorf("fail to get user from db: id(%v) error(%v)", id, err.Error())
			return nil, biz.ErrUserNotFound
		}
		target = user
		r.setUserToCache(ctx, user, key)
	}
	return &biz.User{Phone: target.Phone, Email: target.Email, Wechat: target.Wechat, Github: target.Github}, nil
}

func (r *userRepo) SendCode(ctx context.Context, template int64, account, mode string) (string, error) {
	var err error
	code := util.RandomNumber()
	switch mode {
	case "phone":
		err = r.sendPhoneCode(template, account, code)
	case "email":
		err = r.sendEmailCode(template, account, code)
	}
	if err != nil {
		return "", biz.ErrSendCodeError
	}

	key := mode + "_" + account
	err = r.setUserCodeToCache(ctx, key, code)
	if err != nil {
		return "", biz.ErrUnknownError
	}
	return code, nil
}

func (r *userRepo) sendPhoneCode(template int64, phone, code string) error {
	request := r.data.phoneCodeCli.request
	client := r.data.phoneCodeCli.client
	request.TemplateId = common.StringPtr(util.GetPhoneTemplate(template))
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	_, err := client.SendSms(request)
	if err != nil {
		r.log.Errorf("fail to send phone code: code(%v) error(%v)", code, err)
		return err
	}
	return nil
}

func (r *userRepo) sendEmailCode(template int64, email, code string) error {
	m := r.data.goMailCli.message
	d := r.data.goMailCli.dialer
	m.SetHeader("To", email)
	m.SetHeader("Subject", "cube 魔方技术")
	m.SetBody("text/html", util.GetEmailTemplate(template, code))
	err := d.DialAndSend(m)
	if err != nil {
		r.log.Errorf("fail to send email code: code(%v) error(%v)", code, err)
		return err
	}
	return nil
}

func (r *userRepo) VerifyCode(ctx context.Context, account, code, mode string) error {
	key := mode + "_" + account
	codeInCache, err := r.getUserCodeFromCache(ctx, key)
	if err != nil && err != redis.Nil {
		return biz.ErrUnknownError
	}
	r.removeUserCodeFromCache(ctx, key)
	if code != codeInCache {
		return biz.ErrCodeError
	}
	return nil
}

func (r *userRepo) PasswordModify(ctx context.Context, id int64, password string) error {
	password, err := util.HashPassword(password)
	if err != nil {
		r.log.Errorf("fail to hash password: password(%v) error(%v)", password, err.Error())
		return biz.ErrUnknownError
	}
	if err = r.data.db.Model(&User{}).Where("id = ?", id).Update("password", password).Error; err != nil {
		r.log.Errorf("fail to modify password: password(%v) error(%v)", password, err.Error())
		return biz.ErrUnknownError
	}
	return nil
}

func (r *userRepo) VerifyPassword(ctx context.Context, id int64, password string) error {
	user := &User{}
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		r.log.Errorf("fail to verify password: password(%v) error(%v)", password, err.Error())
		return biz.ErrUnknownError
	}
	if !util.CheckPasswordHash(password, user.Password) {
		return biz.ErrPasswordError
	}
	return nil
}

func (r *userRepo) SetUserPhone(ctx context.Context, id int64) (string, error) {
	return "", nil
}

func (r *userRepo) getUserCodeFromCache(ctx context.Context, key string) (string, error) {
	code, err := r.data.redisCli.Get(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to get code from cache:redis.Get(%v) error(%v)", key, err)
		return "", err
	}
	return code, nil
}

func (r *userRepo) setUserCodeToCache(ctx context.Context, key, code string) error {
	err := r.data.redisCli.Set(ctx, key, code, time.Minute*5).Err()
	if err != nil {
		r.log.Errorf("fail to set code to cache:redis.Set(%v) error(%v)", key, err)
		return err
	}
	return nil
}

func (r *userRepo) removeUserCodeFromCache(ctx context.Context, key string) {
	_, err := r.data.redisCli.Del(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to delete code from cache:redis.Del(key, %v) error(%v)", key, err)
	}
}

func (r *userRepo) getUserFromCache(ctx context.Context, key string) (*User, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to get user from cache:redis.Get(user, %v) error(%v)", key, err)
		return nil, err
	}
	var cacheUser = &User{}
	err = json.Unmarshal([]byte(result), cacheUser)
	if err != nil {
		return nil, err
	}
	return cacheUser, nil
}

func (r *userRepo) setUserToCache(ctx context.Context, user *User, key string) {
	marshal, err := json.Marshal(user)
	if err != nil {
		r.log.Errorf("fail to set user to json:json.Marshal(%v) error(%v)", user, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		r.log.Errorf("fail to set user to cache:redis.Set(%v) error(%v)", user, err)
	}
}
