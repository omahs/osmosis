// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: osmosis/launchpad/v1/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
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

// Params holds parameters for the launchpad module
type Params struct {
	SaleCreationFee               github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=sale_creation_fee,json=saleCreationFee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"sale_creation_fee" yaml:"sale_creation_fee"`
	MinimumDurationUntilStartTime time.Duration                            `protobuf:"bytes,2,opt,name=minimumDurationUntilStartTime,proto3,stdduration" json:"minimumDurationUntilStartTime"`
	MinimumSaleDuration           time.Duration                            `protobuf:"bytes,3,opt,name=minimumSaleDuration,proto3,stdduration" json:"minimumSaleDuration"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_ac6d83bf0ba122b1, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetSaleCreationFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.SaleCreationFee
	}
	return nil
}

func (m *Params) GetMinimumDurationUntilStartTime() time.Duration {
	if m != nil {
		return m.MinimumDurationUntilStartTime
	}
	return 0
}

func (m *Params) GetMinimumSaleDuration() time.Duration {
	if m != nil {
		return m.MinimumSaleDuration
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "osmosis.launchpad.v1.Params")
}

func init() { proto.RegisterFile("osmosis/launchpad/v1/params.proto", fileDescriptor_ac6d83bf0ba122b1) }

var fileDescriptor_ac6d83bf0ba122b1 = []byte{
	// 389 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x6a, 0xea, 0x40,
	0x18, 0x85, 0x13, 0x05, 0xb9, 0xe4, 0x2e, 0x2e, 0x37, 0xd7, 0x45, 0x14, 0x6e, 0x62, 0x5d, 0xb9,
	0x71, 0x86, 0xb4, 0x85, 0x42, 0x97, 0x5a, 0xba, 0x2a, 0xa5, 0x68, 0xdd, 0x74, 0x23, 0x93, 0x38,
	0xc6, 0xa1, 0x99, 0x4c, 0xc8, 0x4c, 0x42, 0x7d, 0x8b, 0x42, 0x37, 0x7d, 0x86, 0x3e, 0x89, 0x4b,
	0x97, 0x5d, 0xa9, 0xe8, 0x1b, 0xf4, 0x09, 0x4a, 0x26, 0x93, 0x22, 0xb6, 0x94, 0xae, 0x92, 0x9f,
	0xf3, 0xfd, 0x67, 0xce, 0x19, 0xc6, 0x38, 0x62, 0x9c, 0x32, 0x4e, 0x38, 0x0c, 0x51, 0x1a, 0xf9,
	0xb3, 0x18, 0x4d, 0x60, 0xe6, 0xc2, 0x18, 0x25, 0x88, 0x72, 0x10, 0x27, 0x4c, 0x30, 0xb3, 0xae,
	0x10, 0xf0, 0x81, 0x80, 0xcc, 0x6d, 0xd6, 0x03, 0x16, 0x30, 0x09, 0xc0, 0xfc, 0xaf, 0x60, 0x9b,
	0x8d, 0x80, 0xb1, 0x20, 0xc4, 0x50, 0x4e, 0x5e, 0x3a, 0x85, 0x28, 0x9a, 0x97, 0x92, 0x2f, 0x7d,
	0xc6, 0xc5, 0x4e, 0x31, 0x28, 0xc9, 0x2e, 0x26, 0xe8, 0x21, 0x8e, 0x61, 0xe6, 0x7a, 0x58, 0x20,
	0x17, 0xfa, 0x8c, 0x44, 0xa5, 0x7e, 0xe8, 0x3a, 0x49, 0x13, 0x24, 0x08, 0x2b, 0x75, 0xe7, 0x50,
	0x17, 0x84, 0x62, 0x2e, 0x10, 0x8d, 0x0b, 0xa0, 0xbd, 0xa9, 0x18, 0xb5, 0x1b, 0xd9, 0xc9, 0x7c,
	0xd2, 0x8d, 0xbf, 0x1c, 0x85, 0x78, 0xec, 0x27, 0x58, 0x7a, 0x8c, 0xa7, 0x18, 0x5b, 0x7a, 0xab,
	0xda, 0xf9, 0x7d, 0xdc, 0x00, 0x2a, 0x56, 0x1e, 0x04, 0xa8, 0x20, 0xa0, 0xcf, 0x48, 0xd4, 0xbb,
	0x5a, 0xac, 0x1c, 0xed, 0x6d, 0xe5, 0x58, 0x73, 0x44, 0xc3, 0xf3, 0xf6, 0x27, 0x87, 0xf6, 0xcb,
	0xda, 0xe9, 0x04, 0x44, 0xcc, 0x52, 0x0f, 0xf8, 0x8c, 0xaa, 0x7e, 0xea, 0xd3, 0xe5, 0x93, 0x7b,
	0x28, 0xe6, 0x31, 0xe6, 0xd2, 0x8c, 0x0f, 0xfe, 0xe4, 0xfb, 0x7d, 0xb5, 0x7e, 0x89, 0xb1, 0x49,
	0x8c, 0xff, 0x94, 0x44, 0x84, 0xa6, 0xf4, 0x42, 0x55, 0x1b, 0x45, 0x82, 0x84, 0x43, 0x81, 0x12,
	0x71, 0x4b, 0x28, 0xb6, 0x2a, 0x2d, 0x5d, 0x06, 0x2c, 0x9a, 0x82, 0xb2, 0x29, 0x28, 0xf1, 0xde,
	0xaf, 0x3c, 0xe0, 0xf3, 0xda, 0xd1, 0x07, 0xdf, 0x3b, 0x99, 0x23, 0xe3, 0x9f, 0x02, 0x86, 0x28,
	0xc4, 0x25, 0x64, 0x55, 0x7f, 0x7e, 0xc0, 0x57, 0xfb, 0xbd, 0xeb, 0xc5, 0xd6, 0xd6, 0x97, 0x5b,
	0x5b, 0xdf, 0x6c, 0x6d, 0xfd, 0x71, 0x67, 0x6b, 0xcb, 0x9d, 0xad, 0xbd, 0xee, 0x6c, 0xed, 0xee,
	0x74, 0xef, 0x5a, 0xd4, 0x53, 0xea, 0x86, 0xc8, 0xe3, 0xe5, 0x00, 0xb3, 0x33, 0xf8, 0xb0, 0xf7,
	0xfe, 0xe4, 0x45, 0x79, 0x35, 0x99, 0xe0, 0xe4, 0x3d, 0x00, 0x00, 0xff, 0xff, 0x91, 0x33, 0x37,
	0xb4, 0xa1, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.MinimumSaleDuration, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.MinimumSaleDuration):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x1a
	n2, err2 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.MinimumDurationUntilStartTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.MinimumDurationUntilStartTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if len(m.SaleCreationFee) > 0 {
		for iNdEx := len(m.SaleCreationFee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SaleCreationFee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SaleCreationFee) > 0 {
		for _, e := range m.SaleCreationFee {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.MinimumDurationUntilStartTime)
	n += 1 + l + sovParams(uint64(l))
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.MinimumSaleDuration)
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SaleCreationFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SaleCreationFee = append(m.SaleCreationFee, types.Coin{})
			if err := m.SaleCreationFee[len(m.SaleCreationFee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumDurationUntilStartTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.MinimumDurationUntilStartTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumSaleDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.MinimumSaleDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)