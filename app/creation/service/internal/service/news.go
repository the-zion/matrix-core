package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

func (s *CreationService) GetNewsFromTianXing(ctx context.Context, req *v1.GetNewsReq) (*v1.GetNewsReply, error) {
	reply := &v1.GetNewsReply{News: make([]*v1.GetNewsReply_News, 0)}
	newsList, err := s.nc.GetNewsFromTianXing(ctx, req.Page, req.Kind)
	if err != nil {
		return nil, err
	}
	for _, item := range newsList {
		reply.News = append(reply.News, &v1.GetNewsReply_News{
			Id:     item.Id,
			Update: item.Update,
			Title:  item.Title,
			Text:   item.Text,
			Tags:   item.Tags,
			Cover:  item.Cover,
			Url:    item.Url,
		})
	}
	return reply, nil
}
