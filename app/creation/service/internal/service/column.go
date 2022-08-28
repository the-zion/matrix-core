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

func (s *CreationService) SubscribeColumn(ctx context.Context, req *v1.SubscribeColumnReq) (*emptypb.Empty, error) {
	err := s.coc.SubscribeColumn(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetColumnSubscribeDbAndCache(ctx context.Context, req *v1.SetColumnSubscribeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnSubscribeDbAndCache(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelSubscribeColumn(ctx context.Context, req *v1.CancelSubscribeColumnReq) (*emptypb.Empty, error) {
	err := s.coc.CancelSubscribeColumn(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) CancelColumnSubscribeDbAndCache(ctx context.Context, req *v1.CancelColumnSubscribeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnSubscribeDbAndCache(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SubscribeJudge(ctx context.Context, req *v1.SubscribeJudgeReq) (*v1.SubscribeJudgeReply, error) {
	subscribe, err := s.coc.SubscribeJudge(ctx, req.Id, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.SubscribeJudgeReply{
		Subscribe: subscribe,
	}, nil
}

func (s *CreationService) SendColumn(ctx context.Context, req *v1.SendColumnReq) (*emptypb.Empty, error) {
	err := s.coc.SendColumn(ctx, req.Id, req.Uuid, req.Ip)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SendColumnEdit(ctx context.Context, req *v1.SendColumnEditReq) (*emptypb.Empty, error) {
	err := s.coc.SendColumnEdit(ctx, req.Id, req.Uuid, req.Ip)
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

func (s *CreationService) CreateColumnDbCacheAndSearch(ctx context.Context, req *v1.CreateColumnDbCacheAndSearchReq) (*emptypb.Empty, error) {
	err := s.coc.CreateColumnDbCacheAndSearch(ctx, req.Id, req.Auth, req.Uuid)
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
	columnStatistic, err := s.coc.GetColumnStatistic(ctx, req.Id, req.Uuid)
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

func (s *CreationService) GetSubscribeList(ctx context.Context, req *v1.GetSubscribeListReq) (*v1.GetSubscribeListReply, error) {
	reply := &v1.GetSubscribeListReply{Subscribe: make([]*v1.GetSubscribeListReply_Subscribe, 0)}
	statisticList, err := s.coc.GetSubscribeList(ctx, req.Page, req.Uuid)
	if err != nil {
		return nil, err
	}
	for _, item := range statisticList {
		reply.Subscribe = append(reply.Subscribe, &v1.GetSubscribeListReply_Subscribe{
			Id:   item.ColumnId,
			Uuid: item.AuthorId,
		})
	}
	return reply, nil
}

func (s *CreationService) GetSubscribeListCount(ctx context.Context, req *v1.GetSubscribeListCountReq) (*v1.GetSubscribeListCountReply, error) {
	count, err := s.coc.GetSubscribeListCount(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetSubscribeListCountReply{
		Count: count,
	}, nil
}

func (s *CreationService) GetColumnSubscribes(ctx context.Context, req *v1.GetColumnSubscribesReq) (*v1.GetColumnSubscribesReply, error) {
	reply := &v1.GetColumnSubscribesReply{Subscribes: make([]*v1.GetColumnSubscribesReply_Subscribes, 0)}
	subscribesList, err := s.coc.GetColumnSubscribes(ctx, req.Uuid, req.Ids)
	if err != nil {
		return nil, err
	}
	for _, item := range subscribesList {
		reply.Subscribes = append(reply.Subscribes, &v1.GetColumnSubscribesReply_Subscribes{
			Id:             item.ColumnId,
			SubscribeJudge: item.Status,
		})
	}
	return reply, nil
}

func (s *CreationService) GetColumnSearch(ctx context.Context, req *v1.GetColumnSearchReq) (*v1.GetColumnSearchReply, error) {
	reply := &v1.GetColumnSearchReply{List: make([]*v1.GetColumnSearchReply_List, 0)}
	columnList, total, err := s.coc.GetColumnSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return reply, err
	}
	for _, item := range columnList {
		reply.List = append(reply.List, &v1.GetColumnSearchReply_List{
			Id:        item.Id,
			Tags:      item.Tags,
			Name:      item.Name,
			Uuid:      item.Uuid,
			Introduce: item.Introduce,
			Cover:     item.Cover,
			Update:    item.Update,
		})
	}
	reply.Total = total
	return reply, nil
}

func (s *CreationService) GetUserColumnAgree(ctx context.Context, req *v1.GetUserColumnAgreeReq) (*v1.GetUserColumnAgreeReply, error) {
	agreeMap, err := s.coc.GetUserColumnAgree(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserColumnAgreeReply{
		Agree: agreeMap,
	}, nil
}

func (s *CreationService) GetUserColumnCollect(ctx context.Context, req *v1.GetUserColumnCollectReq) (*v1.GetUserColumnCollectReply, error) {
	collectMap, err := s.coc.GetUserColumnCollect(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserColumnCollectReply{
		Collect: collectMap,
	}, nil
}

func (s *CreationService) GetUserSubscribeColumn(ctx context.Context, req *v1.GetUserSubscribeColumnReq) (*v1.GetUserSubscribeColumnReply, error) {
	collectMap, err := s.coc.GetUserSubscribeColumn(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserSubscribeColumnReply{
		Subscribe: collectMap,
	}, nil
}

func (s *CreationService) SetColumnAgree(ctx context.Context, req *v1.SetColumnAgreeReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnAgree(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetColumnAgreeDbAndCache(ctx context.Context, req *v1.SetColumnAgreeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
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

func (s *CreationService) CancelColumnAgreeDbAndCache(ctx context.Context, req *v1.CancelColumnAgreeDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnAgreeDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetColumnCollect(ctx context.Context, req *v1.SetColumnCollectReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnCollect(ctx, req.Id, req.CollectionsId, req.Uuid, req.UserUuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *CreationService) SetColumnCollectDbAndCache(ctx context.Context, req *v1.SetColumnCollectDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnCollectDbAndCache(ctx, req.Id, req.CollectionsId, req.Uuid, req.UserUuid)
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

func (s *CreationService) CancelColumnCollectDbAndCache(ctx context.Context, req *v1.CancelColumnCollectDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.CancelColumnCollectDbAndCache(ctx, req.Id, req.Uuid, req.UserUuid)
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

func (s *CreationService) SetColumnViewDbAndCache(ctx context.Context, req *v1.SetColumnViewDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.SetColumnViewDbAndCache(ctx, req.Id, req.Uuid)
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

func (s *CreationService) AddColumnIncludesDbAndCache(ctx context.Context, req *v1.AddColumnIncludesDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.AddColumnIncludesDbAndCache(ctx, req.Id, req.ArticleId, req.Uuid)
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

func (s *CreationService) DeleteColumnIncludesDbAndCache(ctx context.Context, req *v1.DeleteColumnIncludesDbAndCacheReq) (*emptypb.Empty, error) {
	err := s.coc.DeleteColumnIncludesDbAndCache(ctx, req.Id, req.ArticleId, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
