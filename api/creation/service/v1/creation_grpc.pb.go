// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: creation/service/v1/creation.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CreationClient is the client API for Creation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CreationClient interface {
	CreateArticleDraft(ctx context.Context, in *CreateArticleDraftReq, opts ...grpc.CallOption) (*CreateArticleDraftReply, error)
}

type creationClient struct {
	cc grpc.ClientConnInterface
}

func NewCreationClient(cc grpc.ClientConnInterface) CreationClient {
	return &creationClient{cc}
}

func (c *creationClient) CreateArticleDraft(ctx context.Context, in *CreateArticleDraftReq, opts ...grpc.CallOption) (*CreateArticleDraftReply, error) {
	out := new(CreateArticleDraftReply)
	err := c.cc.Invoke(ctx, "/creation.v1.Creation/CreateArticleDraft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreationServer is the server API for Creation service.
// All implementations must embed UnimplementedCreationServer
// for forward compatibility
type CreationServer interface {
	CreateArticleDraft(context.Context, *CreateArticleDraftReq) (*CreateArticleDraftReply, error)
	mustEmbedUnimplementedCreationServer()
}

// UnimplementedCreationServer must be embedded to have forward compatible implementations.
type UnimplementedCreationServer struct {
}

func (UnimplementedCreationServer) CreateArticleDraft(context.Context, *CreateArticleDraftReq) (*CreateArticleDraftReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArticleDraft not implemented")
}
func (UnimplementedCreationServer) mustEmbedUnimplementedCreationServer() {}

// UnsafeCreationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CreationServer will
// result in compilation errors.
type UnsafeCreationServer interface {
	mustEmbedUnimplementedCreationServer()
}

func RegisterCreationServer(s grpc.ServiceRegistrar, srv CreationServer) {
	s.RegisterService(&Creation_ServiceDesc, srv)
}

func _Creation_CreateArticleDraft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateArticleDraftReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CreationServer).CreateArticleDraft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/creation.v1.Creation/CreateArticleDraft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreationServer).CreateArticleDraft(ctx, req.(*CreateArticleDraftReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Creation_ServiceDesc is the grpc.ServiceDesc for Creation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Creation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "creation.v1.Creation",
	HandlerType: (*CreationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateArticleDraft",
			Handler:    _Creation_CreateArticleDraft_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "creation/service/v1/creation.proto",
}