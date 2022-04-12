package service

import (
	"context"
	v1 "github.com/Cube-v2/cube-core/api/user/service/v1"
)

func (s *UserService) SendCode(ctx context.Context, req *v1.SendCodeReq) (*v1.SendCodeReply, error) {
	return s.uc.SendCode(ctx, req)
}
