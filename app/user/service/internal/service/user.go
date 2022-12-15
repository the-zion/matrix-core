package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
	"github.com/the-zion/matrix-core/app/user/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) GetAccount(ctx context.Context, req *v1.GetAccountReq) (*v1.GetAccountReply, error) {
	user, err := s.uc.GetAccount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetAccountReply{
		Phone:    user.Phone,
		Email:    user.Email,
		Qq:       user.Qq,
		Wechat:   user.Wechat,
		Gitee:    user.Gitee,
		Github:   user.Github,
		Password: user.Password,
	}, nil
}

func (s *UserService) GetProfile(ctx context.Context, req *v1.GetProfileReq) (*v1.GetProfileReply, error) {
	profile, err := s.uc.GetProfile(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileReply{
		Uuid:      profile.Uuid,
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		School:    profile.School,
		Company:   profile.Company,
		Job:       profile.Job,
		Homepage:  profile.Homepage,
		Introduce: profile.Introduce,
		Github:    profile.Github,
		Gitee:     profile.Gitee,
		Created:   profile.Created,
	}, nil
}

func (s *UserService) GetProfileList(ctx context.Context, req *v1.GetProfileListReq) (*v1.GetProfileListReply, error) {
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

func (s *UserService) GetProfileUpdate(ctx context.Context, req *v1.GetProfileUpdateReq) (*v1.GetProfileUpdateReply, error) {
	profile, err := s.uc.GetProfileUpdate(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileUpdateReply{
		Username:  profile.Username,
		Avatar:    profile.Avatar,
		School:    profile.School,
		Company:   profile.Company,
		Job:       profile.Job,
		Homepage:  profile.Homepage,
		Github:    profile.Github,
		Gitee:     profile.Gitee,
		Introduce: profile.Introduce,
		Status:    profile.Status,
	}, nil
}

func (s *UserService) GetUserFollow(ctx context.Context, req *v1.GetUserFollowReq) (*v1.GetUserFollowReply, error) {
	follow, err := s.uc.GetUserFollow(ctx, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserFollowReply{
		Follow: follow,
	}, nil
}

func (s *UserService) GetFollowList(ctx context.Context, req *v1.GetFollowListReq) (*v1.GetFollowListReply, error) {
	followList, err := s.uc.GetFollowList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetFollowListReply{Follow: make([]*v1.GetFollowListReply_Follow, 0, len(followList))}
	for _, item := range followList {
		reply.Follow = append(reply.Follow, &v1.GetFollowListReply_Follow{
			Uuid: item.Follow,
		})
	}
	return reply, nil
}

func (s *UserService) GetFollowListCount(ctx context.Context, req *v1.GetFollowListCountReq) (*v1.GetFollowListCountReply, error) {
	count, err := s.uc.GetFollowListCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetFollowListCountReply{
		Count: count,
	}, nil
}

func (s *UserService) GetFollowedList(ctx context.Context, req *v1.GetFollowedListReq) (*v1.GetFollowedListReply, error) {
	followedList, err := s.uc.GetFollowedList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetFollowedListReply{Follow: make([]*v1.GetFollowedListReply_Follow, 0, len(followedList))}
	for _, item := range followedList {
		reply.Follow = append(reply.Follow, &v1.GetFollowedListReply_Follow{
			Uuid: item.Followed,
		})
	}
	return reply, nil
}

func (s *UserService) GetFollowedListCount(ctx context.Context, req *v1.GetFollowedListCountReq) (*v1.GetFollowedListCountReply, error) {
	count, err := s.uc.GetFollowedListCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetFollowedListCountReply{
		Count: count,
	}, nil
}

func (s *UserService) GetUserFollows(ctx context.Context, req *v1.GetUserFollowsReq) (*v1.GetUserFollowsReply, error) {
	followsMap, err := s.uc.GetUserFollows(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserFollowsReply{
		Follows: followsMap,
	}, nil
}

func (s *UserService) GetUserSearch(ctx context.Context, req *v1.GetUserSearchReq) (*v1.GetUserSearchReply, error) {
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
		})
	}
	reply.Total = total
	return reply, nil
}

func (s *UserService) GetAvatarReview(ctx context.Context, req *v1.GetAvatarReviewReq) (*v1.GetAvatarReviewReply, error) {
	reviewList, err := s.uc.GetAvatarReview(ctx, req.Page, req.Uuid)
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

func (s *UserService) GetCoverReview(ctx context.Context, req *v1.GetCoverReviewReq) (*v1.GetCoverReviewReply, error) {
	reviewList, err := s.uc.GetCoverReview(ctx, req.Page, req.Uuid)
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

func (s *UserService) AvatarIrregular(ctx context.Context, req *v1.AvatarIrregularReq) (*emptypb.Empty, error) {
	err := s.uc.AvatarIrregular(ctx, &biz.ImageReview{
		Uuid:     req.Uuid,
		JobId:    req.JobId,
		Url:      req.Url,
		Label:    req.Label,
		Result:   req.Result,
		Score:    req.Score,
		Category: req.Category,
		SubLabel: req.SubLabel,
		Mode:     "add_avatar_review_db_and_cache",
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) AddAvatarReviewDbAndCache(ctx context.Context, req *v1.AddAvatarReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.uc.AddAvatarReviewDbAndCache(ctx, &biz.ImageReview{
		Uuid:     req.Uuid,
		JobId:    req.JobId,
		Url:      req.Url,
		Label:    req.Label,
		Result:   req.Result,
		Score:    req.Score,
		Category: req.Category,
		SubLabel: req.SubLabel,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) CoverIrregular(ctx context.Context, req *v1.CoverIrregularReq) (*emptypb.Empty, error) {
	err := s.uc.CoverIrregular(ctx, &biz.ImageReview{
		Uuid:     req.Uuid,
		JobId:    req.JobId,
		Url:      req.Url,
		Label:    req.Label,
		Result:   req.Result,
		Score:    req.Score,
		Category: req.Category,
		SubLabel: req.SubLabel,
		Mode:     "add_cover_review_db_and_cache",
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) AddCoverReviewDbAndCache(ctx context.Context, req *v1.AddCoverReviewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.uc.AddCoverReviewDbAndCache(ctx, &biz.ImageReview{
		Uuid:     req.Uuid,
		JobId:    req.JobId,
		Url:      req.Url,
		Label:    req.Label,
		Result:   req.Result,
		Score:    req.Score,
		Category: req.Category,
		SubLabel: req.SubLabel,
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SetProfileUpdate(ctx context.Context, req *v1.SetProfileUpdateReq) (*emptypb.Empty, error) {
	profile := &biz.ProfileUpdate{}
	profile.Uuid = req.Uuid
	profile.Username = req.Username
	profile.School = req.School
	profile.Company = req.Company
	profile.Job = req.Job
	profile.Homepage = req.Homepage
	profile.Github = req.Github
	profile.Gitee = req.Gitee
	profile.Introduce = req.Introduce
	err := s.uc.SetProfileUpdate(ctx, profile)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SetUserFollow(ctx context.Context, req *v1.SetUserFollowReq) (*emptypb.Empty, error) {
	err := s.uc.SetUserFollow(ctx, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) SetFollowDbAndCache(ctx context.Context, req *v1.SetFollowDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.uc.SetFollowDbAndCache(ctx, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) CancelUserFollow(ctx context.Context, req *v1.CancelUserFollowReq) (*emptypb.Empty, error) {
	err := s.uc.CancelUserFollow(ctx, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) CancelFollowDbAndCache(ctx context.Context, req *v1.CancelFollowDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.uc.CancelFollowDbAndCache(ctx, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) ProfileReviewPass(ctx context.Context, req *v1.ProfileReviewPassReq) (*emptypb.Empty, error) {
	err := s.uc.ProfileReviewPass(ctx, req.Uuid, req.Update)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) ProfileReviewNotPass(ctx context.Context, req *v1.ProfileReviewNotPassReq) (*emptypb.Empty, error) {
	err := s.uc.ProfileReviewNotPass(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
