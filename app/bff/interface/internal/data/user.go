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

func (r *userRepo) GetCosSessionKey(ctx context.Context, uuid string) (*biz.Credentials, error) {
	reply, err := r.data.uc.GetCosSessionKey(ctx, &userV1.GetCosSessionKeyReq{
		Uuid: uuid,
	})
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
}

func (r *userRepo) GetProfileList(ctx context.Context, uuids []string) ([]*biz.UserProfile, error) {
	reply := make([]*biz.UserProfile, 0)
	profileList, err := r.data.uc.GetProfileList(ctx, &userV1.GetProfileListReq{
		Uuids: uuids,
	})
	if err != nil {
		return nil, err
	}

	for _, item := range profileList.Profile {
		reply = append(reply, &biz.UserProfile{
			Uuid:      item.Uuid,
			Username:  item.Username,
			Introduce: item.Introduce,
		})
	}
	return reply, nil
}

func (r *userRepo) GetUserInfo(ctx context.Context, uuid string) (*biz.UserProfile, error) {
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
			Created:   reply.Created,
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

func (r *userRepo) GetUserFollow(ctx context.Context, uuid, userUuid string) (bool, error) {
	reply, err := r.data.uc.GetUserFollow(ctx, &userV1.GetUserFollowReq{
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return false, err
	}
	return reply.Follow, nil
}

func (r *userRepo) GetUserFollows(ctx context.Context, userId string, uuids []string) ([]*biz.Follows, error) {
	reply := make([]*biz.Follows, 0)
	followsList, err := r.data.uc.GetUserFollows(ctx, &userV1.GetUserFollowsReq{
		Uuids:  uuids,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	for _, item := range followsList.Follows {
		reply = append(reply, &biz.Follows{
			Uuid:   item.Uuid,
			Follow: item.FollowJudge,
		})
	}
	return reply, nil
}

func (r *userRepo) GetFollowList(ctx context.Context, page int32, uuid string) ([]*biz.Follow, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_follow_%s", uuid), func() (interface{}, error) {
		reply := make([]*biz.Follow, 0)
		followList, err := r.data.uc.GetFollowList(ctx, &userV1.GetFollowListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range followList.Follow {
			reply = append(reply, &biz.Follow{
				Follow: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Follow), nil
}

func (r *userRepo) GetFollowListCount(ctx context.Context, uuid string) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_follow_count_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.uc.GetFollowListCount(ctx, &userV1.GetFollowListCountReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *userRepo) GetFollowedList(ctx context.Context, page int32, uuid string) ([]*biz.Follow, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_followed_%s", uuid), func() (interface{}, error) {
		reply := make([]*biz.Follow, 0)
		followedList, err := r.data.uc.GetFollowedList(ctx, &userV1.GetFollowedListReq{
			Page: page,
			Uuid: uuid,
		})
		if err != nil {
			return nil, err
		}

		for _, item := range followedList.Follow {
			reply = append(reply, &biz.Follow{
				Followed: item.Uuid,
			})
		}
		return reply, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]*biz.Follow), nil
}

func (r *userRepo) GetFollowedListCount(ctx context.Context, uuid string) (int32, error) {
	result, err, _ := r.sg.Do(fmt.Sprintf("get_followed_count_%s", uuid), func() (interface{}, error) {
		reply, err := r.data.uc.GetFollowedListCount(ctx, &userV1.GetFollowedListCountReq{
			Uuid: uuid,
		})
		if err != nil {
			return 0, err
		}
		return reply.Count, nil
	})
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func (r *userRepo) GetUserSearch(ctx context.Context, page int32, search string) ([]*biz.UserSearch, int32, error) {
	reply := make([]*biz.UserSearch, 0)
	searchReply, err := r.data.uc.GetUserSearch(ctx, &userV1.GetUserSearchReq{
		Page:   page,
		Search: search,
	})
	if err != nil {
		return nil, 0, err
	}
	for _, item := range searchReply.List {
		reply = append(reply, &biz.UserSearch{
			Uuid:      item.Uuid,
			Username:  item.Username,
			Introduce: item.Introduce,
		})
	}
	return reply, searchReply.Total, nil
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

func (r *userRepo) SetUserFollow(ctx context.Context, uuid, userUuid string) error {
	_, err := r.data.uc.SetUserFollow(ctx, &userV1.SetUserFollowReq{
		Uuid:     uuid,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) CancelUserFollow(ctx context.Context, uuid, userUuid string) error {
	_, err := r.data.uc.CancelUserFollow(ctx, &userV1.CancelUserFollowReq{
		Uuid:     uuid,
		UserUuid: userUuid,
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
