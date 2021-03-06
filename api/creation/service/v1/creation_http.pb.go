// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.3.1

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type CreationHTTPServer interface {
	GetArticleList(context.Context, *GetArticleListReq) (*GetArticleListReply, error)
}

func RegisterCreationHTTPServer(s *http.Server, srv CreationHTTPServer) {
	r := s.Route("/")
	r.GET("/v1/get/article/list", _Creation_GetArticleList0_HTTP_Handler(srv))
}

func _Creation_GetArticleList0_HTTP_Handler(srv CreationHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetArticleListReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/creation.v1.Creation/GetArticleList")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetArticleList(ctx, req.(*GetArticleListReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetArticleListReply)
		return ctx.Result(200, reply)
	}
}

type CreationHTTPClient interface {
	GetArticleList(ctx context.Context, req *GetArticleListReq, opts ...http.CallOption) (rsp *GetArticleListReply, err error)
}

type CreationHTTPClientImpl struct {
	cc *http.Client
}

func NewCreationHTTPClient(client *http.Client) CreationHTTPClient {
	return &CreationHTTPClientImpl{client}
}

func (c *CreationHTTPClientImpl) GetArticleList(ctx context.Context, in *GetArticleListReq, opts ...http.CallOption) (*GetArticleListReply, error) {
	var out GetArticleListReply
	pattern := "/v1/get/article/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/creation.v1.Creation/GetArticleList"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
