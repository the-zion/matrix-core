// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.3
// - protoc             v3.20.0
// source: achievement/service/v1/achievement.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAchievementGetHealth = "/achievement.v1.Achievement/GetHealth"

type AchievementHTTPServer interface {
	GetHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

func RegisterAchievementHTTPServer(s *http.Server, srv AchievementHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/get/achievement/health", _Achievement_GetHealth0_HTTP_Handler(srv))
}

func _Achievement_GetHealth0_HTTP_Handler(srv AchievementHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAchievementGetHealth)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetHealth(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type AchievementHTTPClient interface {
	GetHealth(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type AchievementHTTPClientImpl struct {
	cc *http.Client
}

func NewAchievementHTTPClient(client *http.Client) AchievementHTTPClient {
	return &AchievementHTTPClientImpl{client}
}

func (c *AchievementHTTPClientImpl) GetHealth(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/get/achievement/health"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAchievementGetHealth))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
