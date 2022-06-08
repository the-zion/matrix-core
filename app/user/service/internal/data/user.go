package data

import (
	"context"
	"encoding/json"
	"fmt"
	v2 "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"github.com/the-zion/matrix-core/app/user/service/internal/pkg/util"
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
	Uuid     string `gorm:"uniqueIndex;size:200"`
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

func (r *userRepo) FindByAccount(ctx context.Context, account, types string) (*biz.User, error) {
	user := &User{}
	err := r.data.db.WithContext(ctx).Where(types+" = ?", account).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, v2.NotFound("account not found from db", fmt.Sprintf("account(%s), mode(%s) ", account, types))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: account(%s), type(%s)", account, types))
	}
	return &biz.User{
		Uuid: user.Uuid,
	}, nil
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	key := userCacheKey(strconv.FormatInt(id, 10))
	target, err := r.getUserFromCache(ctx, key)
	if v2.IsNotFound(err) {
		user := &User{}
		err = r.data.db.WithContext(ctx).Where("id = ?", id).First(user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v2.NotFound("user not found from db", fmt.Sprintf("user_id(%v)", id))
		}
		if err != nil {
			return nil, errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v)", id))
		}
		target = user
		r.setUserToCache(ctx, user, key)
	}
	if err != nil {
		return nil, err
	}
	return &biz.User{Phone: target.Phone, Email: target.Email, Wechat: target.Wechat, Github: target.Github}, nil
}

func (r *userRepo) SendCode(ctx context.Context, template int64, account, mode string) error {
	var err error
	code := util.RandomNumber()
	switch mode {
	case "phone":
		err = r.sendPhoneCode(template, account, code)
	case "email":
		err = r.sendEmailCode(template, account, code)
	}
	if err != nil {
		return err
	}

	key := mode + "_" + account
	err = r.setUserCodeToCache(ctx, key, code)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) sendPhoneCode(template int64, phone, code string) error {
	request := r.data.phoneCodeCli.request
	client := r.data.phoneCodeCli.client
	request.TemplateId = common.StringPtr(util.GetPhoneTemplate(string(template)))
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	_, err := client.SendSms(request)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send phone code: code(%v)", code))
	}
	return nil
}

func (r *userRepo) sendEmailCode(template int64, email, code string) error {
	m := r.data.goMailCli.message
	d := r.data.goMailCli.dialer
	m.SetHeader("To", email)
	m.SetHeader("Subject", "cube 魔方技术")
	m.SetBody("text/html", util.GetEmailTemplate(string(template), code))
	err := d.DialAndSend(m)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to send email code: code(%v)", code))
	}
	return nil
}

func (r *userRepo) VerifyCode(ctx context.Context, account, code, mode string) error {
	key := mode + "_" + account
	codeInCache, err := r.getCodeFromCache(ctx, key)
	if !v2.IsNotFound(err) {
		return err
	}
	r.removeUserCodeFromCache(ctx, key)
	if code != codeInCache {
		return errors.Errorf("code error")
	}
	return nil
}

func (r *userRepo) PasswordModify(ctx context.Context, id int64, password string) error {
	password, err := util.HashPassword(password)
	if err != nil {

		//Name interface {}

	}
	if err = r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", password).Error; err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to modify password: password(%s), user_id(%v)", password, id))
	}
	return nil
}

func (r *userRepo) VerifyPassword(ctx context.Context, id int64, password string) error {
	user := &User{}
	err := r.data.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return v2.NotFound("user not found from db", fmt.Sprintf("user_id(%v)", id))
	}
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v)", id))
	}
	if !util.CheckPasswordHash(password, user.Password) {
		return errors.Errorf("password error")
	}
	return nil
}

// SetUserPhone set redis ?
func (r *userRepo) SetUserPhone(ctx context.Context, id int64, phone string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("phone", phone).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v), phone(%s)", id, phone))
	}
	return nil
}

// SetUserEmail set redis ?
func (r *userRepo) SetUserEmail(ctx context.Context, id int64, email string) error {
	err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("email", email).Error
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query system error: user_id(%v), email(%s)", id, email))
	}
	return nil
}

func (r *userRepo) getCodeFromCache(ctx context.Context, key string) (string, error) {
	code, err := r.data.redisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", v2.NotFound("code not found from cache", fmt.Sprintf("key(%s)", key))
	}
	if err != nil {
		return "", errors.Wrapf(err, fmt.Sprintf("fail to get code from cache: redis.Get(%v)", key))
	}
	return code, nil
}

func (r *userRepo) setUserCodeToCache(ctx context.Context, key, code string) error {
	err := r.data.redisCli.Set(ctx, key, code, time.Minute*5).Err()
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("fail to set code to cache: redis.Set(%v), code(%s)", key, code))
	}
	return nil
}

func (r *userRepo) removeUserCodeFromCache(ctx context.Context, key string) {
	_, err := r.data.redisCli.Del(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to delete code from cache: redis.Del(key, %v), error(%v)", key, err)
	}
}

func (r *userRepo) getUserFromCache(ctx context.Context, key string) (*User, error) {
	result, err := r.data.redisCli.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, v2.NotFound("user not found from cache", fmt.Sprintf("key(%s)", key))
	}
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("fail to get user from cache: redis.Set(%v)", key))
	}
	var cacheUser = &User{}
	err = json.Unmarshal([]byte(result), cacheUser)
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("json unmarshal error: user(%v)", result))
	}
	return cacheUser, nil
}

func (r *userRepo) setUserToCache(ctx context.Context, user *User, key string) {
	marshal, err := json.Marshal(user)
	if err != nil {
		r.log.Errorf("json marshal error: json.Marshal(%v), error(%v)", user, err)
	}
	err = r.data.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err()
	if err != nil {
		r.log.Errorf("fail to set user to cache: redis.Set(%v), error(%v)", user, err)
	}
}
