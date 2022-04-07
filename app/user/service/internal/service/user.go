package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {
	//return s.ac.Login(ctx, req)
	return nil, nil
}
