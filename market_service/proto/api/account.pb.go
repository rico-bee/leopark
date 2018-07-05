// Code generated by protoc-gen-go. DO NOT EDIT.
// source: account.proto

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	account.proto

It has these top-level messages:
	CreateAccountRequest
	CreateAccountResponse
*/
package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type CreateAccountRequest struct {
	Name  string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Email string `protobuf:"bytes,2,opt,name=email" json:"email,omitempty"`
}

func (m *CreateAccountRequest) Reset()                    { *m = CreateAccountRequest{} }
func (m *CreateAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateAccountRequest) ProtoMessage()               {}
func (*CreateAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CreateAccountRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateAccountRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

// The response message containing the greetings
type CreateAccountResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *CreateAccountResponse) Reset()                    { *m = CreateAccountResponse{} }
func (m *CreateAccountResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateAccountResponse) ProtoMessage()               {}
func (*CreateAccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateAccountResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateAccountRequest)(nil), "rpc.CreateAccountRequest")
	proto.RegisterType((*CreateAccountResponse)(nil), "rpc.CreateAccountResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Market service

type MarketClient interface {
	// Sends a greeting
	DoCreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
}

type marketClient struct {
	cc *grpc.ClientConn
}

func NewMarketClient(cc *grpc.ClientConn) MarketClient {
	return &marketClient{cc}
}

func (c *marketClient) DoCreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := grpc.Invoke(ctx, "/rpc.market/doCreateAccount", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Market service

type MarketServer interface {
	// Sends a greeting
	DoCreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
}

func RegisterMarketServer(s *grpc.Server, srv MarketServer) {
	s.RegisterService(&_Market_serviceDesc, srv)
}

func _Market_DoCreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketServer).DoCreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.market/DoCreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketServer).DoCreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Market_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.market",
	HandlerType: (*MarketServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "doCreateAccount",
			Handler:    _Market_DoCreateAccount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account.proto",
}

func init() { proto.RegisterFile("account.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 163 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0x4c, 0x4e, 0xce,
	0x2f, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0x72,
	0xe0, 0x12, 0x71, 0x2e, 0x4a, 0x4d, 0x2c, 0x49, 0x75, 0x84, 0xc8, 0x05, 0xa5, 0x16, 0x96, 0xa6,
	0x16, 0x97, 0x08, 0x09, 0x71, 0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70,
	0x06, 0x81, 0xd9, 0x42, 0x22, 0x5c, 0xac, 0xa9, 0xb9, 0x89, 0x99, 0x39, 0x12, 0x4c, 0x60, 0x41,
	0x08, 0x47, 0xc9, 0x90, 0x4b, 0x14, 0xcd, 0x84, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x09,
	0x2e, 0xf6, 0xdc, 0xd4, 0xe2, 0xe2, 0xc4, 0x74, 0x98, 0x29, 0x30, 0xae, 0x51, 0x08, 0x17, 0x5b,
	0x6e, 0x62, 0x51, 0x76, 0x6a, 0x89, 0x90, 0x17, 0x17, 0x7f, 0x4a, 0x3e, 0x8a, 0x76, 0x21, 0x49,
	0xbd, 0xa2, 0x82, 0x64, 0x3d, 0x6c, 0x8e, 0x92, 0x92, 0xc2, 0x26, 0x05, 0xb1, 0x4d, 0x89, 0x21,
	0x89, 0x0d, 0xec, 0x2d, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x92, 0x08, 0xa5, 0x39, 0xe7,
	0x00, 0x00, 0x00,
}