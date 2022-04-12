package data

import (
	"context"
	"github.com/Cube-v2/cube-core/app/user/service/internal/biz"
	"github.com/Cube-v2/cube-core/app/user/service/internal/pkg/util"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"time"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

type Profile struct {
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}
}

func (r *userRepo) FindByUserAccount(ctx context.Context, account, mode string) (*biz.Auth, error) {
	user := &User{}
	result := r.data.db.WithContext(ctx).Where(mode+" = ?", account).First(user)
	if result.Error != nil {
		return nil, biz.ErrUserNotFound
	}
	return &biz.Auth{
		Id: int64(user.Model.ID),
	}, nil
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
	err = r.setUserCodeCache(ctx, key, code)
	if err != nil {
		return "", biz.ErrSendCodeError
	}
	return code, nil
}

func (r *userRepo) sendPhoneCode(template int64, phone, code string) error {
	request := r.data.phoneCodeCli.request
	client := r.data.phoneCodeCli.client
	request.TemplateId = common.StringPtr(util.GetPhoneTemplate(template))
	request.TemplateParamSet = common.StringPtrs([]string{code})
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})
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
		panic(err)
	}
	return nil
}

func (r *userRepo) VerifyCode(ctx context.Context, account, code, mode string) error {
	key := mode + "_" + account
	codeInCache, err := r.getUserCodeCache(ctx, key)
	if err != nil || code != codeInCache {
		return biz.ErrCodeError
	}
	return nil
}

func (r *userRepo) VerifyPassword(ctx context.Context, id int64, password string) error {
	user := &User{}
	result := r.data.db.WithContext(ctx).Where("id = ? AND password = ?", id, password).First(&user)
	if result.Error != nil {
		return biz.ErrPasswordError
	}
	return nil
}

func (r *userRepo) setUserCodeCache(ctx context.Context, key, code string) error {
	err := r.data.redisCli.Set(ctx, key, code, time.Minute*5).Err()
	if err != nil {
		r.log.Errorf("fail to set code cache:redis.Set(%v) error(%v)", key, err)
		return err
	}
	return nil
}

func (r *userRepo) getUserCodeCache(ctx context.Context, key string) (string, error) {
	code, err := r.data.redisCli.Get(ctx, key).Result()
	if err != nil {
		r.log.Errorf("fail to get code cache:redis.Get(%v) error(%v)", key, err)
		return "", err
	}
	return code, nil
}
