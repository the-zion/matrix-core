package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) SendCode(ctx context.Context, req *v1.SendCodeReq) (*v1.SendCodeReply, error) {
	return s.uc.SendCode(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.GetUserReply, error) {
	return s.uc.GetUser(ctx, req)
}
