// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: fin/v1/finpb.proto

package finv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Transaction_Category int32

const (
	Transaction_UNKNOWN Transaction_Category = 0
	Transaction_CASH    Transaction_Category = 10
	Transaction_ING     Transaction_Category = 11
	Transaction_WISE    Transaction_Category = 12
	Transaction_REVOLUT Transaction_Category = 13
	Transaction_DEGIRO  Transaction_Category = 14
	Transaction_BITTREX Transaction_Category = 15
	Transaction_HSBC    Transaction_Category = 16
	Transaction_NOVUS   Transaction_Category = 17
	// = 18;
	Transaction_FOREX         Transaction_Category = 19
	Transaction_SALARY        Transaction_Category = 20
	Transaction_ALLOWANCE     Transaction_Category = 21
	Transaction_CREDIT        Transaction_Category = 22
	Transaction_IN_OTHER      Transaction_Category = 23
	Transaction_FOOD          Transaction_Category = 30
	Transaction_TRANSPORT     Transaction_Category = 31
	Transaction_HOUSING       Transaction_Category = 32
	Transaction_TECH          Transaction_Category = 33
	Transaction_FINANCE       Transaction_Category = 34
	Transaction_HEALTH        Transaction_Category = 35
	Transaction_CLOTHING      Transaction_Category = 36
	Transaction_ENTERTAINMENT Transaction_Category = 37
	Transaction_PERSONAL      Transaction_Category = 38
	Transaction_TRAVEL        Transaction_Category = 39
	Transaction_SOCIAL        Transaction_Category = 40
	Transaction_EDUCATION     Transaction_Category = 41
)

// Enum value maps for Transaction_Category.
var (
	Transaction_Category_name = map[int32]string{
		0:  "UNKNOWN",
		10: "CASH",
		11: "ING",
		12: "WISE",
		13: "REVOLUT",
		14: "DEGIRO",
		15: "BITTREX",
		16: "HSBC",
		17: "NOVUS",
		19: "FOREX",
		20: "SALARY",
		21: "ALLOWANCE",
		22: "CREDIT",
		23: "IN_OTHER",
		30: "FOOD",
		31: "TRANSPORT",
		32: "HOUSING",
		33: "TECH",
		34: "FINANCE",
		35: "HEALTH",
		36: "CLOTHING",
		37: "ENTERTAINMENT",
		38: "PERSONAL",
		39: "TRAVEL",
		40: "SOCIAL",
		41: "EDUCATION",
	}
	Transaction_Category_value = map[string]int32{
		"UNKNOWN":       0,
		"CASH":          10,
		"ING":           11,
		"WISE":          12,
		"REVOLUT":       13,
		"DEGIRO":        14,
		"BITTREX":       15,
		"HSBC":          16,
		"NOVUS":         17,
		"FOREX":         19,
		"SALARY":        20,
		"ALLOWANCE":     21,
		"CREDIT":        22,
		"IN_OTHER":      23,
		"FOOD":          30,
		"TRANSPORT":     31,
		"HOUSING":       32,
		"TECH":          33,
		"FINANCE":       34,
		"HEALTH":        35,
		"CLOTHING":      36,
		"ENTERTAINMENT": 37,
		"PERSONAL":      38,
		"TRAVEL":        39,
		"SOCIAL":        40,
		"EDUCATION":     41,
	}
)

func (x Transaction_Category) Enum() *Transaction_Category {
	p := new(Transaction_Category)
	*p = x
	return p
}

func (x Transaction_Category) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Transaction_Category) Descriptor() protoreflect.EnumDescriptor {
	return file_fin_v1_finpb_proto_enumTypes[0].Descriptor()
}

func (Transaction_Category) Type() protoreflect.EnumType {
	return &file_fin_v1_finpb_proto_enumTypes[0]
}

func (x Transaction_Category) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Transaction_Category.Descriptor instead.
func (Transaction_Category) EnumDescriptor() ([]byte, []int) {
	return file_fin_v1_finpb_proto_rawDescGZIP(), []int{2, 0}
}

type All struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Currency string   `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	Months   []*Month `protobuf:"bytes,3,rep,name=months,proto3" json:"months,omitempty"`
}

func (x *All) Reset() {
	*x = All{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fin_v1_finpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *All) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*All) ProtoMessage() {}

func (x *All) ProtoReflect() protoreflect.Message {
	mi := &file_fin_v1_finpb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use All.ProtoReflect.Descriptor instead.
func (*All) Descriptor() ([]byte, []int) {
	return file_fin_v1_finpb_proto_rawDescGZIP(), []int{0}
}

func (x *All) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *All) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *All) GetMonths() []*Month {
	if x != nil {
		return x.Months
	}
	return nil
}

type Month struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Year         int32          `protobuf:"varint,1,opt,name=year,proto3" json:"year,omitempty"`
	Month        int32          `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Transactions []*Transaction `protobuf:"bytes,3,rep,name=transactions,proto3" json:"transactions,omitempty"`
}

func (x *Month) Reset() {
	*x = Month{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fin_v1_finpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Month) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Month) ProtoMessage() {}

func (x *Month) ProtoReflect() protoreflect.Message {
	mi := &file_fin_v1_finpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Month.ProtoReflect.Descriptor instead.
func (*Month) Descriptor() ([]byte, []int) {
	return file_fin_v1_finpb_proto_rawDescGZIP(), []int{1}
}

func (x *Month) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *Month) GetMonth() int32 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *Month) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Amount int64                `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	Src    Transaction_Category `protobuf:"varint,2,opt,name=src,proto3,enum=apis.fin.v1.Transaction_Category" json:"src,omitempty"`
	Dst    Transaction_Category `protobuf:"varint,3,opt,name=dst,proto3,enum=apis.fin.v1.Transaction_Category" json:"dst,omitempty"`
	Note   string               `protobuf:"bytes,4,opt,name=note,proto3" json:"note,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_fin_v1_finpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_fin_v1_finpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_fin_v1_finpb_proto_rawDescGZIP(), []int{2}
}

func (x *Transaction) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Transaction) GetSrc() Transaction_Category {
	if x != nil {
		return x.Src
	}
	return Transaction_UNKNOWN
}

func (x *Transaction) GetDst() Transaction_Category {
	if x != nil {
		return x.Dst
	}
	return Transaction_UNKNOWN
}

func (x *Transaction) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

var File_fin_v1_finpb_proto protoreflect.FileDescriptor

var file_fin_v1_finpb_proto_rawDesc = []byte{
	0x0a, 0x12, 0x66, 0x69, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x69, 0x6e, 0x70, 0x62, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x66, 0x69, 0x6e, 0x2e, 0x76,
	0x31, 0x22, 0x61, 0x0a, 0x03, 0x41, 0x6c, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x2a, 0x0a, 0x06, 0x6d, 0x6f, 0x6e, 0x74,
	0x68, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e,
	0x66, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x52, 0x06, 0x6d, 0x6f,
	0x6e, 0x74, 0x68, 0x73, 0x22, 0x6f, 0x0a, 0x05, 0x4d, 0x6f, 0x6e, 0x74, 0x68, 0x12, 0x12, 0x0a,
	0x04, 0x79, 0x65, 0x61, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x79, 0x65, 0x61,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x12, 0x3c, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x61, 0x70, 0x69, 0x73, 0x2e, 0x66, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0xf4, 0x03, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x33, 0x0a,
	0x03, 0x73, 0x72, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x61, 0x70, 0x69,
	0x73, 0x2e, 0x66, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x52, 0x03, 0x73,
	0x72, 0x63, 0x12, 0x33, 0x0a, 0x03, 0x64, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x21, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x66, 0x69, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x52, 0x03, 0x64, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x22, 0xce, 0x02, 0x0a, 0x08,
	0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x43, 0x41, 0x53, 0x48, 0x10, 0x0a, 0x12,
	0x07, 0x0a, 0x03, 0x49, 0x4e, 0x47, 0x10, 0x0b, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x49, 0x53, 0x45,
	0x10, 0x0c, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x56, 0x4f, 0x4c, 0x55, 0x54, 0x10, 0x0d, 0x12,
	0x0a, 0x0a, 0x06, 0x44, 0x45, 0x47, 0x49, 0x52, 0x4f, 0x10, 0x0e, 0x12, 0x0b, 0x0a, 0x07, 0x42,
	0x49, 0x54, 0x54, 0x52, 0x45, 0x58, 0x10, 0x0f, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x53, 0x42, 0x43,
	0x10, 0x10, 0x12, 0x09, 0x0a, 0x05, 0x4e, 0x4f, 0x56, 0x55, 0x53, 0x10, 0x11, 0x12, 0x09, 0x0a,
	0x05, 0x46, 0x4f, 0x52, 0x45, 0x58, 0x10, 0x13, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x41, 0x4c, 0x41,
	0x52, 0x59, 0x10, 0x14, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x4c, 0x4c, 0x4f, 0x57, 0x41, 0x4e, 0x43,
	0x45, 0x10, 0x15, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x52, 0x45, 0x44, 0x49, 0x54, 0x10, 0x16, 0x12,
	0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x5f, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x17, 0x12, 0x08, 0x0a,
	0x04, 0x46, 0x4f, 0x4f, 0x44, 0x10, 0x1e, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x52, 0x41, 0x4e, 0x53,
	0x50, 0x4f, 0x52, 0x54, 0x10, 0x1f, 0x12, 0x0b, 0x0a, 0x07, 0x48, 0x4f, 0x55, 0x53, 0x49, 0x4e,
	0x47, 0x10, 0x20, 0x12, 0x08, 0x0a, 0x04, 0x54, 0x45, 0x43, 0x48, 0x10, 0x21, 0x12, 0x0b, 0x0a,
	0x07, 0x46, 0x49, 0x4e, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x22, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x45,
	0x41, 0x4c, 0x54, 0x48, 0x10, 0x23, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4c, 0x4f, 0x54, 0x48, 0x49,
	0x4e, 0x47, 0x10, 0x24, 0x12, 0x11, 0x0a, 0x0d, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x54, 0x41, 0x49,
	0x4e, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x25, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x45, 0x52, 0x53, 0x4f,
	0x4e, 0x41, 0x4c, 0x10, 0x26, 0x12, 0x0a, 0x0a, 0x06, 0x54, 0x52, 0x41, 0x56, 0x45, 0x4c, 0x10,
	0x27, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x4f, 0x43, 0x49, 0x41, 0x4c, 0x10, 0x28, 0x12, 0x0d, 0x0a,
	0x09, 0x45, 0x44, 0x55, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x29, 0x42, 0x95, 0x01, 0x0a,
	0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x66, 0x69, 0x6e, 0x2e, 0x76, 0x31,
	0x42, 0x0a, 0x46, 0x69, 0x6e, 0x70, 0x62, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x28,
	0x67, 0x6f, 0x2e, 0x73, 0x65, 0x61, 0x6e, 0x6b, 0x68, 0x6c, 0x69, 0x61, 0x6f, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x66, 0x69, 0x6e, 0x2f,
	0x76, 0x31, 0x3b, 0x66, 0x69, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x41, 0x46, 0x58, 0xaa, 0x02,
	0x0b, 0x41, 0x70, 0x69, 0x73, 0x2e, 0x46, 0x69, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x41,
	0x70, 0x69, 0x73, 0x5c, 0x46, 0x69, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x41, 0x70, 0x69,
	0x73, 0x5c, 0x46, 0x69, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x41, 0x70, 0x69, 0x73, 0x3a, 0x3a, 0x46, 0x69, 0x6e,
	0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_fin_v1_finpb_proto_rawDescOnce sync.Once
	file_fin_v1_finpb_proto_rawDescData = file_fin_v1_finpb_proto_rawDesc
)

func file_fin_v1_finpb_proto_rawDescGZIP() []byte {
	file_fin_v1_finpb_proto_rawDescOnce.Do(func() {
		file_fin_v1_finpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_fin_v1_finpb_proto_rawDescData)
	})
	return file_fin_v1_finpb_proto_rawDescData
}

var file_fin_v1_finpb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_fin_v1_finpb_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_fin_v1_finpb_proto_goTypes = []interface{}{
	(Transaction_Category)(0), // 0: apis.fin.v1.Transaction.Category
	(*All)(nil),               // 1: apis.fin.v1.All
	(*Month)(nil),             // 2: apis.fin.v1.Month
	(*Transaction)(nil),       // 3: apis.fin.v1.Transaction
}
var file_fin_v1_finpb_proto_depIdxs = []int32{
	2, // 0: apis.fin.v1.All.months:type_name -> apis.fin.v1.Month
	3, // 1: apis.fin.v1.Month.transactions:type_name -> apis.fin.v1.Transaction
	0, // 2: apis.fin.v1.Transaction.src:type_name -> apis.fin.v1.Transaction.Category
	0, // 3: apis.fin.v1.Transaction.dst:type_name -> apis.fin.v1.Transaction.Category
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_fin_v1_finpb_proto_init() }
func file_fin_v1_finpb_proto_init() {
	if File_fin_v1_finpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_fin_v1_finpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*All); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_fin_v1_finpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Month); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_fin_v1_finpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_fin_v1_finpb_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_fin_v1_finpb_proto_goTypes,
		DependencyIndexes: file_fin_v1_finpb_proto_depIdxs,
		EnumInfos:         file_fin_v1_finpb_proto_enumTypes,
		MessageInfos:      file_fin_v1_finpb_proto_msgTypes,
	}.Build()
	File_fin_v1_finpb_proto = out.File
	file_fin_v1_finpb_proto_rawDesc = nil
	file_fin_v1_finpb_proto_goTypes = nil
	file_fin_v1_finpb_proto_depIdxs = nil
}
