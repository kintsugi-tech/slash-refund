// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: slashrefund/unbonding_deposit.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "github.com/regen-network/cosmos-proto"
	_ "google.golang.org/protobuf/types/known/durationpb"
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

type UnbondingDeposit struct {
	Id               uint64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UnbondingStart   time.Time  `protobuf:"bytes,2,opt,name=unbondingStart,proto3,stdtime" json:"unbondingStart"`
	DepositorAddress string     `protobuf:"bytes,3,opt,name=depositorAddress,proto3" json:"depositorAddress,omitempty"`
	ValidatorAddress string     `protobuf:"bytes,4,opt,name=validatorAddress,proto3" json:"validatorAddress,omitempty"`
	Balance          types.Coin `protobuf:"bytes,5,opt,name=balance,proto3" json:"balance"`
}

func (m *UnbondingDeposit) Reset()         { *m = UnbondingDeposit{} }
func (m *UnbondingDeposit) String() string { return proto.CompactTextString(m) }
func (*UnbondingDeposit) ProtoMessage()    {}
func (*UnbondingDeposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_00b62f420100ad0b, []int{0}
}
func (m *UnbondingDeposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnbondingDeposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnbondingDeposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnbondingDeposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnbondingDeposit.Merge(m, src)
}
func (m *UnbondingDeposit) XXX_Size() int {
	return m.Size()
}
func (m *UnbondingDeposit) XXX_DiscardUnknown() {
	xxx_messageInfo_UnbondingDeposit.DiscardUnknown(m)
}

var xxx_messageInfo_UnbondingDeposit proto.InternalMessageInfo

func (m *UnbondingDeposit) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UnbondingDeposit) GetUnbondingStart() time.Time {
	if m != nil {
		return m.UnbondingStart
	}
	return time.Time{}
}

func (m *UnbondingDeposit) GetDepositorAddress() string {
	if m != nil {
		return m.DepositorAddress
	}
	return ""
}

func (m *UnbondingDeposit) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *UnbondingDeposit) GetBalance() types.Coin {
	if m != nil {
		return m.Balance
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*UnbondingDeposit)(nil), "madeinblock.slashrefund.slashrefund.UnbondingDeposit")
}

func init() {
	proto.RegisterFile("slashrefund/unbonding_deposit.proto", fileDescriptor_00b62f420100ad0b)
}

var fileDescriptor_00b62f420100ad0b = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xbf, 0xae, 0xd3, 0x30,
	0x18, 0xc5, 0xe3, 0x50, 0xfe, 0x05, 0xe9, 0xea, 0x2a, 0xba, 0x43, 0xda, 0x21, 0xad, 0xb8, 0x4b,
	0x97, 0xd8, 0xba, 0x30, 0xdd, 0x91, 0x70, 0x47, 0xa6, 0x14, 0x16, 0x96, 0xca, 0x8e, 0xdd, 0xd4,
	0x22, 0xf1, 0x17, 0xc5, 0x4e, 0x45, 0xdf, 0xa2, 0x0f, 0x83, 0xc4, 0x2b, 0x74, 0xac, 0x98, 0x98,
	0x00, 0xb5, 0x2f, 0x82, 0x92, 0x38, 0x28, 0x84, 0x81, 0xcd, 0xd6, 0xcf, 0xe7, 0xe8, 0x7c, 0xe7,
	0xb3, 0x77, 0xab, 0x73, 0xaa, 0xb7, 0x95, 0xd8, 0xd4, 0x8a, 0x93, 0x5a, 0x31, 0x50, 0x5c, 0xaa,
	0x6c, 0xcd, 0x45, 0x09, 0x5a, 0x1a, 0x5c, 0x56, 0x60, 0xc0, 0xbf, 0x2d, 0x28, 0x17, 0x52, 0xb1,
	0x1c, 0xd2, 0x4f, 0x78, 0x20, 0x18, 0x9e, 0x67, 0x37, 0x19, 0x64, 0xd0, 0xbe, 0x27, 0xcd, 0xa9,
	0x93, 0xce, 0xa6, 0x19, 0x40, 0x96, 0x0b, 0xd2, 0xde, 0x58, 0xbd, 0x21, 0x54, 0xed, 0x2d, 0x0a,
	0xc7, 0x88, 0xd7, 0x15, 0x35, 0x12, 0x94, 0xe5, 0xf3, 0x31, 0x37, 0xb2, 0x10, 0xda, 0xd0, 0xa2,
	0xec, 0x0d, 0x52, 0xd0, 0x05, 0x68, 0xc2, 0xa8, 0x16, 0x64, 0x77, 0xc7, 0x84, 0xa1, 0x77, 0x24,
	0x05, 0xd9, 0x1b, 0x4c, 0x3b, 0xbe, 0xee, 0x42, 0x75, 0x97, 0x1e, 0x0d, 0xc7, 0xfe, 0x6b, 0xd8,
	0x97, 0x5f, 0x5d, 0xef, 0xfa, 0x43, 0x5f, 0xc4, 0x43, 0x87, 0xfc, 0x2b, 0xcf, 0x95, 0x3c, 0x40,
	0x0b, 0xb4, 0x9c, 0x24, 0xae, 0xe4, 0xfe, 0x3b, 0xef, 0xea, 0x4f, 0x59, 0x2b, 0x43, 0x2b, 0x13,
	0xb8, 0x0b, 0xb4, 0x7c, 0xf1, 0x6a, 0x86, 0xbb, 0xd0, 0xb8, 0x0f, 0x8d, 0xdf, 0xf7, 0xa1, 0xe3,
	0x67, 0xc7, 0x1f, 0x73, 0xe7, 0xf0, 0x73, 0x8e, 0x92, 0x91, 0xd6, 0x7f, 0xf0, 0xae, 0x6d, 0x06,
	0xa8, 0xde, 0x70, 0x5e, 0x09, 0xad, 0x83, 0x47, 0x0b, 0xb4, 0x7c, 0x1e, 0x07, 0xdf, 0xbe, 0x44,
	0x37, 0x36, 0xb9, 0x25, 0x2b, 0x53, 0x49, 0x95, 0x25, 0xff, 0x28, 0x1a, 0x97, 0x1d, 0xcd, 0x25,
	0xa7, 0x03, 0x97, 0xc9, 0xff, 0x5c, 0xc6, 0x0a, 0xff, 0xde, 0x7b, 0xca, 0x68, 0x4e, 0x55, 0x2a,
	0x82, 0xc7, 0xed, 0x48, 0x53, 0x6c, 0x95, 0x4d, 0xcd, 0xd8, 0xd6, 0x8c, 0xdf, 0x82, 0x54, 0xf1,
	0xa4, 0x99, 0x28, 0xe9, 0xdf, 0xc7, 0xab, 0xe3, 0x39, 0x44, 0xa7, 0x73, 0x88, 0x7e, 0x9d, 0x43,
	0x74, 0xb8, 0x84, 0xce, 0xe9, 0x12, 0x3a, 0xdf, 0x2f, 0xa1, 0xf3, 0xf1, 0x3e, 0x93, 0x66, 0x5b,
	0x33, 0x9c, 0x42, 0x41, 0x9a, 0xbf, 0x14, 0x49, 0x15, 0xb5, 0xbf, 0x89, 0xb4, 0x7b, 0x88, 0xec,
	0x22, 0x3e, 0x93, 0xe1, 0x5a, 0xcc, 0xbe, 0x14, 0x9a, 0x3d, 0x69, 0x9b, 0x7c, 0xfd, 0x3b, 0x00,
	0x00, 0xff, 0xff, 0x27, 0x64, 0x2f, 0xc4, 0xa9, 0x02, 0x00, 0x00,
}

func (m *UnbondingDeposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnbondingDeposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnbondingDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Balance.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.DepositorAddress) > 0 {
		i -= len(m.DepositorAddress)
		copy(dAtA[i:], m.DepositorAddress)
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(len(m.DepositorAddress)))
		i--
		dAtA[i] = 0x1a
	}
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UnbondingStart, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UnbondingStart):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintUnbondingDeposit(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if m.Id != 0 {
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintUnbondingDeposit(dAtA []byte, offset int, v uint64) int {
	offset -= sovUnbondingDeposit(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *UnbondingDeposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovUnbondingDeposit(uint64(m.Id))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UnbondingStart)
	n += 1 + l + sovUnbondingDeposit(uint64(l))
	l = len(m.DepositorAddress)
	if l > 0 {
		n += 1 + l + sovUnbondingDeposit(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovUnbondingDeposit(uint64(l))
	}
	l = m.Balance.Size()
	n += 1 + l + sovUnbondingDeposit(uint64(l))
	return n
}

func sovUnbondingDeposit(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozUnbondingDeposit(x uint64) (n int) {
	return sovUnbondingDeposit(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *UnbondingDeposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowUnbondingDeposit
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
			return fmt.Errorf("proto: UnbondingDeposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnbondingDeposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbondingDeposit
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingStart", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbondingDeposit
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
				return ErrInvalidLengthUnbondingDeposit
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthUnbondingDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UnbondingStart, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbondingDeposit
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
				return ErrInvalidLengthUnbondingDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUnbondingDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbondingDeposit
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
				return ErrInvalidLengthUnbondingDeposit
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthUnbondingDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Balance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowUnbondingDeposit
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
				return ErrInvalidLengthUnbondingDeposit
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthUnbondingDeposit
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Balance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipUnbondingDeposit(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthUnbondingDeposit
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
func skipUnbondingDeposit(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowUnbondingDeposit
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
					return 0, ErrIntOverflowUnbondingDeposit
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
					return 0, ErrIntOverflowUnbondingDeposit
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
				return 0, ErrInvalidLengthUnbondingDeposit
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupUnbondingDeposit
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthUnbondingDeposit
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthUnbondingDeposit        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowUnbondingDeposit          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupUnbondingDeposit = fmt.Errorf("proto: unexpected end of group")
)
