package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

func (s *CreationService) GetNews(ctx context.Context, req *v1.GetNewsReq) (*v1.GetNewsReply, error) {
	newsList, err := s.nc.GetNews(ctx, req.Page)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetNewsReply{News: make([]*v1.GetNewsReply_News, 0, len(newsList))}
	for _, item := range newsList {
		reply.News = append(reply.News, &v1.GetNewsReply_News{
			Id:     item.Id,
			Update: item.Update,
			Title:  item.Title,
			Author: item.Author,
			Text:   item.Text,
			Tags:   item.Tags,
			Cover:  item.Cover,
			Url:    item.Url,
		})
	}
	return reply, nil
}

func (s *CreationService) GetNewsSearch(ctx context.Context, req *v1.GetNewsSearchReq) (*v1.GetNewsSearchReply, error) {
	newsList, total, err := s.nc.GetNewsSearch(ctx, req.Page, req.Search, req.Time)
	if err != nil {
		return nil, err
	}
	reply := &v1.GetNewsSearchReply{List: make([]*v1.GetNewsSearchReply_List, 0, len(newsList))}
	for _, item := range newsList {
		reply.List = append(reply.List, &v1.GetNewsSearchReply_List{
			Id:     item.Id,
			Tags:   item.Tags,
			Title:  item.Title,
			Author: item.Author,
			Update: item.Update,
			Url:    item.Url,
		})
	}
	reply.Total = total
	return reply, nil
}
