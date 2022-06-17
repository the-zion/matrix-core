package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/user/service/v1"
)

func (s *UserService) SetUserProfile(ctx context.Context, req *v1.SetUserProfileReq) (*v1.SetUserProfileReply, error) {
	//id := GetUserIdFromContext(ctx)
	//err := s.pc.SetUserProfile(ctx, id, req.Sex, req.Introduce, req.Industry, req.Address, req.Profile, req.Tag)
	//if err != nil {
	//	return nil, err
	//}
	//return &v1.SetUserProfileReply{
	//	Success: true,
	//}, nil
	return nil, nil
}
