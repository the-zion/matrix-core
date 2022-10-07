package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *BffService) GetMailBoxLastTime(ctx context.Context, _ *emptypb.Empty) (*v1.GetMailBoxLastTimeReply, error) {
	mailbox, err := s.mc.GetMailBoxLastTime(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetMailBoxLastTimeReply{
		Time: mailbox.Time,
	}, nil
}

func (s *BffService) SetMailBoxLastTime(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.mc.SetMailBoxLastTime(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
