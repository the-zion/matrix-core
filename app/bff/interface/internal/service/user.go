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

func (s *BffService) GetUserProfile(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserProfileReply, error) {
	userProfile, err := s.uc.GetUserProfile(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserProfileReply{
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

func (s *BffService) GetUserProfileUpdate(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserProfileUpdateReply, error) {
	userProfile, err := s.uc.GetUserProfileUpdate(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserProfileUpdateReply{
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

func (s *BffService) SetUserProfile(ctx context.Context, req *v1.SetUserProfileReq) (*emptypb.Empty, error) {
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

func (s *BffService) ProfileReview(ctx context.Context, req *v1.ProfileReviewReq) (*emptypb.Empty, error) {
	tr := &biz.TextReview{
		Code:         req.JobsDetail.Code,
		Message:      req.JobsDetail.Message,
		JobId:        req.JobsDetail.JobId,
		DataId:       req.JobsDetail.DataId,
		State:        req.JobsDetail.State,
		CreationTime: req.JobsDetail.CreationTime,
		Object:       req.JobsDetail.Object,
		Label:        req.JobsDetail.Label,
		Result:       req.JobsDetail.Result,
		BucketId:     req.JobsDetail.BucketId,
		Region:       req.JobsDetail.Region,
		CosHeaders:   req.JobsDetail.CosHeaders,
	}

	var section []*biz.Section

	for _, item := range req.JobsDetail.Section {
		se := &biz.Section{
			Label:  item.Label,
			Result: item.Result,
			PornInfo: &biz.SectionPornInfo{
				HitFlag:  item.PornInfo.HitFlag,
				Score:    item.PornInfo.Score,
				Keywords: item.PornInfo.Keywords,
			},
			AdsInfo: &biz.SectionAdsInfo{
				HitFlag:  item.AdsInfo.HitFlag,
				Score:    item.AdsInfo.Score,
				Keywords: item.AdsInfo.Keywords,
			},
			IllegalInfo: &biz.SectionIllegalInfo{
				HitFlag:  item.IllegalInfo.HitFlag,
				Score:    item.IllegalInfo.Score,
				Keywords: item.IllegalInfo.Keywords,
			},
			AbuseInfo: &biz.SectionAbuseInfo{
				HitFlag:  item.AbuseInfo.HitFlag,
				Score:    item.AbuseInfo.Score,
				Keywords: item.AbuseInfo.Keywords,
			},
		}
		section = append(section, se)
	}

	tr.Section = section

	err := s.uc.ProfileReview(ctx, tr)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
