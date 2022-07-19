package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CreationService) GetLastColumnDraft(ctx context.Context, req *v1.GetLastColumnDraftReq) (*v1.GetLastColumnDraftReply, error) {
	draft, err := s.coc.GetLastColumnDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetLastColumnDraftReply{
		Id:     draft.Id,
		Status: draft.Status,
	}, nil
}

func (s *CreationService) CreateColumnDraft(ctx context.Context, req *v1.CreateColumnDraftReq) (*v1.CreateColumnDraftReply, error) {
	id, err := s.coc.CreateColumnDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateColumnDraftReply{
		Id: id,
	}, nil
}

func (s *CreationService) CreateColumn(ctx context.Context, req *v1.CreateColumnReq) (*emptypb.Empty, error) {
	err := s.coc.CreateColumn(ctx, req.Uuid, req.Name, req.Introduce, req.Cover, req.Auth)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateColumnCacheAndSearch(ctx context.Context, req *v1.CreateColumnCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.coc.CreateColumnCacheAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) GetColumnList(ctx context.Context, req *v1.GetColumnListReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetColumnList(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:   item.Id,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetColumnListHot(ctx context.Context, req *v1.GetColumnListHotReq) (*v1.GetColumnListHotReply, error) {
	reply := &v1.GetColumnListHotReply{Column: make([]*v1.GetColumnListHotReply_Column, 0)}
	columnList, err := s.coc.GetColumnListHot(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListHotReply_Column{
			Id:   item.ColumnId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetColumnListStatistic(ctx context.Context, req *v1.GetColumnListStatisticReq) (*v1.GetColumnListStatisticReply, error) {
	reply := &v1.GetColumnListStatisticReply{Count: make([]*v1.GetColumnListStatisticReply_Count, 0)}
	statisticList, err := s.coc.GetColumnListStatistic(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Count = append(reply.Count, &v1.GetColumnListStatisticReply_Count{
			Id:      item.ColumnId,
			Agree:   item.Agree,
			Collect: item.Collect,
			View:    item.View,
		})
	}
	return reply, nil
}
