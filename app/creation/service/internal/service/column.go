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

func (s *CreationService) SendColumn(ctx context.Context, req *v1.SendColumnReq) (*emptypb.Empty, error) {
	err := s.coc.SendColumn(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SendColumnEdit(ctx context.Context, req *v1.SendColumnEditReq) (*emptypb.Empty, error) {
	err := s.coc.SendColumnEdit(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteColumn(ctx context.Context, req *v1.DeleteColumnReq) (*emptypb.Empty, error) {
	err := s.coc.DeleteColumn(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CreateColumn(ctx context.Context, req *v1.CreateColumnReq) (*emptypb.Empty, error) {
	err := s.coc.CreateColumn(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) EditColumn(ctx context.Context, req *v1.EditColumnReq) (*emptypb.Empty, error) {
	err := s.coc.EditColumn(ctx, req.Id, req.Auth, req.Uuid)
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

func (s *CreationService) EditColumnCosAndSearch(ctx context.Context, req *v1.EditColumnCosAndSearchReq) (*emptypb.Empty, error) {
	err := s.coc.EditColumnCosAndSearch(ctx, req.Id, req.Auth, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteColumnCacheAndSearch(ctx context.Context, req *v1.DeleteColumnCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.coc.DeleteColumnCacheAndSearch(ctx, req.Id, req.Uuid)
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
			Id:   item.ColumnId,
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

func (s *CreationService) GetUserColumnList(ctx context.Context, req *v1.GetUserColumnListReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetUserColumnList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:   item.ColumnId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetUserColumnListVisitor(ctx context.Context, req *v1.GetUserColumnListVisitorReq) (*v1.GetColumnListReply, error) {
	reply := &v1.GetColumnListReply{Column: make([]*v1.GetColumnListReply_Column, 0)}
	columnList, err := s.coc.GetUserColumnListVisitor(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range columnList {
		reply.Column = append(reply.Column, &v1.GetColumnListReply_Column{
			Id:   item.ColumnId,
			Uuid: item.Uuid,
		})
	}
	return reply, nil
}

func (s *CreationService) GetColumnCount(ctx context.Context, req *v1.GetColumnCountReq) (*v1.GetColumnCountReply, error) {
	count, err := s.coc.GetColumnCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetColumnCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetColumnCountVisitor(ctx context.Context, req *v1.GetColumnCountVisitorReq) (*v1.GetColumnCountReply, error) {
	count, err := s.coc.GetColumnCountVisitor(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetColumnCountReply{
		Count: count,
	}, nil
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

func (s *CreationService) GetColumnStatistic(ctx context.Context, req *v1.GetColumnStatisticReq) (*v1.GetColumnStatisticReply, error) {
	columnStatistic, err := s.coc.GetColumnStatistic(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetColumnStatisticReply{
		Uuid:    columnStatistic.Uuid,
		Agree:   columnStatistic.Agree,
		Collect: columnStatistic.Collect,
		View:    columnStatistic.View,
	}, nil
}

func (s *CreationService) ColumnStatisticJudge(ctx context.Context, req *v1.ColumnStatisticJudgeReq) (*v1.ColumnStatisticJudgeReply, error) {
	judge, err := s.coc.ColumnStatisticJudge(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.ColumnStatisticJudgeReply{
		Agree:   judge.Agree,
		Collect: judge.Collect,
	}, nil
}

func (s *CreationService) SetColumnAgree(ctx context.Context, req *v1.SetColumnAgreeReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelColumnAgree(ctx context.Context, req *v1.CancelColumnAgreeReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelColumnCollect(ctx context.Context, req *v1.CancelColumnCollectReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnCollect(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetColumnView(ctx context.Context, req *v1.SetColumnViewReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnView(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) AddColumnIncludes(ctx context.Context, req *v1.AddColumnIncludesReq) (*emptypb.Empty, error) {
	err := s.coc.AddColumnIncludes(ctx, req.Id, req.ArticleId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) DeleteColumnIncludes(ctx context.Context, req *v1.DeleteColumnIncludesReq) (*emptypb.Empty, error) {
	err := s.coc.DeleteColumnIncludes(ctx, req.Id, req.ArticleId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
