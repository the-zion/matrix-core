package service

import (
	"context"
	"fmt"
	"github.com/the-zion/matrix-core/api/message/service/v1"
)

func (s *MessageService) ProfileReview(ctx context.Context, req *v1.ProfileReviewReq) (*v1.ProfileReviewReply, error) {
	fmt.Println(req)
	return &v1.ProfileReviewReply{}, nil
}
