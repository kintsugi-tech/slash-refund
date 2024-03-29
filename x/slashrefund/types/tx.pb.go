// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: slashrefund/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MsgDeposit struct {
	// option (gogoproto.equal)           = false;
	// option (gogoproto.goproto_getters) = false;
	DepositorAddress string     `protobuf:"bytes,1,opt,name=depositor_address,json=depositorAddress,proto3" json:"depositor_address,omitempty"`
	ValidatorAddress string     `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Amount           types.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgDeposit) Reset()         { *m = MsgDeposit{} }
func (m *MsgDeposit) String() string { return proto.CompactTextString(m) }
func (*MsgDeposit) ProtoMessage()    {}
func (*MsgDeposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_15788a07da2b95d1, []int{0}
}
func (m *MsgDeposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeposit.Merge(m, src)
}
func (m *MsgDeposit) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeposit) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeposit.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeposit proto.InternalMessageInfo

func (m *MsgDeposit) GetDepositorAddress() string {
	if m != nil {
		return m.DepositorAddress
	}
	return ""
}

func (m *MsgDeposit) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *MsgDeposit) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

type MsgDepositResponse struct {
}

func (m *MsgDepositResponse) Reset()         { *m = MsgDepositResponse{} }
func (m *MsgDepositResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDepositResponse) ProtoMessage()    {}
func (*MsgDepositResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_15788a07da2b95d1, []int{1}
}
func (m *MsgDepositResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDepositResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDepositResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDepositResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDepositResponse.Merge(m, src)
}
func (m *MsgDepositResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDepositResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDepositResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDepositResponse proto.InternalMessageInfo

type MsgWithdraw struct {
	DepositorAddress string     `protobuf:"bytes,1,opt,name=depositor_address,json=depositorAddress,proto3" json:"depositor_address,omitempty"`
	ValidatorAddress string     `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Amount           types.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
}

func (m *MsgWithdraw) Reset()         { *m = MsgWithdraw{} }
func (m *MsgWithdraw) String() string { return proto.CompactTextString(m) }
func (*MsgWithdraw) ProtoMessage()    {}
func (*MsgWithdraw) Descriptor() ([]byte, []int) {
	return fileDescriptor_15788a07da2b95d1, []int{2}
}
func (m *MsgWithdraw) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdraw) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdraw.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdraw) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdraw.Merge(m, src)
}
func (m *MsgWithdraw) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdraw) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdraw.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdraw proto.InternalMessageInfo

func (m *MsgWithdraw) GetDepositorAddress() string {
	if m != nil {
		return m.DepositorAddress
	}
	return ""
}

func (m *MsgWithdraw) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *MsgWithdraw) GetAmount() types.Coin {
	if m != nil {
		return m.Amount
	}
	return types.Coin{}
}

type MsgWithdrawResponse struct {
	CompletionTime time.Time `protobuf:"bytes,1,opt,name=completion_time,json=completionTime,proto3,stdtime" json:"completion_time"`
}

func (m *MsgWithdrawResponse) Reset()         { *m = MsgWithdrawResponse{} }
func (m *MsgWithdrawResponse) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawResponse) ProtoMessage()    {}
func (*MsgWithdrawResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_15788a07da2b95d1, []int{3}
}
func (m *MsgWithdrawResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawResponse.Merge(m, src)
}
func (m *MsgWithdrawResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawResponse proto.InternalMessageInfo

func (m *MsgWithdrawResponse) GetCompletionTime() time.Time {
	if m != nil {
		return m.CompletionTime
	}
	return time.Time{}
}

type MsgClaim struct {
	// option (gogoproto.equal)           = false;
	// option (gogoproto.goproto_getters) = false;
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
}

func (m *MsgClaim) Reset()         { *m = MsgClaim{} }
func (m *MsgClaim) String() string { return proto.CompactTextString(m) }
func (*MsgClaim) ProtoMessage()    {}
func (*MsgClaim) Descriptor() ([]byte, []int) {
	return fileDescriptor_15788a07da2b95d1, []int{4}
}
func (m *MsgClaim) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaim.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaim.Merge(m, src)
}
func (m *MsgClaim) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaim) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaim.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaim proto.InternalMessageInfo

func (m *MsgClaim) GetDelegatorAddress() string {
	if m != nil {
		return m.DelegatorAddress
	}
	return ""
}

func (m *MsgClaim) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

type MsgClaimResponse struct {
}

func (m *MsgClaimResponse) Reset()         { *m = MsgClaimResponse{} }
func (m *MsgClaimResponse) String() string { return proto.CompactTextString(m) }
func (*MsgClaimResponse) ProtoMessage()    {}
func (*MsgClaimResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_15788a07da2b95d1, []int{5}
}
func (m *MsgClaimResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimResponse.Merge(m, src)
}
func (m *MsgClaimResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgDeposit)(nil), "madeinblock.slashrefund.slashrefund.MsgDeposit")
	proto.RegisterType((*MsgDepositResponse)(nil), "madeinblock.slashrefund.slashrefund.MsgDepositResponse")
	proto.RegisterType((*MsgWithdraw)(nil), "madeinblock.slashrefund.slashrefund.MsgWithdraw")
	proto.RegisterType((*MsgWithdrawResponse)(nil), "madeinblock.slashrefund.slashrefund.MsgWithdrawResponse")
	proto.RegisterType((*MsgClaim)(nil), "madeinblock.slashrefund.slashrefund.MsgClaim")
	proto.RegisterType((*MsgClaimResponse)(nil), "madeinblock.slashrefund.slashrefund.MsgClaimResponse")
}

func init() { proto.RegisterFile("slashrefund/tx.proto", fileDescriptor_15788a07da2b95d1) }

var fileDescriptor_15788a07da2b95d1 = []byte{
	// 524 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x54, 0x4f, 0x6f, 0xd3, 0x30,
	0x1c, 0xad, 0x37, 0x18, 0xc3, 0x95, 0x60, 0x0b, 0x15, 0x74, 0x39, 0xa4, 0x53, 0xb9, 0x4c, 0x48,
	0xb5, 0x69, 0x11, 0x1a, 0x70, 0xa3, 0x83, 0x63, 0x2e, 0x1d, 0x12, 0x12, 0x97, 0xca, 0x49, 0x3c,
	0xd7, 0x22, 0x8e, 0xa3, 0xd8, 0x2d, 0xe3, 0xca, 0x27, 0xd8, 0x47, 0xd9, 0x81, 0x0f, 0xb1, 0xe3,
	0xc4, 0x89, 0x13, 0x7f, 0xda, 0x03, 0x07, 0xf8, 0x10, 0x28, 0xb1, 0x93, 0x06, 0x24, 0x44, 0x2b,
	0x71, 0xe2, 0x14, 0x3b, 0xef, 0xf7, 0x9e, 0xfd, 0x7b, 0x7e, 0x36, 0x6c, 0xa9, 0x98, 0xa8, 0x49,
	0x46, 0x4f, 0xa6, 0x49, 0x84, 0xf5, 0x29, 0x4a, 0x33, 0xa9, 0xa5, 0x73, 0x57, 0x90, 0x88, 0xf2,
	0x24, 0x88, 0x65, 0xf8, 0x1a, 0xd5, 0x2a, 0xea, 0x63, 0xb7, 0xc5, 0x24, 0x93, 0x45, 0x3d, 0xce,
	0x47, 0x86, 0xea, 0x7a, 0xa1, 0x54, 0x42, 0x2a, 0x1c, 0x10, 0x45, 0xf1, 0xac, 0x1f, 0x50, 0x4d,
	0xfa, 0x38, 0x94, 0x3c, 0xb1, 0xf8, 0x1d, 0x8b, 0x0b, 0xc5, 0xf0, 0xac, 0x9f, 0x7f, 0x2c, 0xb0,
	0x67, 0x80, 0xb1, 0x51, 0x34, 0x13, 0x0b, 0x75, 0x98, 0x94, 0x2c, 0xa6, 0xb8, 0x98, 0x05, 0xd3,
	0x13, 0xac, 0xb9, 0xa0, 0x4a, 0x13, 0x91, 0x9a, 0x82, 0xee, 0x77, 0x00, 0xa1, 0xaf, 0xd8, 0x33,
	0x9a, 0x4a, 0xc5, 0xb5, 0xf3, 0x1c, 0xee, 0x46, 0x66, 0x28, 0xb3, 0x31, 0x89, 0xa2, 0x8c, 0x2a,
	0xd5, 0x06, 0xfb, 0xe0, 0xe0, 0xfa, 0xb0, 0xfd, 0xe1, 0x7d, 0xaf, 0x65, 0xc5, 0x9f, 0x1a, 0xe4,
	0x58, 0x67, 0x3c, 0x61, 0xa3, 0x9d, 0x8a, 0x62, 0xff, 0xe7, 0x32, 0x33, 0x12, 0xf3, 0x88, 0xd4,
	0x65, 0x36, 0xfe, 0x26, 0x53, 0x51, 0x4a, 0x99, 0x43, 0xb8, 0x45, 0x84, 0x9c, 0x26, 0xba, 0xbd,
	0xb9, 0x0f, 0x0e, 0x9a, 0x83, 0x3d, 0x64, 0x89, 0xb9, 0x45, 0xc8, 0x5a, 0x84, 0x8e, 0x24, 0x4f,
	0x86, 0x57, 0x2e, 0x3e, 0x75, 0x1a, 0x23, 0x5b, 0xfe, 0xe4, 0xf6, 0xbb, 0x6f, 0xe7, 0xf7, 0x76,
	0x23, 0x1a, 0x53, 0x56, 0xdf, 0x42, 0xb7, 0x05, 0x9d, 0x65, 0xb3, 0x23, 0xaa, 0x52, 0x99, 0x28,
	0xda, 0xfd, 0x01, 0x60, 0xd3, 0x57, 0xec, 0x25, 0xd7, 0x93, 0x28, 0x23, 0x6f, 0xfe, 0x77, 0x13,
	0x22, 0x78, 0xab, 0xd6, 0x6d, 0xe9, 0x82, 0xe3, 0xc3, 0x9b, 0xa1, 0x14, 0x69, 0x4c, 0x35, 0x97,
	0xc9, 0x38, 0xcf, 0x49, 0xd1, 0x73, 0x73, 0xe0, 0x22, 0x13, 0x22, 0x54, 0x86, 0x08, 0xbd, 0x28,
	0x43, 0x34, 0xdc, 0xce, 0x57, 0x3c, 0xfb, 0xdc, 0x01, 0xa3, 0x1b, 0x4b, 0x72, 0x0e, 0x77, 0xcf,
	0x01, 0xdc, 0xf6, 0x15, 0x3b, 0x8a, 0x09, 0x17, 0xc6, 0xd1, 0xdf, 0xf6, 0xb1, 0x8a, 0xa3, 0x96,
	0xf2, 0x6f, 0x1d, 0xfd, 0xa3, 0x31, 0x0e, 0xdc, 0x29, 0x77, 0x5c, 0xba, 0x32, 0xf8, 0xba, 0x01,
	0x37, 0x7d, 0xc5, 0x1c, 0x05, 0xaf, 0x95, 0x77, 0x04, 0xa3, 0x15, 0xee, 0x38, 0x5a, 0xe6, 0xcc,
	0x3d, 0x5c, 0x93, 0x50, 0x1d, 0xc9, 0x0c, 0x6e, 0x57, 0xa1, 0xbc, 0xbf, 0xaa, 0x48, 0xc9, 0x70,
	0x1f, 0xad, 0xcb, 0xa8, 0xd6, 0x15, 0xf0, 0xaa, 0x39, 0xb7, 0xde, 0xaa, 0x12, 0x45, 0xb9, 0xfb,
	0x70, 0xad, 0xf2, 0x72, 0xb9, 0xe1, 0xf1, 0xc5, 0xdc, 0x03, 0x97, 0x73, 0x0f, 0x7c, 0x99, 0x7b,
	0xe0, 0x6c, 0xe1, 0x35, 0x2e, 0x17, 0x5e, 0xe3, 0xe3, 0xc2, 0x6b, 0xbc, 0x7a, 0xcc, 0xb8, 0x9e,
	0x4c, 0x03, 0x14, 0x4a, 0x81, 0x73, 0xe9, 0x1e, 0x4f, 0x7a, 0x85, 0x38, 0x2e, 0x04, 0x7b, 0xf6,
	0xf5, 0x3d, 0xc5, 0xbf, 0xbc, 0xc5, 0x6f, 0x53, 0xaa, 0x82, 0xad, 0x22, 0xad, 0x0f, 0x7e, 0x06,
	0x00, 0x00, 0xff, 0xff, 0x61, 0x4e, 0x8a, 0x33, 0xa7, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// Deposit defines a method to make a deposit into the module
	Deposit(ctx context.Context, in *MsgDeposit, opts ...grpc.CallOption) (*MsgDepositResponse, error)
	// Withdraw defines a method to withdraw a previously deposited amount
	Withdraw(ctx context.Context, in *MsgWithdraw, opts ...grpc.CallOption) (*MsgWithdrawResponse, error)
	Claim(ctx context.Context, in *MsgClaim, opts ...grpc.CallOption) (*MsgClaimResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) Deposit(ctx context.Context, in *MsgDeposit, opts ...grpc.CallOption) (*MsgDepositResponse, error) {
	out := new(MsgDepositResponse)
	err := c.cc.Invoke(ctx, "/madeinblock.slashrefund.slashrefund.Msg/Deposit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Withdraw(ctx context.Context, in *MsgWithdraw, opts ...grpc.CallOption) (*MsgWithdrawResponse, error) {
	out := new(MsgWithdrawResponse)
	err := c.cc.Invoke(ctx, "/madeinblock.slashrefund.slashrefund.Msg/Withdraw", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Claim(ctx context.Context, in *MsgClaim, opts ...grpc.CallOption) (*MsgClaimResponse, error) {
	out := new(MsgClaimResponse)
	err := c.cc.Invoke(ctx, "/madeinblock.slashrefund.slashrefund.Msg/Claim", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// Deposit defines a method to make a deposit into the module
	Deposit(context.Context, *MsgDeposit) (*MsgDepositResponse, error)
	// Withdraw defines a method to withdraw a previously deposited amount
	Withdraw(context.Context, *MsgWithdraw) (*MsgWithdrawResponse, error)
	Claim(context.Context, *MsgClaim) (*MsgClaimResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) Deposit(ctx context.Context, req *MsgDeposit) (*MsgDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}
func (*UnimplementedMsgServer) Withdraw(ctx context.Context, req *MsgWithdraw) (*MsgWithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (*UnimplementedMsgServer) Claim(ctx context.Context, req *MsgClaim) (*MsgClaimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Claim not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeposit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/madeinblock.slashrefund.slashrefund.Msg/Deposit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Deposit(ctx, req.(*MsgDeposit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdraw)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/madeinblock.slashrefund.slashrefund.Msg/Withdraw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Withdraw(ctx, req.(*MsgWithdraw))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Claim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaim)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Claim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/madeinblock.slashrefund.slashrefund.Msg/Claim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Claim(ctx, req.(*MsgClaim))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "madeinblock.slashrefund.slashrefund.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deposit",
			Handler:    _Msg_Deposit_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _Msg_Withdraw_Handler,
		},
		{
			MethodName: "Claim",
			Handler:    _Msg_Claim_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "slashrefund/tx.proto",
}

func (m *MsgDeposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DepositorAddress) > 0 {
		i -= len(m.DepositorAddress)
		copy(dAtA[i:], m.DepositorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DepositorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDepositResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDepositResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDepositResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgWithdraw) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdraw) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdraw) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Amount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DepositorAddress) > 0 {
		i -= len(m.DepositorAddress)
		copy(dAtA[i:], m.DepositorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DepositorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n3, err3 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.CompletionTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.CompletionTime):])
	if err3 != nil {
		return 0, err3
	}
	i -= n3
	i = encodeVarintTx(dAtA, i, uint64(n3))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *MsgClaim) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaim) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaim) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgClaimResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgDeposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DepositorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgDepositResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgWithdraw) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DepositorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgWithdrawResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.CompletionTime)
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgClaim) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgClaimResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgDeposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgDeposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgDepositResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgDepositResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDepositResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgWithdraw) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgWithdraw: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdraw: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgWithdrawResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgWithdrawResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CompletionTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.CompletionTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgClaim) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgClaim: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaim: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgClaimResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgClaimResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
