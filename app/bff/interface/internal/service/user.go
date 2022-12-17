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

func (s *BffService) LoginByWeChat(ctx context.Context, req *v1.LoginByWeChatReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByWechat(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *BffService) LoginByQQ(ctx context.Context, req *v1.LoginByQQReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByQQ(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

func (s *BffService) LoginByGithub(ctx context.Context, req *v1.LoginByGithubReq) (*v1.LoginReply, error) {
	github, err := s.uc.LoginByGithub(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: github.Token,
	}, nil
}

func (s *BffService) LoginByGitee(ctx context.Context, req *v1.LoginByGiteeReq) (*v1.LoginReply, error) {
	token, err := s.uc.LoginByGitee(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
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
		Gitee:    userAccount.Gitee,
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
		Github:    userProfile.Github,
		Gitee:     userProfile.Gitee,
		Introduce: userProfile.Introduce,
	}, nil
}

func (s *BffService) GetProfileList(ctx context.Context, req *v1.GetProfileListReq) (*v1.GetProfileListReply, error) {
	profileList, err := s.uc.GetProfileList(ctx, req.Uuids)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetProfileListReply{Profile: make([]*v1.GetProfileListReply_Profile, 0, len(profileList))}
	for _, item := range profileList {
		reply.Profile = append(reply.Profile, &v1.GetProfileListReply_Profile{
			Uuid:      item.Uuid,
			Username:  item.Username,
			Introduce: item.Introduce,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserInfo(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserInfoReply, error) {
	userProfile, err := s.uc.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserInfoReply{
		Username:    userProfile.Username,
		School:      userProfile.School,
		Company:     userProfile.Company,
		Job:         userProfile.Job,
		Gitee:       userProfile.Gitee,
		Github:      userProfile.Github,
		Homepage:    userProfile.Homepage,
		Introduce:   userProfile.Introduce,
		Created:     userProfile.Created,
		Score:       userProfile.Score,
		Agree:       userProfile.Agree,
		Collect:     userProfile.Collect,
		View:        userProfile.View,
		Follow:      userProfile.Follow,
		Followed:    userProfile.Followed,
		Article:     userProfile.Article,
		Column:      userProfile.Column,
		Talk:        userProfile.Talk,
		Collections: userProfile.Collections,
		Subscribe:   userProfile.Subscribe,
	}, nil
}

func (s *BffService) GetUserInfoVisitor(ctx context.Context, req *v1.GetUserInfoVisitorReq) (*v1.GetUserInfoReply, error) {
	userProfile, err := s.uc.GetUserInfoVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserInfoReply{
		Username:    userProfile.Username,
		School:      userProfile.School,
		Company:     userProfile.Company,
		Job:         userProfile.Job,
		Homepage:    userProfile.Homepage,
		Github:      userProfile.Github,
		Gitee:       userProfile.Gitee,
		Introduce:   userProfile.Introduce,
		Created:     userProfile.Created,
		Score:       userProfile.Score,
		Agree:       userProfile.Agree,
		Collect:     userProfile.Collect,
		View:        userProfile.View,
		Follow:      userProfile.Follow,
		Followed:    userProfile.Followed,
		Article:     userProfile.Article,
		Column:      userProfile.Column,
		Talk:        userProfile.Talk,
		Collections: userProfile.Collections,
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
		Github:    userProfile.Github,
		Gitee:     userProfile.Gitee,
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
	followList, err := s.uc.GetFollowList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetFollowListReply{Follow: make([]*v1.GetFollowListReply_Follow, 0, len(followList))}
	for _, item := range followList {
		reply.Follow = append(reply.Follow, &v1.GetFollowListReply_Follow{
			Uuid:      item.Follow,
			Username:  item.Username,
			Introduce: item.Introduce,
			Agree:     item.Agree,
			View:      item.View,
			Follow:    item.FollowNum,
			Followed:  item.FollowedNum,
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
	followedList, err := s.uc.GetFollowedList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetFollowedListReply{Follow: make([]*v1.GetFollowedListReply_Follow, 0, len(followedList))}
	for _, item := range followedList {
		reply.Follow = append(reply.Follow, &v1.GetFollowedListReply_Follow{
			Uuid:      item.Followed,
			Username:  item.Username,
			Introduce: item.Introduce,
			Agree:     item.Agree,
			View:      item.View,
			Follow:    item.FollowNum,
			Followed:  item.FollowedNum,
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

func (s *BffService) GetUserFollows(ctx context.Context, _ *emptypb.Empty) (*v1.GetUserFollowsReply, error) {
	followsMap, err := s.uc.GetUserFollows(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserFollowsReply{
		Follows: followsMap,
	}, nil
}

func (s *BffService) GetTimeLineUsers(ctx context.Context, _ *emptypb.Empty) (*v1.GetTimeLineUsersReply, error) {
	followList, err := s.uc.GetTimeLineUsers(ctx)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetTimeLineUsersReply{Follows: make([]*v1.GetTimeLineUsersReply_Follows, 0, len(followList))}
	for _, item := range followList {
		reply.Follows = append(reply.Follows, &v1.GetTimeLineUsersReply_Follows{
			Uuid:     item.Uuid,
			Username: item.Username,
		})
	}
	return reply, nil
}

func (s *BffService) GetUserSearch(ctx context.Context, req *v1.GetUserSearchReq) (*v1.GetUserSearchReply, error) {
	userList, total, err := s.uc.GetUserSearch(ctx, req.Page, req.Search)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetUserSearchReply{List: make([]*v1.GetUserSearchReply_List, 0, len(userList))}
	for _, item := range userList {
		reply.List = append(reply.List, &v1.GetUserSearchReply_List{
			Uuid:      item.Uuid,
			Username:  item.Username,
			Introduce: item.Introduce,
			Agree:     item.Agree,
			View:      item.View,
			Follow:    item.FollowNum,
			Followed:  item.FollowedNum,
		})
	}
	reply.Total = total
	return reply, nil
}

func (s *BffService) GetAvatarReview(ctx context.Context, req *v1.GetAvatarReviewReq) (*v1.GetAvatarReviewReply, error) {
	reviewList, err := s.uc.GetAvatarReview(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetAvatarReviewReply{Review: make([]*v1.GetAvatarReviewReply_Review, 0, len(reviewList))}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetAvatarReviewReply_Review{
			Id:       item.Id,
			Uuid:     item.Uuid,
			CreateAt: item.CreateAt,
			JobId:    item.JobId,
			Url:      item.Url,
			Label:    item.Label,
			Result:   item.Result,
			Score:    item.Score,
			Category: item.Category,
			SubLabel: item.SubLabel,
		})
	}
	return reply, nil
}

func (s *BffService) GetCoverReview(ctx context.Context, req *v1.GetCoverReviewReq) (*v1.GetCoverReviewReply, error) {
	reviewList, err := s.uc.GetCoverReview(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetCoverReviewReply{Review: make([]*v1.GetCoverReviewReply_Review, 0, len(reviewList))}
	for _, item := range reviewList {
		reply.Review = append(reply.Review, &v1.GetCoverReviewReply_Review{
			Id:       item.Id,
			Uuid:     item.Uuid,
			CreateAt: item.CreateAt,
			JobId:    item.JobId,
			Url:      item.Url,
			Label:    item.Label,
			Result:   item.Result,
			Score:    item.Score,
			Category: item.Category,
			SubLabel: item.SubLabel,
		})
	}
	return reply, nil
}

func (s *BffService) SetProfileUpdate(ctx context.Context, req *v1.SetProfileUpdateReq) (*emptypb.Empty, error) {
	profile := &biz.UserProfileUpdate{}
	profile.Username = req.Username
	profile.School = req.School
	profile.Company = req.Company
	profile.Job = req.Job
	profile.Homepage = req.Homepage
	profile.Github = req.Github
	profile.Gitee = req.Gitee
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

func (s *BffService) SetUserWechat(ctx context.Context, req *v1.SetUserWechatReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserWechat(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserQQ(ctx context.Context, req *v1.SetUserQQReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserQQ(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserGitee(ctx context.Context, req *v1.SetUserGiteeReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserGitee(ctx, req.Code, req.RedirectUrl)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) SetUserGithub(ctx context.Context, req *v1.SetUserGithubReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserGithub(ctx, req.Code, req.RedirectUrl)
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

func (s *BffService) UnbindUserPhone(ctx context.Context, req *v1.UnbindUserAccountReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserPhone(ctx, req.Phone, req.Email, req.Account, req.Password, req.Code, req.Choose, req.Mode, req.RedirectUri)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserEmail(ctx context.Context, req *v1.UnbindUserAccountReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserEmail(ctx, req.Phone, req.Email, req.Account, req.Password, req.Code, req.Choose, req.Mode, req.RedirectUri)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserWechat(ctx context.Context, req *v1.UnbindUserAccountReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserWechat(ctx, req.Phone, req.Email, req.Account, req.Password, req.Code, req.Choose, req.Mode, req.RedirectUri)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserQQ(ctx context.Context, req *v1.UnbindUserAccountReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserQQ(ctx, req.Phone, req.Email, req.Account, req.Password, req.Code, req.Choose, req.Mode, req.RedirectUri)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserGitee(ctx context.Context, req *v1.UnbindUserAccountReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserGitee(ctx, req.Phone, req.Email, req.Account, req.Password, req.Code, req.Choose, req.Mode, req.RedirectUri)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *BffService) UnbindUserGithub(ctx context.Context, req *v1.UnbindUserAccountReq) (*emptypb.Empty, error) {
	err := s.uc.UnbindUserGithub(ctx, req.Phone, req.Email, req.Account, req.Password, req.Code, req.Choose, req.Mode, req.RedirectUri)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
