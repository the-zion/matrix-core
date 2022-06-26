package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	userV1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"golang.org/x/sync/singleflight"
	"google.golang.org/protobuf/types/known/emptypb"
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
		log:  log.NewHelper(log.With(logger, "module", "bff/data/user")),
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
	reply, err := r.data.uc.GetCosSessionKey(ctx, &emptypb.Empty{})
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

func (r *userRepo) GetAccount(ctx context.Context, uuid string) (*biz.UserAccount, error) {
	account, err := r.data.uc.GetAccount(ctx, &userV1.GetAccountReq{
		Uuid: uuid,
	})
	if err != nil {
		return nil, err
	}
	return &biz.UserAccount{
		Phone:    account.Phone,
		Email:    account.Email,
		Qq:       account.Qq,
		Wechat:   account.Wechat,
		Weibo:    account.Weibo,
		Github:   account.Github,
		Password: account.Password,
	}, nil
}

func (r *userRepo) GetProfile(ctx context.Context, uuid string) (*biz.UserProfile, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_user_profile_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.uc.GetProfile(ctx, &userV1.GetProfileReq{
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

func (r *userRepo) GetProfileUpdate(ctx context.Context, uuid string) (*biz.UserProfileUpdate, error) {
	pu := &biz.UserProfileUpdate{}
	reply, err := r.data.uc.GetProfileUpdate(ctx, &userV1.GetProfileUpdateReq{
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

func (r *userRepo) SetProfileUpdate(ctx context.Context, profile *biz.UserProfileUpdate) error {
	_, err := r.data.uc.SetProfileUpdate(ctx, &userV1.SetProfileUpdateReq{
		Uuid:      profile.Uuid,
		Username:  profile.Username,
		School:    profile.School,
		Company:   profile.Company,
		Job:       profile.Job,
		Homepage:  profile.Homepage,
		Introduce: profile.Introduce,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetUserPhone(ctx context.Context, uuid, phone, code string) error {
	_, err := r.data.uc.SetUserPhone(ctx, &userV1.SetUserPhoneReq{
		Uuid:  uuid,
		Phone: phone,
		Code:  code,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetUserEmail(ctx context.Context, uuid, email, code string) error {
	_, err := r.data.uc.SetUserEmail(ctx, &userV1.SetUserEmailReq{
		Uuid:  uuid,
		Email: email,
		Code:  code,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetUserPassword(ctx context.Context, uuid, password string) error {
	_, err := r.data.uc.SetUserPassword(ctx, &userV1.SetUserPasswordReq{
		Uuid:     uuid,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) ChangeUserPassword(ctx context.Context, uuid, oldpassword, password string) error {
	_, err := r.data.uc.ChangeUserPassword(ctx, &userV1.ChangeUserPasswordReq{
		Uuid:        uuid,
		Oldpassword: oldpassword,
		Password:    password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) UnbindUserPhone(ctx context.Context, uuid, phone, code string) error {
	_, err := r.data.uc.UnbindUserPhone(ctx, &userV1.UnbindUserPhoneReq{
		Uuid:  uuid,
		Phone: phone,
		Code:  code,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) UnbindUserEmail(ctx context.Context, uuid, email, code string) error {
	_, err := r.data.uc.UnbindUserEmail(ctx, &userV1.UnbindUserEmailReq{
		Uuid:  uuid,
		Email: email,
		Code:  code,
	})
	if err != nil {
		return err
	}
	return nil
}
