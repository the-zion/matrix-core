package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"github.com/the-zion/matrix-core/app/bff/interface/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *BffService) UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*emptypb.Empty, error) {
	err := s.uc.UserRegister(ctx, req.Email, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) LoginByPassword(ctx context.Context, req *v1.LoginByPasswordReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByPassword(ctx, req.Account, req.Password, req.Mode)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *BffService) LoginByCode(ctx context.Context, req *v1.LoginByCodeReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByCode(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *BffService) LoginPasswordReset(ctx context.Context, req *v1.LoginPasswordResetReq) (*emptypb.Empty, error) {
	err := s.uc.LoginPasswordReset(ctx, req.Account, req.Password, req.Code, req.Mode)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendPhoneCode(ctx context.Context, req *v1.SendPhoneCodeReq) (*emptypb.Empty, error) {
	err := s.uc.SendPhoneCode(ctx, req.Template, req.Phone)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SendEmailCode(ctx context.Context, req *v1.SendEmailCodeReq) (*emptypb.Empty, error) {
	err := s.uc.SendEmailCode(ctx, req.Template, req.Email)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) GetCosSessionKey(ctx context.Context, _ *emptypb.Empty) (*v1.GetCosSessionKeyReply, error) {
	credentials, err := s.uc.GetCosSessionKey(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetCosSessionKeyReply{
		TmpSecretId:  credentials.TmpSecretID,
		TmpSecretKey: credentials.TmpSecretKey,
		SessionToken: credentials.SessionToken,
		StartTime:    credentials.StartTime,
		ExpiredTime:  credentials.ExpiredTime,
	}, nil
}

func (s *BffService) GetAccount(ctx context.Context, _ *emptypb.Empty) (*v1.GetAccountReply, error) {
	userAccount, err := s.uc.GetAccount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetAccountReply{
		Phone:    userAccount.Phone,
		Email:    userAccount.Email,
		Qq:       userAccount.Qq,
		Wechat:   userAccount.Wechat,
		Weibo:    userAccount.Weibo,
		Github:   userAccount.Github,
		Password: userAccount.Password,
	}, nil
}

func (s *BffService) GetProfile(ctx context.Context, _ *emptypb.Empty) (*v1.GetProfileReply, error) {
	userProfile, err := s.uc.GetProfile(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileReply{
		Uuid:      userProfile.Uuid,
		Username:  userProfile.Username,
		Avatar:    userProfile.Avatar,
		School:    userProfile.School,
		Company:   userProfile.Company,
		Job:       userProfile.Job,
		Homepage:  userProfile.Homepage,
		Introduce: userProfile.Introduce,
	}, nil
}

func (s *BffService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoReq) (*v1.GetUserInfoReply, error) {
	userProfile, err := s.uc.GetUserInfo(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserInfoReply{
		Username:  userProfile.Username,
		Avatar:    userProfile.Avatar,
		School:    userProfile.School,
		Company:   userProfile.Company,
		Job:       userProfile.Job,
		Homepage:  userProfile.Homepage,
		Introduce: userProfile.Introduce,
		Created:   userProfile.Created,
	}, nil
}

func (s *BffService) GetProfileUpdate(ctx context.Context, _ *emptypb.Empty) (*v1.GetProfileUpdateReply, error) {
	userProfile, err := s.uc.GetProfileUpdate(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileUpdateReply{
		Username:  userProfile.Username,
		Avatar:    userProfile.Avatar,
		School:    userProfile.School,
		Company:   userProfile.Company,
		Job:       userProfile.Job,
		Homepage:  userProfile.Homepage,
		Introduce: userProfile.Introduce,
		Status:    userProfile.Status,
	}, nil
}

func (s *BffService) GetUserFollow(ctx context.Context, req *v1.GetUserFollowReq) (*v1.GetUserFollowReply, error) {
	follow, err := s.uc.GetUserFollow(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserFollowReply{
		Follow: follow,
	}, nil
}

func (s *BffService) GetFollowList(ctx context.Context, req *v1.GetFollowListReq) (*v1.GetFollowListReply, error) {
	reply := &v1.GetFollowListReply{Follow: make([]*v1.GetFollowListReply_Follow, 0)}
	followList, err := s.uc.GetFollowList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range followList {
		reply.Follow = append(reply.Follow, &v1.GetFollowListReply_Follow{
			Uuid: item.Follow,
		})
	}
	return reply, nil
}

func (s *BffService) GetFollowListCount(ctx context.Context, req *v1.GetFollowListCountReq) (*v1.GetFollowListCountReply, error) {
	count, err := s.uc.GetFollowListCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetFollowListCountReply{
		Count: count,
	}, nil
}

func (s *BffService) GetFollowedList(ctx context.Context, req *v1.GetFollowedListReq) (*v1.GetFollowedListReply, error) {
	reply := &v1.GetFollowedListReply{Follow: make([]*v1.GetFollowedListReply_Follow, 0)}
	followedList, err := s.uc.GetFollowList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range followedList {
		reply.Follow = append(reply.Follow, &v1.GetFollowedListReply_Follow{
			Uuid: item.Followed,
		})
	}
	return reply, nil
}

func (s *BffService) GetFollowedListCount(ctx context.Context, req *v1.GetFollowedListCountReq) (*v1.GetFollowedListCountReply, error) {
	count, err := s.uc.GetFollowedListCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetFollowedListCountReply{
		Count: count,
	}, nil
}

func (s *BffService) SetProfileUpdate(ctx context.Context, req *v1.SetProfileUpdateReq) (*emptypb.Empty, error) {
	profile := &biz.UserProfileUpdate{}
	profile.Username = req.Username
	profile.School = req.School
	profile.Company = req.Company
	profile.Job = req.Job
	profile.Homepage = req.Homepage
	profile.Introduce = req.Introduce
	err := s.uc.SetUserProfile(ctx, profile)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserPhone(ctx context.Context, req *v1.SetUserPhoneReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserPhone(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserEmail(ctx context.Context, req *v1.SetUserEmailReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserEmail(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserPassword(ctx context.Context, req *v1.SetUserPasswordReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserPassword(ctx, req.Password)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserFollow(ctx context.Context, req *v1.SetUserFollowReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserFollow(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) CancelUserFollow(ctx context.Context, req *v1.CancelUserFollowReq) (*emptypb.Empty, error) {
	err := s.uc.CancelUserFollow(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) ChangeUserPassword(ctx context.Context, req *v1.ChangeUserPasswordReq) (*emptypb.Empty, error) {
	err := s.uc.ChangeUserPassword(ctx, req.Oldpassword, req.Password)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserPhone(ctx context.Context, req *v1.UnbindUserPhoneReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserPhone(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserEmail(ctx context.Context, req *v1.UnbindUserEmailReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserEmail(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
