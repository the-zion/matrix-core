package service

import (
	"context"
	"github.com/the-zion/matrix-core/api/bff/interface/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *BffService) GetMessageNotification(ctx context.Context, _ *emptypb.Empty) (*v1.GetMessageNotificationReply, error) {
	notification, err := s.mc.GetMessageNotification(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetMessageNotificationReply{
		Timeline:   notification.Timeline,
		Comment:    notification.Comment,
		SubComment: notification.SubComment,
		System:     notification.SystemNotification,
	}, nil
}

func (s *BffService) GetMailBoxLastTime(ctx context.Context, _ *emptypb.Empty) (*v1.GetMailBoxLastTimeReply, error) {
	mailbox, err := s.mc.GetMailBoxLastTime(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetMailBoxLastTimeReply{
		Time: mailbox.Time,
	}, nil
}

func (s *BffService) SetMailBoxLastTime(ctx context.Context, req *v1.SetMailBoxLastTimeReq) (*emptypb.Empty, error) {
	err := s.mc.SetMailBoxLastTime(ctx, req.Time)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
