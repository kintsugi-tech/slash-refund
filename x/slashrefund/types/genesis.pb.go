// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: slashrefund/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

// GenesisState defines the slashrefund module's genesis state.
type GenesisState struct {
	Params                Params             `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	DepositList           []Deposit          `protobuf:"bytes,2,rep,name=depositList,proto3" json:"depositList"`
	UnbondingDepositList  []UnbondingDeposit `protobuf:"bytes,3,rep,name=unbondingDepositList,proto3" json:"unbondingDepositList"`
	UnbondingDepositCount uint64             `protobuf:"varint,4,opt,name=unbondingDepositCount,proto3" json:"unbondingDepositCount,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_55476ec977a1d0a4, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetDepositList() []Deposit {
	if m != nil {
		return m.DepositList
	}
	return nil
}

func (m *GenesisState) GetUnbondingDepositList() []UnbondingDeposit {
	if m != nil {
		return m.UnbondingDepositList
	}
	return nil
}

func (m *GenesisState) GetUnbondingDepositCount() uint64 {
	if m != nil {
		return m.UnbondingDepositCount
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "madeinblock.slashrefund.slashrefund.GenesisState")
}

func init() { proto.RegisterFile("slashrefund/genesis.proto", fileDescriptor_55476ec977a1d0a4) }

var fileDescriptor_55476ec977a1d0a4 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0xce, 0x49, 0x2c,
	0xce, 0x28, 0x4a, 0x4d, 0x2b, 0xcd, 0x4b, 0xd1, 0x4f, 0x4f, 0xcd, 0x4b, 0x2d, 0xce, 0x2c, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x52, 0xce, 0x4d, 0x4c, 0x49, 0xcd, 0xcc, 0x4b, 0xca, 0xc9,
	0x4f, 0xce, 0xd6, 0x43, 0x52, 0x86, 0xcc, 0x96, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0xab, 0xd7,
	0x07, 0xb1, 0x20, 0x5a, 0xa5, 0x24, 0x90, 0x4d, 0x2d, 0x48, 0x2c, 0x4a, 0xcc, 0x85, 0x1a, 0x2a,
	0x85, 0x62, 0x5f, 0x4a, 0x6a, 0x41, 0x7e, 0x71, 0x66, 0x09, 0x54, 0x4a, 0x19, 0x59, 0xaa, 0x34,
	0x2f, 0x29, 0x3f, 0x2f, 0x25, 0x33, 0x2f, 0x3d, 0x1e, 0x45, 0x91, 0xd2, 0x35, 0x26, 0x2e, 0x1e,
	0x77, 0x88, 0x33, 0x83, 0x4b, 0x12, 0x4b, 0x52, 0x85, 0x3c, 0xb9, 0xd8, 0x20, 0x16, 0x48, 0x30,
	0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x69, 0xeb, 0x11, 0xe1, 0x6c, 0xbd, 0x00, 0xb0, 0x16, 0x27, 0x96,
	0x13, 0xf7, 0xe4, 0x19, 0x82, 0xa0, 0x06, 0x08, 0x85, 0x70, 0x71, 0x43, 0x2d, 0xf3, 0xc9, 0x2c,
	0x2e, 0x91, 0x60, 0x52, 0x60, 0xd6, 0xe0, 0x36, 0xd2, 0x21, 0xca, 0x3c, 0x17, 0x88, 0x3e, 0xa8,
	0x81, 0xc8, 0xc6, 0x08, 0xe5, 0x73, 0x89, 0xc0, 0x3d, 0xe3, 0x82, 0x64, 0x3c, 0x33, 0xd8, 0x78,
	0x53, 0xa2, 0x8c, 0x0f, 0x45, 0x33, 0x00, 0x6a, 0x0f, 0x56, 0x83, 0x85, 0x4c, 0xb8, 0x44, 0xd1,
	0xc5, 0x9d, 0xf3, 0x4b, 0xf3, 0x4a, 0x24, 0x58, 0x14, 0x18, 0x35, 0x58, 0x82, 0xb0, 0x4b, 0x3a,
	0x05, 0x9f, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e,
	0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x65, 0x7a, 0x66, 0x49,
	0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x3e, 0xc8, 0xb1, 0xba, 0x99, 0x79, 0xba, 0x60, 0xe7,
	0xea, 0x83, 0x9d, 0xa8, 0x0b, 0x8d, 0xb1, 0x0a, 0x7d, 0xe4, 0xf8, 0x2b, 0xa9, 0x2c, 0x48, 0x2d,
	0x4e, 0x62, 0x03, 0x47, 0x9a, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x5d, 0xa3, 0x92, 0x66,
	0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.UnbondingDepositCount != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.UnbondingDepositCount))
		i--
		dAtA[i] = 0x20
	}
	if len(m.UnbondingDepositList) > 0 {
		for iNdEx := len(m.UnbondingDepositList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UnbondingDepositList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.DepositList) > 0 {
		for iNdEx := len(m.DepositList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.DepositList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.DepositList) > 0 {
		for _, e := range m.DepositList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.UnbondingDepositList) > 0 {
		for _, e := range m.UnbondingDepositList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.UnbondingDepositCount != 0 {
		n += 1 + sovGenesis(uint64(m.UnbondingDepositCount))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DepositList = append(m.DepositList, Deposit{})
			if err := m.DepositList[len(m.DepositList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingDepositList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UnbondingDepositList = append(m.UnbondingDepositList, UnbondingDeposit{})
			if err := m.UnbondingDepositList[len(m.UnbondingDepositList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingDepositCount", wireType)
			}
			m.UnbondingDepositCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnbondingDepositCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
