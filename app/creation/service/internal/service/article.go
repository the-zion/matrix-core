package service

import (
	"context"
	v1 "github.com/the-zion/matrix-core/api/creation/service/v1"
)

func (s *CreationService) CreateArticleDraft(ctx context.Context, req *v1.CreateArticleDraftReq) (*v1.CreateArticleDraftReply, error) {
	id, err := s.ac.CreateArticleDraft(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	return &v1.CreateArticleDraftReply{
		Id: id,
	}, nil
}
