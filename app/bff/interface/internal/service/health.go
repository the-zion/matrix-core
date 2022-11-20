package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *BffService) GetHealth(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
