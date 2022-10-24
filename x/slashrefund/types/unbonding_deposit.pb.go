// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: slashrefund/unbonding_deposit.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
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

type UnbondingDeposit struct {
	DelegatorAddress      string `protobuf:"bytes,1,opt,name=delegatorAddress,proto3" json:"delegatorAddress,omitempty"`
	ValidatorAddress      string `protobuf:"bytes,2,opt,name=validatorAddress,proto3" json:"validatorAddress,omitempty"`
	UnbondingDepositEntry string `protobuf:"bytes,3,opt,name=unbondingDepositEntry,proto3" json:"unbondingDepositEntry,omitempty"`
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

func (m *UnbondingDeposit) GetDelegatorAddress() string {
	if m != nil {
		return m.DelegatorAddress
	}
	return ""
}

func (m *UnbondingDeposit) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *UnbondingDeposit) GetUnbondingDepositEntry() string {
	if m != nil {
		return m.UnbondingDepositEntry
	}
	return ""
}

func init() {
	proto.RegisterType((*UnbondingDeposit)(nil), "madeinblock.slashrefund.slashrefund.UnbondingDeposit")
}

func init() {
	proto.RegisterFile("slashrefund/unbonding_deposit.proto", fileDescriptor_00b62f420100ad0b)
}

var fileDescriptor_00b62f420100ad0b = []byte{
	// 222 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2e, 0xce, 0x49, 0x2c,
	0xce, 0x28, 0x4a, 0x4d, 0x2b, 0xcd, 0x4b, 0xd1, 0x2f, 0xcd, 0x4b, 0xca, 0xcf, 0x4b, 0xc9, 0xcc,
	0x4b, 0x8f, 0x4f, 0x49, 0x2d, 0xc8, 0x2f, 0xce, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17,
	0x52, 0xce, 0x4d, 0x4c, 0x49, 0xcd, 0xcc, 0x4b, 0xca, 0xc9, 0x4f, 0xce, 0xd6, 0x43, 0xd2, 0x80,
	0xcc, 0x56, 0x5a, 0xc0, 0xc8, 0x25, 0x10, 0x0a, 0x33, 0xc0, 0x05, 0xa2, 0x5f, 0x48, 0x8b, 0x4b,
	0x20, 0x25, 0x35, 0x27, 0x35, 0x3d, 0xb1, 0x24, 0xbf, 0xc8, 0x31, 0x25, 0xa5, 0x28, 0xb5, 0xb8,
	0x58, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0x43, 0x1c, 0xa4, 0xb6, 0x2c, 0x31, 0x27, 0x33,
	0x05, 0x59, 0x2d, 0x13, 0x44, 0x2d, 0xba, 0xb8, 0x90, 0x09, 0x97, 0x68, 0x29, 0x9a, 0x5d, 0xae,
	0x79, 0x25, 0x45, 0x95, 0x12, 0xcc, 0x60, 0x0d, 0xd8, 0x25, 0x9d, 0x82, 0x4f, 0x3c, 0x92, 0x63,
	0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96,
	0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x32, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39,
	0x3f, 0x57, 0x1f, 0xe4, 0x59, 0xdd, 0xcc, 0x3c, 0x5d, 0xb0, 0x77, 0xf5, 0xc1, 0x5e, 0xd4, 0x85,
	0x06, 0x50, 0x85, 0x3e, 0x72, 0x70, 0x95, 0x54, 0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81, 0xc3, 0xc8,
	0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x76, 0x3d, 0xd3, 0x53, 0x4a, 0x01, 0x00, 0x00,
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
	if len(m.UnbondingDepositEntry) > 0 {
		i -= len(m.UnbondingDepositEntry)
		copy(dAtA[i:], m.UnbondingDepositEntry)
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(len(m.UnbondingDepositEntry)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintUnbondingDeposit(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
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
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovUnbondingDeposit(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovUnbondingDeposit(uint64(l))
	}
	l = len(m.UnbondingDepositEntry)
	if l > 0 {
		n += 1 + l + sovUnbondingDeposit(uint64(l))
	}
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
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
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingDepositEntry", wireType)
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
			m.UnbondingDepositEntry = string(dAtA[iNdEx:postIndex])
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
