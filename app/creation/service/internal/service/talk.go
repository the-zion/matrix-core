package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CreationService) GetLastTalkDraft(ctx context.Context, req *v1.GetLastTalkDraftReq) (*v1.GetLastTalkDraftReply, error) {
	draft, err := s.tc.GetLastTalkDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastTalkDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *CreationService) GetTalkList(ctx context.Context, req *v1.GetTalkListReq) (*v1.GetTalkListReply, error) {
	reply := &v1.GetTalkListReply{Talk: make([]*v1.GetTalkListReply_Talk, 0)}
	talkList, err := s.tc.GetTalkList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range talkList {
		reply.Talk = append(reply.Talk, &v1.GetTalkListReply_Talk{
			Id:   item.TalkId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) CreateTalkDraft(ctx context.Context, req *v1.CreateTalkDraftReq) (*v1.CreateTalkDraftReply, error) {
	id, err := s.tc.CreateTalkDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateTalkDraftReply{
		Id: id,
	}, nil
}

func (s *CreationService) SendTalk(ctx context.Context, req *v1.SendTalkReq) (*emptypb.Empty, error) {
	err := s.tc.SendTalk(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SendTalkEdit(ctx context.Context, req *v1.SendTalkEditReq) (*emptypb.Empty, error) {
	err := s.tc.SendTalkEdit(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateTalk(ctx context.Context, req *v1.CreateTalkReq) (*emptypb.Empty, error) {
	err := s.tc.CreateTalk(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditTalk(ctx context.Context, req *v1.EditTalkReq) (*emptypb.Empty, error) {
	err := s.tc.EditTalk(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateTalkCacheAndSearch(ctx context.Context, req *v1.CreateTalkCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.tc.CreateTalkCacheAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditTalkCosAndSearch(ctx context.Context, req *v1.EditTalkCosAndSearchReq) (*emptypb.Empty, error) {
	err := s.tc.EditTalkCosAndSearch(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
