// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: bff/interface/v1/bff.proto

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

// BffClient is the client API for Bff service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BffClient interface {
	UserRegister(ctx context.Context, in *UserRegisterReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	LoginByPassword(ctx context.Context, in *LoginByPasswordReq, opts ...grpc.CallOption) (*LoginReply, error)
	LoginByCode(ctx context.Context, in *LoginByCodeReq, opts ...grpc.CallOption) (*LoginReply, error)
	LoginByWeChat(ctx context.Context, in *LoginByWeChatReq, opts ...grpc.CallOption) (*LoginReply, error)
	LoginByGithub(ctx context.Context, in *LoginByGithubReq, opts ...grpc.CallOption) (*LoginReply, error)
	LoginPasswordReset(ctx context.Context, in *LoginPasswordResetReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendPhoneCode(ctx context.Context, in *SendPhoneCodeReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SendEmailCode(ctx context.Context, in *SendEmailCodeReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetCosSessionKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetCosSessionKeyReply, error)
	GetAccount(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAccountReply, error)
	GetProfile(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetProfileReply, error)
	GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoReply, error)
	GetProfileUpdate(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetProfileUpdateReply, error)
	SetProfileUpdate(ctx context.Context, in *SetProfileUpdateReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetUserPhone(ctx context.Context, in *SetUserPhoneReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetUserEmail(ctx context.Context, in *SetUserEmailReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetUserPassword(ctx context.Context, in *SetUserPasswordReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ChangeUserPassword(ctx context.Context, in *ChangeUserPasswordReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnbindUserPhone(ctx context.Context, in *UnbindUserPhoneReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	UnbindUserEmail(ctx context.Context, in *UnbindUserEmailReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetArticleList(ctx context.Context, in *GetArticleListReq, opts ...grpc.CallOption) (*GetArticleListReply, error)
	GetLastArticleDraft(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetLastArticleDraftReply, error)
	CreateArticleDraft(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CreateArticleDraftReply, error)
	ArticleDraftMark(ctx context.Context, in *ArticleDraftMarkReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetArticleDraftList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetArticleDraftListReply, error)
	SendArticle(ctx context.Context, in *SendArticleReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type bffClient struct {
	cc grpc.ClientConnInterface
}

func NewBffClient(cc grpc.ClientConnInterface) BffClient {
	return &bffClient{cc}
}

func (c *bffClient) UserRegister(ctx context.Context, in *UserRegisterReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/UserRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) LoginByPassword(ctx context.Context, in *LoginByPasswordReq, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/LoginByPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) LoginByCode(ctx context.Context, in *LoginByCodeReq, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/LoginByCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) LoginByWeChat(ctx context.Context, in *LoginByWeChatReq, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/LoginByWeChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) LoginByGithub(ctx context.Context, in *LoginByGithubReq, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/LoginByGithub", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) LoginPasswordReset(ctx context.Context, in *LoginPasswordResetReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/LoginPasswordReset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SendPhoneCode(ctx context.Context, in *SendPhoneCodeReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SendPhoneCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SendEmailCode(ctx context.Context, in *SendEmailCodeReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SendEmailCode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetCosSessionKey(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetCosSessionKeyReply, error) {
	out := new(GetCosSessionKeyReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetCosSessionKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetAccount(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAccountReply, error) {
	out := new(GetAccountReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetProfile(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetProfileReply, error) {
	out := new(GetProfileReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoReply, error) {
	out := new(GetUserInfoReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetUserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetProfileUpdate(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetProfileUpdateReply, error) {
	out := new(GetProfileUpdateReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetProfileUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SetProfileUpdate(ctx context.Context, in *SetProfileUpdateReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SetProfileUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SetUserPhone(ctx context.Context, in *SetUserPhoneReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SetUserPhone", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SetUserEmail(ctx context.Context, in *SetUserEmailReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SetUserEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SetUserPassword(ctx context.Context, in *SetUserPasswordReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SetUserPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) ChangeUserPassword(ctx context.Context, in *ChangeUserPasswordReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/ChangeUserPassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) UnbindUserPhone(ctx context.Context, in *UnbindUserPhoneReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/UnbindUserPhone", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) UnbindUserEmail(ctx context.Context, in *UnbindUserEmailReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/UnbindUserEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetArticleList(ctx context.Context, in *GetArticleListReq, opts ...grpc.CallOption) (*GetArticleListReply, error) {
	out := new(GetArticleListReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetArticleList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetLastArticleDraft(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetLastArticleDraftReply, error) {
	out := new(GetLastArticleDraftReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetLastArticleDraft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) CreateArticleDraft(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CreateArticleDraftReply, error) {
	out := new(CreateArticleDraftReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/CreateArticleDraft", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) ArticleDraftMark(ctx context.Context, in *ArticleDraftMarkReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/ArticleDraftMark", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) GetArticleDraftList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetArticleDraftListReply, error) {
	out := new(GetArticleDraftListReply)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/GetArticleDraftList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bffClient) SendArticle(ctx context.Context, in *SendArticleReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bff.v1.Bff/SendArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BffServer is the server API for Bff service.
// All implementations must embed UnimplementedBffServer
// for forward compatibility
type BffServer interface {
	UserRegister(context.Context, *UserRegisterReq) (*emptypb.Empty, error)
	LoginByPassword(context.Context, *LoginByPasswordReq) (*LoginReply, error)
	LoginByCode(context.Context, *LoginByCodeReq) (*LoginReply, error)
	LoginByWeChat(context.Context, *LoginByWeChatReq) (*LoginReply, error)
	LoginByGithub(context.Context, *LoginByGithubReq) (*LoginReply, error)
	LoginPasswordReset(context.Context, *LoginPasswordResetReq) (*emptypb.Empty, error)
	SendPhoneCode(context.Context, *SendPhoneCodeReq) (*emptypb.Empty, error)
	SendEmailCode(context.Context, *SendEmailCodeReq) (*emptypb.Empty, error)
	GetCosSessionKey(context.Context, *emptypb.Empty) (*GetCosSessionKeyReply, error)
	GetAccount(context.Context, *emptypb.Empty) (*GetAccountReply, error)
	GetProfile(context.Context, *emptypb.Empty) (*GetProfileReply, error)
	GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoReply, error)
	GetProfileUpdate(context.Context, *emptypb.Empty) (*GetProfileUpdateReply, error)
	SetProfileUpdate(context.Context, *SetProfileUpdateReq) (*emptypb.Empty, error)
	SetUserPhone(context.Context, *SetUserPhoneReq) (*emptypb.Empty, error)
	SetUserEmail(context.Context, *SetUserEmailReq) (*emptypb.Empty, error)
	SetUserPassword(context.Context, *SetUserPasswordReq) (*emptypb.Empty, error)
	ChangeUserPassword(context.Context, *ChangeUserPasswordReq) (*emptypb.Empty, error)
	UnbindUserPhone(context.Context, *UnbindUserPhoneReq) (*emptypb.Empty, error)
	UnbindUserEmail(context.Context, *UnbindUserEmailReq) (*emptypb.Empty, error)
	GetArticleList(context.Context, *GetArticleListReq) (*GetArticleListReply, error)
	GetLastArticleDraft(context.Context, *emptypb.Empty) (*GetLastArticleDraftReply, error)
	CreateArticleDraft(context.Context, *emptypb.Empty) (*CreateArticleDraftReply, error)
	ArticleDraftMark(context.Context, *ArticleDraftMarkReq) (*emptypb.Empty, error)
	GetArticleDraftList(context.Context, *emptypb.Empty) (*GetArticleDraftListReply, error)
	SendArticle(context.Context, *SendArticleReq) (*emptypb.Empty, error)
	mustEmbedUnimplementedBffServer()
}

// UnimplementedBffServer must be embedded to have forward compatible implementations.
type UnimplementedBffServer struct {
}

func (UnimplementedBffServer) UserRegister(context.Context, *UserRegisterReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserRegister not implemented")
}
func (UnimplementedBffServer) LoginByPassword(context.Context, *LoginByPasswordReq) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByPassword not implemented")
}
func (UnimplementedBffServer) LoginByCode(context.Context, *LoginByCodeReq) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByCode not implemented")
}
func (UnimplementedBffServer) LoginByWeChat(context.Context, *LoginByWeChatReq) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByWeChat not implemented")
}
func (UnimplementedBffServer) LoginByGithub(context.Context, *LoginByGithubReq) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginByGithub not implemented")
}
func (UnimplementedBffServer) LoginPasswordReset(context.Context, *LoginPasswordResetReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginPasswordReset not implemented")
}
func (UnimplementedBffServer) SendPhoneCode(context.Context, *SendPhoneCodeReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPhoneCode not implemented")
}
func (UnimplementedBffServer) SendEmailCode(context.Context, *SendEmailCodeReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmailCode not implemented")
}
func (UnimplementedBffServer) GetCosSessionKey(context.Context, *emptypb.Empty) (*GetCosSessionKeyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCosSessionKey not implemented")
}
func (UnimplementedBffServer) GetAccount(context.Context, *emptypb.Empty) (*GetAccountReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedBffServer) GetProfile(context.Context, *emptypb.Empty) (*GetProfileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}
func (UnimplementedBffServer) GetUserInfo(context.Context, *GetUserInfoReq) (*GetUserInfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedBffServer) GetProfileUpdate(context.Context, *emptypb.Empty) (*GetProfileUpdateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfileUpdate not implemented")
}
func (UnimplementedBffServer) SetProfileUpdate(context.Context, *SetProfileUpdateReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetProfileUpdate not implemented")
}
func (UnimplementedBffServer) SetUserPhone(context.Context, *SetUserPhoneReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserPhone not implemented")
}
func (UnimplementedBffServer) SetUserEmail(context.Context, *SetUserEmailReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserEmail not implemented")
}
func (UnimplementedBffServer) SetUserPassword(context.Context, *SetUserPasswordReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUserPassword not implemented")
}
func (UnimplementedBffServer) ChangeUserPassword(context.Context, *ChangeUserPasswordReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeUserPassword not implemented")
}
func (UnimplementedBffServer) UnbindUserPhone(context.Context, *UnbindUserPhoneReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbindUserPhone not implemented")
}
func (UnimplementedBffServer) UnbindUserEmail(context.Context, *UnbindUserEmailReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnbindUserEmail not implemented")
}
func (UnimplementedBffServer) GetArticleList(context.Context, *GetArticleListReq) (*GetArticleListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleList not implemented")
}
func (UnimplementedBffServer) GetLastArticleDraft(context.Context, *emptypb.Empty) (*GetLastArticleDraftReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastArticleDraft not implemented")
}
func (UnimplementedBffServer) CreateArticleDraft(context.Context, *emptypb.Empty) (*CreateArticleDraftReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArticleDraft not implemented")
}
func (UnimplementedBffServer) ArticleDraftMark(context.Context, *ArticleDraftMarkReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArticleDraftMark not implemented")
}
func (UnimplementedBffServer) GetArticleDraftList(context.Context, *emptypb.Empty) (*GetArticleDraftListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticleDraftList not implemented")
}
func (UnimplementedBffServer) SendArticle(context.Context, *SendArticleReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendArticle not implemented")
}
func (UnimplementedBffServer) mustEmbedUnimplementedBffServer() {}

// UnsafeBffServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BffServer will
// result in compilation errors.
type UnsafeBffServer interface {
	mustEmbedUnimplementedBffServer()
}

func RegisterBffServer(s grpc.ServiceRegistrar, srv BffServer) {
	s.RegisterService(&Bff_ServiceDesc, srv)
}

func _Bff_UserRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).UserRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/UserRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).UserRegister(ctx, req.(*UserRegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_LoginByPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByPasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).LoginByPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/LoginByPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).LoginByPassword(ctx, req.(*LoginByPasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_LoginByCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByCodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).LoginByCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/LoginByCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).LoginByCode(ctx, req.(*LoginByCodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_LoginByWeChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByWeChatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).LoginByWeChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/LoginByWeChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).LoginByWeChat(ctx, req.(*LoginByWeChatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_LoginByGithub_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginByGithubReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).LoginByGithub(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/LoginByGithub",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).LoginByGithub(ctx, req.(*LoginByGithubReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_LoginPasswordReset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginPasswordResetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).LoginPasswordReset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/LoginPasswordReset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).LoginPasswordReset(ctx, req.(*LoginPasswordResetReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SendPhoneCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPhoneCodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SendPhoneCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SendPhoneCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SendPhoneCode(ctx, req.(*SendPhoneCodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SendEmailCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEmailCodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SendEmailCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SendEmailCode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SendEmailCode(ctx, req.(*SendEmailCodeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetCosSessionKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetCosSessionKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetCosSessionKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetCosSessionKey(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetAccount(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetProfile(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetUserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetUserInfo(ctx, req.(*GetUserInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetProfileUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetProfileUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetProfileUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetProfileUpdate(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SetProfileUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetProfileUpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SetProfileUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SetProfileUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SetProfileUpdate(ctx, req.(*SetProfileUpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SetUserPhone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUserPhoneReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SetUserPhone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SetUserPhone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SetUserPhone(ctx, req.(*SetUserPhoneReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SetUserEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUserEmailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SetUserEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SetUserEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SetUserEmail(ctx, req.(*SetUserEmailReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SetUserPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUserPasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SetUserPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SetUserPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SetUserPassword(ctx, req.(*SetUserPasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_ChangeUserPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeUserPasswordReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).ChangeUserPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/ChangeUserPassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).ChangeUserPassword(ctx, req.(*ChangeUserPasswordReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_UnbindUserPhone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnbindUserPhoneReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).UnbindUserPhone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/UnbindUserPhone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).UnbindUserPhone(ctx, req.(*UnbindUserPhoneReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_UnbindUserEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnbindUserEmailReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).UnbindUserEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/UnbindUserEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).UnbindUserEmail(ctx, req.(*UnbindUserEmailReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetArticleList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArticleListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetArticleList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetArticleList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetArticleList(ctx, req.(*GetArticleListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetLastArticleDraft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetLastArticleDraft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetLastArticleDraft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetLastArticleDraft(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_CreateArticleDraft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).CreateArticleDraft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/CreateArticleDraft",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).CreateArticleDraft(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_ArticleDraftMark_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArticleDraftMarkReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).ArticleDraftMark(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/ArticleDraftMark",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).ArticleDraftMark(ctx, req.(*ArticleDraftMarkReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_GetArticleDraftList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).GetArticleDraftList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/GetArticleDraftList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).GetArticleDraftList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bff_SendArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendArticleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BffServer).SendArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bff.v1.Bff/SendArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BffServer).SendArticle(ctx, req.(*SendArticleReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Bff_ServiceDesc is the grpc.ServiceDesc for Bff service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bff_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bff.v1.Bff",
	HandlerType: (*BffServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserRegister",
			Handler:    _Bff_UserRegister_Handler,
		},
		{
			MethodName: "LoginByPassword",
			Handler:    _Bff_LoginByPassword_Handler,
		},
		{
			MethodName: "LoginByCode",
			Handler:    _Bff_LoginByCode_Handler,
		},
		{
			MethodName: "LoginByWeChat",
			Handler:    _Bff_LoginByWeChat_Handler,
		},
		{
			MethodName: "LoginByGithub",
			Handler:    _Bff_LoginByGithub_Handler,
		},
		{
			MethodName: "LoginPasswordReset",
			Handler:    _Bff_LoginPasswordReset_Handler,
		},
		{
			MethodName: "SendPhoneCode",
			Handler:    _Bff_SendPhoneCode_Handler,
		},
		{
			MethodName: "SendEmailCode",
			Handler:    _Bff_SendEmailCode_Handler,
		},
		{
			MethodName: "GetCosSessionKey",
			Handler:    _Bff_GetCosSessionKey_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _Bff_GetAccount_Handler,
		},
		{
			MethodName: "GetProfile",
			Handler:    _Bff_GetProfile_Handler,
		},
		{
			MethodName: "GetUserInfo",
			Handler:    _Bff_GetUserInfo_Handler,
		},
		{
			MethodName: "GetProfileUpdate",
			Handler:    _Bff_GetProfileUpdate_Handler,
		},
		{
			MethodName: "SetProfileUpdate",
			Handler:    _Bff_SetProfileUpdate_Handler,
		},
		{
			MethodName: "SetUserPhone",
			Handler:    _Bff_SetUserPhone_Handler,
		},
		{
			MethodName: "SetUserEmail",
			Handler:    _Bff_SetUserEmail_Handler,
		},
		{
			MethodName: "SetUserPassword",
			Handler:    _Bff_SetUserPassword_Handler,
		},
		{
			MethodName: "ChangeUserPassword",
			Handler:    _Bff_ChangeUserPassword_Handler,
		},
		{
			MethodName: "UnbindUserPhone",
			Handler:    _Bff_UnbindUserPhone_Handler,
		},
		{
			MethodName: "UnbindUserEmail",
			Handler:    _Bff_UnbindUserEmail_Handler,
		},
		{
			MethodName: "GetArticleList",
			Handler:    _Bff_GetArticleList_Handler,
		},
		{
			MethodName: "GetLastArticleDraft",
			Handler:    _Bff_GetLastArticleDraft_Handler,
		},
		{
			MethodName: "CreateArticleDraft",
			Handler:    _Bff_CreateArticleDraft_Handler,
		},
		{
			MethodName: "ArticleDraftMark",
			Handler:    _Bff_ArticleDraftMark_Handler,
		},
		{
			MethodName: "GetArticleDraftList",
			Handler:    _Bff_GetArticleDraftList_Handler,
		},
		{
			MethodName: "SendArticle",
			Handler:    _Bff_SendArticle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bff/interface/v1/bff.proto",
}
