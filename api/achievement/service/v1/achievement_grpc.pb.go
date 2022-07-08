// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: achievement/service/v1/achievement.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AchievementClient is the client API for Achievement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AchievementClient interface {
	SetAchievementAgree(ctx context.Context, in *SetAchievementAgreeReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetAchievementView(ctx context.Context, in *SetAchievementViewReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetAchievementCollect(ctx context.Context, in *SetAchievementCollectReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type achievementClient struct {
	cc grpc.ClientConnInterface
}

func NewAchievementClient(cc grpc.ClientConnInterface) AchievementClient {
	return &achievementClient{cc}
}

func (c *achievementClient) SetAchievementAgree(ctx context.Context, in *SetAchievementAgreeReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/achievement.v1.Achievement/SetAchievementAgree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *achievementClient) SetAchievementView(ctx context.Context, in *SetAchievementViewReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/achievement.v1.Achievement/SetAchievementView", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *achievementClient) SetAchievementCollect(ctx context.Context, in *SetAchievementCollectReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/achievement.v1.Achievement/SetAchievementCollect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *achievementClient) GetHealth(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/achievement.v1.Achievement/GetHealth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AchievementServer is the server API for Achievement service.
// All implementations must embed UnimplementedAchievementServer
// for forward compatibility
type AchievementServer interface {
	SetAchievementAgree(context.Context, *SetAchievementAgreeReq) (*emptypb.Empty, error)
	SetAchievementView(context.Context, *SetAchievementViewReq) (*emptypb.Empty, error)
	SetAchievementCollect(context.Context, *SetAchievementCollectReq) (*emptypb.Empty, error)
	GetHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedAchievementServer()
}

// UnimplementedAchievementServer must be embedded to have forward compatible implementations.
type UnimplementedAchievementServer struct {
}

func (UnimplementedAchievementServer) SetAchievementAgree(context.Context, *SetAchievementAgreeReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAchievementAgree not implemented")
}
func (UnimplementedAchievementServer) SetAchievementView(context.Context, *SetAchievementViewReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAchievementView not implemented")
}
func (UnimplementedAchievementServer) SetAchievementCollect(context.Context, *SetAchievementCollectReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAchievementCollect not implemented")
}
func (UnimplementedAchievementServer) GetHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHealth not implemented")
}
func (UnimplementedAchievementServer) mustEmbedUnimplementedAchievementServer() {}

// UnsafeAchievementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AchievementServer will
// result in compilation errors.
type UnsafeAchievementServer interface {
	mustEmbedUnimplementedAchievementServer()
}

func RegisterAchievementServer(s grpc.ServiceRegistrar, srv AchievementServer) {
	s.RegisterService(&Achievement_ServiceDesc, srv)
}

func _Achievement_SetAchievementAgree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAchievementAgreeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AchievementServer).SetAchievementAgree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/achievement.v1.Achievement/SetAchievementAgree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AchievementServer).SetAchievementAgree(ctx, req.(*SetAchievementAgreeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Achievement_SetAchievementView_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAchievementViewReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AchievementServer).SetAchievementView(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/achievement.v1.Achievement/SetAchievementView",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AchievementServer).SetAchievementView(ctx, req.(*SetAchievementViewReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Achievement_SetAchievementCollect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAchievementCollectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AchievementServer).SetAchievementCollect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/achievement.v1.Achievement/SetAchievementCollect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AchievementServer).SetAchievementCollect(ctx, req.(*SetAchievementCollectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Achievement_GetHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AchievementServer).GetHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/achievement.v1.Achievement/GetHealth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AchievementServer).GetHealth(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Achievement_ServiceDesc is the grpc.ServiceDesc for Achievement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Achievement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "achievement.v1.Achievement",
	HandlerType: (*AchievementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetAchievementAgree",
			Handler:    _Achievement_SetAchievementAgree_Handler,
		},
		{
			MethodName: "SetAchievementView",
			Handler:    _Achievement_SetAchievementView_Handler,
		},
		{
			MethodName: "SetAchievementCollect",
			Handler:    _Achievement_SetAchievementCollect_Handler,
		},
		{
			MethodName: "GetHealth",
			Handler:    _Achievement_GetHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "achievement/service/v1/achievement.proto",
}