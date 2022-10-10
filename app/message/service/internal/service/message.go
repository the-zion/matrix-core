package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/message/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *MessageService) GetMailBoxLastTime(ctx context.Context, req *v1.GetMailBoxLastTimeReq) (*v1.GetMailBoxLastTimeReply, error) {
	mailbox, err := s.mc.GetMailBoxLastTime(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetMailBoxLastTimeReply{
		Time: mailbox.Time,
	}, nil
}

func (s *MessageService) GetMessageNotification(ctx context.Context, req *v1.GetMessageNotificationReq) (*v1.GetMessageNotificationReply, error) {
	notification, err := s.mc.GetMessageNotification(ctx, req.Uuid, req.Follows)
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

func (s *MessageService) SetMailBoxLastTime(ctx context.Context, req *v1.SetMailBoxLastTimeReq) (*emptypb.Empty, error) {
	err := s.mc.SetMailBoxLastTime(ctx, req.Uuid, req.Time)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) RemoveMailBoxCommentCount(ctx context.Context, req *v1.RemoveMailBoxCommentCountReq) (*emptypb.Empty, error) {
	err := s.mc.RemoveMailBoxCommentCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MessageService) RemoveMailBoxSubCommentCount(ctx context.Context, req *v1.RemoveMailBoxSubCommentCountReq) (*emptypb.Empty, error) {
	err := s.mc.RemoveMailBoxSubCommentCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
