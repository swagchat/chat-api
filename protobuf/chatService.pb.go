// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chatService.proto

package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf1 "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ChatIncoming service

type ChatIncomingClient interface {
	PostMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type chatIncomingClient struct {
	cc *grpc.ClientConn
}

func NewChatIncomingClient(cc *grpc.ClientConn) ChatIncomingClient {
	return &chatIncomingClient{cc}
}

func (c *chatIncomingClient) PostMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.ChatIncoming/PostMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChatIncoming service

type ChatIncomingServer interface {
	PostMessage(context.Context, *Message) (*google_protobuf1.Empty, error)
}

func RegisterChatIncomingServer(s *grpc.Server, srv ChatIncomingServer) {
	s.RegisterService(&_ChatIncoming_serviceDesc, srv)
}

func _ChatIncoming_PostMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatIncomingServer).PostMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.ChatIncoming/PostMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatIncomingServer).PostMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChatIncoming_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swagchat.protobuf.ChatIncoming",
	HandlerType: (*ChatIncomingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostMessage",
			Handler:    _ChatIncoming_PostMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chatService.proto",
}

// Client API for ChatOutgoing service

type ChatOutgoingClient interface {
	PostWebhookRoom(ctx context.Context, in *Room, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	PostWebhookMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type chatOutgoingClient struct {
	cc *grpc.ClientConn
}

func NewChatOutgoingClient(cc *grpc.ClientConn) ChatOutgoingClient {
	return &chatOutgoingClient{cc}
}

func (c *chatOutgoingClient) PostWebhookRoom(ctx context.Context, in *Room, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.ChatOutgoing/PostWebhookRoom", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatOutgoingClient) PostWebhookMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.ChatOutgoing/PostWebhookMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChatOutgoing service

type ChatOutgoingServer interface {
	PostWebhookRoom(context.Context, *Room) (*google_protobuf1.Empty, error)
	PostWebhookMessage(context.Context, *Message) (*google_protobuf1.Empty, error)
}

func RegisterChatOutgoingServer(s *grpc.Server, srv ChatOutgoingServer) {
	s.RegisterService(&_ChatOutgoing_serviceDesc, srv)
}

func _ChatOutgoing_PostWebhookRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Room)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatOutgoingServer).PostWebhookRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.ChatOutgoing/PostWebhookRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatOutgoingServer).PostWebhookRoom(ctx, req.(*Room))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatOutgoing_PostWebhookMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatOutgoingServer).PostWebhookMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.ChatOutgoing/PostWebhookMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatOutgoingServer).PostWebhookMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _ChatOutgoing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swagchat.protobuf.ChatOutgoing",
	HandlerType: (*ChatOutgoingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostWebhookRoom",
			Handler:    _ChatOutgoing_PostWebhookRoom_Handler,
		},
		{
			MethodName: "PostWebhookMessage",
			Handler:    _ChatOutgoing_PostWebhookMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chatService.proto",
}

// Client API for UserRoleService service

type UserRoleServiceClient interface {
	PostUserRole(ctx context.Context, in *PostUserRoleReq, opts ...grpc.CallOption) (*UserRole, error)
	GetRoleIDsOfUserRole(ctx context.Context, in *GetRoleIDsOfUserRoleReq, opts ...grpc.CallOption) (*RoleIDs, error)
	GetUserIDsOfUserRole(ctx context.Context, in *GetUserIDsOfUserRoleReq, opts ...grpc.CallOption) (*UserIDs, error)
	DeleteUserRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type userRoleServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserRoleServiceClient(cc *grpc.ClientConn) UserRoleServiceClient {
	return &userRoleServiceClient{cc}
}

func (c *userRoleServiceClient) PostUserRole(ctx context.Context, in *PostUserRoleReq, opts ...grpc.CallOption) (*UserRole, error) {
	out := new(UserRole)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserRoleService/PostUserRole", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRoleServiceClient) GetRoleIDsOfUserRole(ctx context.Context, in *GetRoleIDsOfUserRoleReq, opts ...grpc.CallOption) (*RoleIDs, error) {
	out := new(RoleIDs)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserRoleService/GetRoleIDsOfUserRole", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRoleServiceClient) GetUserIDsOfUserRole(ctx context.Context, in *GetUserIDsOfUserRoleReq, opts ...grpc.CallOption) (*UserIDs, error) {
	out := new(UserIDs)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserRoleService/GetUserIDsOfUserRole", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userRoleServiceClient) DeleteUserRole(ctx context.Context, in *UserRole, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.UserRoleService/DeleteUserRole", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserRoleService service

type UserRoleServiceServer interface {
	PostUserRole(context.Context, *PostUserRoleReq) (*UserRole, error)
	GetRoleIDsOfUserRole(context.Context, *GetRoleIDsOfUserRoleReq) (*RoleIDs, error)
	GetUserIDsOfUserRole(context.Context, *GetUserIDsOfUserRoleReq) (*UserIDs, error)
	DeleteUserRole(context.Context, *UserRole) (*google_protobuf1.Empty, error)
}

func RegisterUserRoleServiceServer(s *grpc.Server, srv UserRoleServiceServer) {
	s.RegisterService(&_UserRoleService_serviceDesc, srv)
}

func _UserRoleService_PostUserRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostUserRoleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRoleServiceServer).PostUserRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserRoleService/PostUserRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRoleServiceServer).PostUserRole(ctx, req.(*PostUserRoleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRoleService_GetRoleIDsOfUserRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleIDsOfUserRoleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRoleServiceServer).GetRoleIDsOfUserRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserRoleService/GetRoleIDsOfUserRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRoleServiceServer).GetRoleIDsOfUserRole(ctx, req.(*GetRoleIDsOfUserRoleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRoleService_GetUserIDsOfUserRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserIDsOfUserRoleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRoleServiceServer).GetUserIDsOfUserRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserRoleService/GetUserIDsOfUserRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRoleServiceServer).GetUserIDsOfUserRole(ctx, req.(*GetUserIDsOfUserRoleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserRoleService_DeleteUserRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRole)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserRoleServiceServer).DeleteUserRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.UserRoleService/DeleteUserRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserRoleServiceServer).DeleteUserRole(ctx, req.(*UserRole))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserRoleService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swagchat.protobuf.UserRoleService",
	HandlerType: (*UserRoleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostUserRole",
			Handler:    _UserRoleService_PostUserRole_Handler,
		},
		{
			MethodName: "GetRoleIDsOfUserRole",
			Handler:    _UserRoleService_GetRoleIDsOfUserRole_Handler,
		},
		{
			MethodName: "GetUserIDsOfUserRole",
			Handler:    _UserRoleService_GetUserIDsOfUserRole_Handler,
		},
		{
			MethodName: "DeleteUserRole",
			Handler:    _UserRoleService_DeleteUserRole_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chatService.proto",
}

// Client API for RoomUserService service

type RoomUserServiceClient interface {
	PostRoomUser(ctx context.Context, in *PostRoomUserReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
	DeleteRoomUser(ctx context.Context, in *DeleteRoomUserReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error)
}

type roomUserServiceClient struct {
	cc *grpc.ClientConn
}

func NewRoomUserServiceClient(cc *grpc.ClientConn) RoomUserServiceClient {
	return &roomUserServiceClient{cc}
}

func (c *roomUserServiceClient) PostRoomUser(ctx context.Context, in *PostRoomUserReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.RoomUserService/PostRoomUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomUserServiceClient) DeleteRoomUser(ctx context.Context, in *DeleteRoomUserReq, opts ...grpc.CallOption) (*google_protobuf1.Empty, error) {
	out := new(google_protobuf1.Empty)
	err := grpc.Invoke(ctx, "/swagchat.protobuf.RoomUserService/DeleteRoomUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RoomUserService service

type RoomUserServiceServer interface {
	PostRoomUser(context.Context, *PostRoomUserReq) (*google_protobuf1.Empty, error)
	DeleteRoomUser(context.Context, *DeleteRoomUserReq) (*google_protobuf1.Empty, error)
}

func RegisterRoomUserServiceServer(s *grpc.Server, srv RoomUserServiceServer) {
	s.RegisterService(&_RoomUserService_serviceDesc, srv)
}

func _RoomUserService_PostRoomUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostRoomUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomUserServiceServer).PostRoomUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.RoomUserService/PostRoomUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomUserServiceServer).PostRoomUser(ctx, req.(*PostRoomUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomUserService_DeleteRoomUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoomUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomUserServiceServer).DeleteRoomUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swagchat.protobuf.RoomUserService/DeleteRoomUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomUserServiceServer).DeleteRoomUser(ctx, req.(*DeleteRoomUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _RoomUserService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swagchat.protobuf.RoomUserService",
	HandlerType: (*RoomUserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostRoomUser",
			Handler:    _RoomUserService_PostRoomUser_Handler,
		},
		{
			MethodName: "DeleteRoomUser",
			Handler:    _RoomUserService_DeleteRoomUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chatService.proto",
}

func init() { proto.RegisterFile("chatService.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0x4f, 0xeb, 0xd3, 0x30,
	0x18, 0xc7, 0xad, 0x82, 0x87, 0x38, 0x56, 0x17, 0x75, 0x6a, 0x37, 0x3c, 0x14, 0x4f, 0x3d, 0x34,
	0x50, 0x6f, 0x3b, 0xba, 0x89, 0xec, 0x20, 0x93, 0x0d, 0x11, 0x05, 0x0f, 0x69, 0x49, 0xff, 0xe0,
	0xda, 0x67, 0x36, 0xa9, 0xd3, 0xab, 0x6f, 0xc1, 0x37, 0xe5, 0x5d, 0x04, 0x5f, 0x80, 0x2f, 0x44,
	0x92, 0x26, 0xed, 0x7e, 0x2c, 0x1d, 0x3f, 0xf8, 0x5d, 0x46, 0xf2, 0x3c, 0xdf, 0x7c, 0x3f, 0x79,
	0xf2, 0x1d, 0x45, 0x93, 0x24, 0xa7, 0x62, 0xc7, 0xea, 0xaf, 0x45, 0xc2, 0xc2, 0x43, 0x0d, 0x02,
	0xf0, 0x84, 0x1f, 0x69, 0x26, 0xcb, 0xed, 0x3e, 0x6e, 0x52, 0x6f, 0x9e, 0x01, 0x64, 0x7b, 0x46,
	0xe8, 0xa1, 0x20, 0xb4, 0xaa, 0x40, 0x50, 0x51, 0x40, 0xc5, 0x5b, 0x81, 0x37, 0xd3, 0x5d, 0x23,
	0x27, 0xac, 0x3c, 0x88, 0xef, 0xba, 0xa9, 0x00, 0x6f, 0x18, 0xe7, 0x34, 0xd3, 0x80, 0x68, 0x87,
	0x46, 0xcb, 0x9c, 0x8a, 0x75, 0x95, 0x40, 0x59, 0x54, 0x19, 0x5e, 0xa2, 0x7b, 0x6f, 0x81, 0x1b,
	0x11, 0xf6, 0xc2, 0xb3, 0x0b, 0x84, 0xba, 0xe7, 0x4d, 0xc3, 0x96, 0xd5, 0x77, 0x5e, 0x49, 0x96,
	0x7f, 0x2b, 0xfa, 0xe3, 0xb4, 0xae, 0x9b, 0x46, 0x64, 0x20, 0x5d, 0x3f, 0x21, 0x57, 0xba, 0xbe,
	0x67, 0x71, 0x0e, 0xf0, 0x79, 0x0b, 0x50, 0xe2, 0xc7, 0x16, 0x67, 0xd9, 0x18, 0xb4, 0x7d, 0xfa,
	0xe3, 0xf7, 0xbf, 0x9f, 0xb7, 0x1f, 0xf8, 0x63, 0x72, 0x6c, 0x6d, 0x38, 0xa9, 0x01, 0xca, 0x85,
	0x13, 0xe0, 0x14, 0xe1, 0x13, 0xfb, 0x9b, 0xdc, 0x7d, 0xae, 0x20, 0x53, 0x7f, 0xd2, 0x43, 0xca,
	0xf6, 0xc8, 0xc2, 0x09, 0xa2, 0x5f, 0x77, 0x90, 0xfb, 0x8e, 0xb3, 0x7a, 0x0b, 0x7b, 0xa6, 0x73,
	0xc2, 0x29, 0x1a, 0x49, 0xb6, 0x29, 0x63, 0xdf, 0x42, 0x3d, 0x15, 0x6c, 0xd9, 0x17, 0x6f, 0x66,
	0xd1, 0x98, 0xbe, 0xff, 0x48, 0x5d, 0xc1, 0xf5, 0x11, 0x69, 0x74, 0x89, 0xcb, 0x19, 0xbf, 0xa1,
	0x87, 0xaf, 0x99, 0x90, 0xdb, 0xf5, 0x8a, 0x6f, 0xd2, 0x8e, 0x17, 0x58, 0xbc, 0x6c, 0x42, 0xc9,
	0xf5, 0xac, 0x6f, 0xae, 0x84, 0xfe, 0x13, 0x85, 0xc5, 0xf8, 0x7e, 0x8f, 0x25, 0xb5, 0xfc, 0xd5,
	0x64, 0xe9, 0x73, 0x2d, 0xf2, 0x99, 0x70, 0x88, 0xac, 0x85, 0x56, 0xb2, 0x5c, 0x71, 0xfc, 0x01,
	0x8d, 0x57, 0x6c, 0xcf, 0x04, 0xeb, 0x98, 0x97, 0x5e, 0x6e, 0x30, 0x54, 0xac, 0x00, 0xa3, 0xe0,
	0xe4, 0x45, 0xa3, 0xbf, 0x0e, 0x72, 0xe5, 0xdf, 0x4d, 0x1e, 0x36, 0x51, 0xd2, 0x36, 0x4a, 0x53,
	0x1e, 0x8c, 0xd2, 0x08, 0xe4, 0x60, 0x43, 0xcc, 0x3e, 0xc5, 0x5a, 0xab, 0x55, 0x8a, 0xb1, 0x99,
	0xa8, 0x83, 0x3c, 0xb7, 0x40, 0xae, 0x4a, 0x2e, 0x61, 0xfa, 0xd1, 0x3a, 0xcc, 0xcb, 0x67, 0x1f,
	0xe7, 0x59, 0x21, 0xf2, 0x26, 0x0e, 0x13, 0x28, 0x89, 0x71, 0xef, 0xbe, 0x08, 0xf1, 0x5d, 0xb5,
	0x7a, 0xf1, 0x3f, 0x00, 0x00, 0xff, 0xff, 0x54, 0x8a, 0xe3, 0xa0, 0x6f, 0x04, 0x00, 0x00,
}