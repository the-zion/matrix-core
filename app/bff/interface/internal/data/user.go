package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	userV1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
	sg   *singleflight.Group
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "info/data/user")),
		sg:   &singleflight.Group{},
	}
}

func (r *userRepo) UserRegister(ctx context.Context, email, password, code string) error {
	_, err := r.data.uc.UserRegister(ctx, &userV1.UserRegisterReq{
		Email:    email,
		Password: password,
		Code:     code,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) LoginByPassword(ctx context.Context, account, password, mode string) (string, error) {
	reply, err := r.data.uc.LoginByPassword(ctx, &userV1.LoginByPasswordReq{
		Account:  account,
		Password: password,
		Mode:     mode,
	})
	if err != nil {
		return "", err
	}
	return reply.Token, nil
}

func (r *userRepo) LoginByCode(ctx context.Context, phone, code string) (string, error) {
	reply, err := r.data.uc.LoginByCode(ctx, &userV1.LoginByCodeReq{
		Phone: phone,
		Code:  code,
	})
	if err != nil {
		return "", err
	}
	return reply.Token, nil
}

func (r *userRepo) LoginPasswordReset(ctx context.Context, account, password, code, mode string) error {
	_, err := r.data.uc.LoginPasswordReset(ctx, &userV1.LoginPasswordResetReq{
		Account:  account,
		Password: password,
		Code:     code,
		Mode:     mode,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SendPhoneCode(ctx context.Context, template, phone string) error {
	_, err := r.data.uc.SendPhoneCode(ctx, &userV1.SendPhoneCodeReq{
		Template: template,
		Phone:    phone,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SendEmailCode(ctx context.Context, template, email string) error {
	_, err := r.data.uc.SendEmailCode(ctx, &userV1.SendEmailCodeReq{
		Template: template,
		Email:    email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) GetCosSessionKey(ctx context.Context) (*biz.Credentials, error) {
	reply, err := r.data.uc.GetCosSessionKey(ctx, &userV1.GetCosSessionKeyReq{})
	if err != nil {
		return nil, err
	}
	return &biz.Credentials{
		TmpSecretKey: reply.TmpSecretKey,
		TmpSecretID:  reply.TmpSecretId,
		SessionToken: reply.SessionToken,
		StartTime:    reply.StartTime,
		ExpiredTime:  reply.ExpiredTime,
	}, nil
}

func (r *userRepo) GetUserProfile(ctx context.Context, uuid string) (*biz.UserProfile, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_profile_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.uc.GetUserProfile(ctx, &userV1.GetUserProfileReq{
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}
		return &biz.UserProfile{
			Uuid:      reply.Uuid,
			Username:  reply.Username,
			Avatar:    reply.Avatar,
			School:    reply.School,
			Company:   reply.Company,
			Job:       reply.Job,
			Homepage:  reply.Homepage,
			Introduce: reply.Introduce,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*biz.UserProfile), nil
}

func (r *userRepo) GetUserProfileUpdate(ctx context.Context, uuid string) (*biz.UserProfileUpdate, error) {
	pu := &biz.UserProfileUpdate{}
	reply, err := r.data.uc.GetUserProfileUpdate(ctx, &userV1.GetUserProfileUpdateReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	pu.Username = reply.Username
	pu.Avatar = reply.Avatar
	pu.School = reply.School
	pu.Company = reply.Company
	pu.Job = reply.Job
	pu.Homepage = reply.Homepage
	pu.Introduce = reply.Introduce
	pu.Status = reply.Status
	return pu, nil
}
