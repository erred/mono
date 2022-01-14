// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: auth/authnbpb/authnbpb.proto

package authnbpb

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

// AuthnBClient is the client API for AuthnB service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthnBClient interface {
	GetSession(ctx context.Context, in *GetSessionRequest, opts ...grpc.CallOption) (*GetSessionResponse, error)
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error)
	DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*DeleteSessionResponse, error)
	GetUserAuth(ctx context.Context, in *GetUserAuthRequest, opts ...grpc.CallOption) (*GetUserAuthResponse, error)
}

type authnBClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthnBClient(cc grpc.ClientConnInterface) AuthnBClient {
	return &authnBClient{cc}
}

func (c *authnBClient) GetSession(ctx context.Context, in *GetSessionRequest, opts ...grpc.CallOption) (*GetSessionResponse, error) {
	out := new(GetSessionResponse)
	err := c.cc.Invoke(ctx, "/seankhliao.auth.authnbpb.AuthnB/GetSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authnBClient) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error) {
	out := new(CreateSessionResponse)
	err := c.cc.Invoke(ctx, "/seankhliao.auth.authnbpb.AuthnB/CreateSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authnBClient) DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*DeleteSessionResponse, error) {
	out := new(DeleteSessionResponse)
	err := c.cc.Invoke(ctx, "/seankhliao.auth.authnbpb.AuthnB/DeleteSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authnBClient) GetUserAuth(ctx context.Context, in *GetUserAuthRequest, opts ...grpc.CallOption) (*GetUserAuthResponse, error) {
	out := new(GetUserAuthResponse)
	err := c.cc.Invoke(ctx, "/seankhliao.auth.authnbpb.AuthnB/GetUserAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthnBServer is the server API for AuthnB service.
// All implementations must embed UnimplementedAuthnBServer
// for forward compatibility
type AuthnBServer interface {
	GetSession(context.Context, *GetSessionRequest) (*GetSessionResponse, error)
	CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error)
	DeleteSession(context.Context, *DeleteSessionRequest) (*DeleteSessionResponse, error)
	GetUserAuth(context.Context, *GetUserAuthRequest) (*GetUserAuthResponse, error)
	mustEmbedUnimplementedAuthnBServer()
}

// UnimplementedAuthnBServer must be embedded to have forward compatible implementations.
type UnimplementedAuthnBServer struct {
}

func (UnimplementedAuthnBServer) GetSession(context.Context, *GetSessionRequest) (*GetSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSession not implemented")
}
func (UnimplementedAuthnBServer) CreateSession(context.Context, *CreateSessionRequest) (*CreateSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (UnimplementedAuthnBServer) DeleteSession(context.Context, *DeleteSessionRequest) (*DeleteSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSession not implemented")
}
func (UnimplementedAuthnBServer) GetUserAuth(context.Context, *GetUserAuthRequest) (*GetUserAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserAuth not implemented")
}
func (UnimplementedAuthnBServer) mustEmbedUnimplementedAuthnBServer() {}

// UnsafeAuthnBServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthnBServer will
// result in compilation errors.
type UnsafeAuthnBServer interface {
	mustEmbedUnimplementedAuthnBServer()
}

func RegisterAuthnBServer(s grpc.ServiceRegistrar, srv AuthnBServer) {
	s.RegisterService(&AuthnB_ServiceDesc, srv)
}

func _AuthnB_GetSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthnBServer).GetSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seankhliao.auth.authnbpb.AuthnB/GetSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthnBServer).GetSession(ctx, req.(*GetSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthnB_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthnBServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seankhliao.auth.authnbpb.AuthnB/CreateSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthnBServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthnB_DeleteSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthnBServer).DeleteSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seankhliao.auth.authnbpb.AuthnB/DeleteSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthnBServer).DeleteSession(ctx, req.(*DeleteSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuthnB_GetUserAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthnBServer).GetUserAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/seankhliao.auth.authnbpb.AuthnB/GetUserAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthnBServer).GetUserAuth(ctx, req.(*GetUserAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AuthnB_ServiceDesc is the grpc.ServiceDesc for AuthnB service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuthnB_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "seankhliao.auth.authnbpb.AuthnB",
	HandlerType: (*AuthnBServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSession",
			Handler:    _AuthnB_GetSession_Handler,
		},
		{
			MethodName: "CreateSession",
			Handler:    _AuthnB_CreateSession_Handler,
		},
		{
			MethodName: "DeleteSession",
			Handler:    _AuthnB_DeleteSession_Handler,
		},
		{
			MethodName: "GetUserAuth",
			Handler:    _AuthnB_GetUserAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth/authnbpb/authnbpb.proto",
}
