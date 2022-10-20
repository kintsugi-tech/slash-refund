// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: slashrefund/deposit_pool.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	_ "github.com/tendermint/tendermint/proto/tendermint/types"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// TODO: to account for more than one token, Tokens and Shares must be a struct.
type DepositPool struct {
	OperatorAddress string                                 `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
	Tokens          types.Coin                             `protobuf:"bytes,2,opt,name=tokens,proto3" json:"tokens"`
	Shares          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares" yaml:"depositor_shares"`
}

func (m *DepositPool) Reset()         { *m = DepositPool{} }
func (m *DepositPool) String() string { return proto.CompactTextString(m) }
func (*DepositPool) ProtoMessage()    {}
func (*DepositPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_6b70b28e450891ac, []int{0}
}
func (m *DepositPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DepositPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DepositPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DepositPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DepositPool.Merge(m, src)
}
func (m *DepositPool) XXX_Size() int {
	return m.Size()
}
func (m *DepositPool) XXX_DiscardUnknown() {
	xxx_messageInfo_DepositPool.DiscardUnknown(m)
}

var xxx_messageInfo_DepositPool proto.InternalMessageInfo

func (m *DepositPool) GetOperatorAddress() string {
	if m != nil {
		return m.OperatorAddress
	}
	return ""
}

func (m *DepositPool) GetTokens() types.Coin {
	if m != nil {
		return m.Tokens
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*DepositPool)(nil), "madeinblock.slashrefund.slashrefund.DepositPool")
}

func init() { proto.RegisterFile("slashrefund/deposit_pool.proto", fileDescriptor_6b70b28e450891ac) }

var fileDescriptor_6b70b28e450891ac = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x52, 0xb1, 0x6a, 0xeb, 0x30,
	0x14, 0xb5, 0xdf, 0x7b, 0x04, 0x9e, 0x33, 0xb4, 0x98, 0x40, 0x9d, 0x50, 0x94, 0x90, 0x42, 0xc9,
	0x62, 0x8b, 0xb4, 0x43, 0x69, 0xb7, 0x26, 0x59, 0xba, 0x95, 0x64, 0xeb, 0x62, 0x64, 0x5b, 0x71,
	0x44, 0x6c, 0x5d, 0x23, 0x29, 0xa5, 0xf9, 0x8b, 0x7e, 0x4c, 0x3f, 0x22, 0x63, 0xe8, 0x54, 0x3a,
	0x84, 0x92, 0xec, 0x1d, 0xfa, 0x05, 0xc5, 0x96, 0x02, 0x5e, 0xec, 0x7b, 0x75, 0xee, 0xb9, 0xe7,
	0x5e, 0x1d, 0x39, 0x48, 0x66, 0x44, 0x2e, 0x04, 0x9d, 0xaf, 0x78, 0x82, 0x13, 0x5a, 0x80, 0x64,
	0x2a, 0x2c, 0x00, 0xb2, 0xa0, 0x10, 0xa0, 0xc0, 0xbd, 0xc8, 0x49, 0x42, 0x19, 0x8f, 0x32, 0x88,
	0x97, 0x41, 0xad, 0xb6, 0x1e, 0x77, 0x5a, 0x29, 0xa4, 0x50, 0xd5, 0xe3, 0x32, 0xd2, 0xd4, 0x4e,
	0x3b, 0x05, 0x48, 0x33, 0x8a, 0xab, 0x2c, 0x5a, 0xcd, 0x31, 0xe1, 0xeb, 0x23, 0x14, 0x83, 0xcc,
	0x41, 0x86, 0x9a, 0xa3, 0x13, 0x03, 0x21, 0x9d, 0xe1, 0x88, 0x48, 0x8a, 0x9f, 0x87, 0x11, 0x55,
	0x64, 0x88, 0x63, 0x60, 0xdc, 0xe0, 0xe7, 0x8a, 0xf2, 0x84, 0x8a, 0x9c, 0x71, 0x85, 0xd5, 0xba,
	0xa0, 0x52, 0x7f, 0x35, 0xda, 0xff, 0xb6, 0x9d, 0xe6, 0x44, 0x6f, 0xf1, 0x08, 0x90, 0xb9, 0x63,
	0xe7, 0x14, 0x0a, 0x2a, 0x88, 0x02, 0x11, 0x92, 0x24, 0x11, 0x54, 0x4a, 0xcf, 0xee, 0xd9, 0x83,
	0xff, 0x23, 0xef, 0xfd, 0xcd, 0x6f, 0x19, 0xe5, 0x7b, 0x8d, 0xcc, 0x94, 0x60, 0x3c, 0x9d, 0x9e,
	0x1c, 0x19, 0xe6, 0xd8, 0xbd, 0x71, 0x1a, 0x0a, 0x96, 0x94, 0x4b, 0xef, 0x4f, 0xcf, 0x1e, 0x34,
	0xaf, 0xda, 0x81, 0xe1, 0x95, 0x33, 0x06, 0x66, 0xc6, 0x60, 0x0c, 0x8c, 0x8f, 0xfe, 0x6d, 0x76,
	0x5d, 0x6b, 0x6a, 0xca, 0x5d, 0xe2, 0x34, 0xe4, 0x82, 0x08, 0x2a, 0xbd, 0xbf, 0x95, 0xe6, 0x43,
	0x89, 0x7e, 0xee, 0xba, 0x97, 0x29, 0x53, 0x8b, 0x55, 0x14, 0xc4, 0x90, 0x9b, 0xe5, 0xcd, 0xcf,
	0x97, 0xc9, 0xd2, 0xec, 0x33, 0xa1, 0xf1, 0xcf, 0xae, 0x7b, 0xb6, 0x26, 0x79, 0x76, 0xd7, 0x37,
	0xc6, 0x80, 0x08, 0x75, 0xbf, 0xfe, 0xd4, 0x34, 0x1e, 0xcd, 0x36, 0x7b, 0x64, 0x6f, 0xf7, 0xc8,
	0xfe, 0xda, 0x23, 0xfb, 0xf5, 0x80, 0xac, 0xed, 0x01, 0x59, 0x1f, 0x07, 0x64, 0x3d, 0xdd, 0xd6,
	0x44, 0x4a, 0x13, 0x7d, 0xc6, 0xfd, 0xca, 0x46, 0x5c, 0x59, 0xe7, 0x1b, 0xcf, 0x5f, 0x70, 0xfd,
	0x05, 0x54, 0xda, 0x51, 0xa3, 0xba, 0xcc, 0xeb, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5c, 0x9d,
	0x30, 0x3b, 0x1d, 0x02, 0x00, 0x00,
}

func (m *DepositPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DepositPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DepositPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Shares.Size()
		i -= size
		if _, err := m.Shares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDepositPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.Tokens.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintDepositPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.OperatorAddress) > 0 {
		i -= len(m.OperatorAddress)
		copy(dAtA[i:], m.OperatorAddress)
		i = encodeVarintDepositPool(dAtA, i, uint64(len(m.OperatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDepositPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovDepositPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DepositPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OperatorAddress)
	if l > 0 {
		n += 1 + l + sovDepositPool(uint64(l))
	}
	l = m.Tokens.Size()
	n += 1 + l + sovDepositPool(uint64(l))
	l = m.Shares.Size()
	n += 1 + l + sovDepositPool(uint64(l))
	return n
}

func sovDepositPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDepositPool(x uint64) (n int) {
	return sovDepositPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DepositPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDepositPool
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
			return fmt.Errorf("proto: DepositPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DepositPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OperatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDepositPool
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
				return ErrInvalidLengthDepositPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDepositPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OperatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tokens", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDepositPool
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
				return ErrInvalidLengthDepositPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthDepositPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tokens.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDepositPool
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
				return ErrInvalidLengthDepositPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDepositPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDepositPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDepositPool
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
func skipDepositPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDepositPool
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
					return 0, ErrIntOverflowDepositPool
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
					return 0, ErrIntOverflowDepositPool
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
				return 0, ErrInvalidLengthDepositPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDepositPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDepositPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDepositPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDepositPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDepositPool = fmt.Errorf("proto: unexpected end of group")
)