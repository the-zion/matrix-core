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

func (s *MessageService) GetMessageSystemNotification(ctx context.Context, req *v1.GetMessageSystemNotificationReq) (*v1.GetMessageSystemNotificationReply, error) {
	reply := &v1.GetMessageSystemNotificationReply{List: make([]*v1.GetMessageSystemNotificationReply_List, 0)}
	notificationList, err := s.mc.GetMessageSystemNotification(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range notificationList {
		reply.List = append(reply.List, &v1.GetMessageSystemNotificationReply_List{
			Id:               item.Id,
			ContentId:        item.ContentId,
			CreatedAt:        item.CreatedAt,
			NotificationType: item.NotificationType,
			Title:            item.Title,
			Uid:              item.Uid,
			Uuid:             item.Uuid,
			Label:            item.Label,
			Result:           item.Result,
			Section:          item.Section,
			Text:             item.Text,
			Comment:          item.Comment,
		})
	}
	return reply, nil
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

func (s *MessageService) RemoveMailBoxSystemNotificationCount(ctx context.Context, req *v1.RemoveMailBoxSystemNotificationCountReq) (*emptypb.Empty, error) {
	err := s.mc.RemoveMailBoxSystemNotificationCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
