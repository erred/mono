// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: authd/v1/authd.proto

package authdv1

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

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// host: path regex
	Allowlist map[string]*AllowedPaths `protobuf:"bytes,1,rep,name=allowlist,proto3" json:"allowlist,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// host:
	Tokens   map[string]*Tokens `protobuf:"bytes,2,rep,name=tokens,proto3" json:"tokens,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Htpasswd map[string]string  `protobuf:"bytes,3,rep,name=htpasswd,proto3" json:"htpasswd,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// entire file, lower priority
	HtpasswdFile string `protobuf:"bytes,4,opt,name=htpasswd_file,json=htpasswdFile,proto3" json:"htpasswd_file,omitempty"`
	SessionStore string `protobuf:"bytes,5,opt,name=session_store,json=sessionStore,proto3" json:"session_store,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authd_v1_authd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_authd_v1_authd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_authd_v1_authd_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetAllowlist() map[string]*AllowedPaths {
	if x != nil {
		return x.Allowlist
	}
	return nil
}

func (x *Config) GetTokens() map[string]*Tokens {
	if x != nil {
		return x.Tokens
	}
	return nil
}

func (x *Config) GetHtpasswd() map[string]string {
	if x != nil {
		return x.Htpasswd
	}
	return nil
}

func (x *Config) GetHtpasswdFile() string {
	if x != nil {
		return x.HtpasswdFile
	}
	return ""
}

func (x *Config) GetSessionStore() string {
	if x != nil {
		return x.SessionStore
	}
	return ""
}

type AllowedPaths struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PathRe []string `protobuf:"bytes,1,rep,name=path_re,json=pathRe,proto3" json:"path_re,omitempty"`
}

func (x *AllowedPaths) Reset() {
	*x = AllowedPaths{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authd_v1_authd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllowedPaths) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllowedPaths) ProtoMessage() {}

func (x *AllowedPaths) ProtoReflect() protoreflect.Message {
	mi := &file_authd_v1_authd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllowedPaths.ProtoReflect.Descriptor instead.
func (*AllowedPaths) Descriptor() ([]byte, []int) {
	return file_authd_v1_authd_proto_rawDescGZIP(), []int{1}
}

func (x *AllowedPaths) GetPathRe() []string {
	if x != nil {
		return x.PathRe
	}
	return nil
}

type Tokens struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tokens []*Token `protobuf:"bytes,1,rep,name=tokens,proto3" json:"tokens,omitempty"`
}

func (x *Tokens) Reset() {
	*x = Tokens{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authd_v1_authd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tokens) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tokens) ProtoMessage() {}

func (x *Tokens) ProtoReflect() protoreflect.Message {
	mi := &file_authd_v1_authd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tokens.ProtoReflect.Descriptor instead.
func (*Tokens) Descriptor() ([]byte, []int) {
	return file_authd_v1_authd_proto_rawDescGZIP(), []int{2}
}

func (x *Tokens) GetTokens() []*Token {
	if x != nil {
		return x.Tokens
	}
	return nil
}

type Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// token used in "authorization: Bearer $token"
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	// name used to identify it
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Token) Reset() {
	*x = Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authd_v1_authd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Token) ProtoMessage() {}

func (x *Token) ProtoReflect() protoreflect.Message {
	mi := &file_authd_v1_authd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Token.ProtoReflect.Descriptor instead.
func (*Token) Descriptor() ([]byte, []int) {
	return file_authd_v1_authd_proto_rawDescGZIP(), []int{3}
}

func (x *Token) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *Token) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

var File_authd_v1_authd_proto protoreflect.FileDescriptor

var file_authd_v1_authd_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x64,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x75, 0x74, 0x68, 0x64, 0x2e, 0x76, 0x31,
	0x22, 0xe3, 0x03, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3d, 0x0a, 0x09, 0x61,
	0x6c, 0x6c, 0x6f, 0x77, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x61, 0x75, 0x74, 0x68, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x6c, 0x69, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x09, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x06, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x73,
	0x12, 0x3a, 0x0a, 0x08, 0x68, 0x74, 0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x48, 0x74, 0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x08, 0x68, 0x74, 0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x68, 0x74, 0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x68, 0x74, 0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x1a, 0x54, 0x0a, 0x0e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x6c,
	0x69, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x64, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x64, 0x50, 0x61, 0x74, 0x68,
	0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x4b, 0x0a, 0x0b,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x3b, 0x0a, 0x0d, 0x48, 0x74, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x27, 0x0a, 0x0c, 0x41, 0x6c, 0x6c, 0x6f, 0x77, 0x65,
	0x64, 0x50, 0x61, 0x74, 0x68, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x72,
	0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x70, 0x61, 0x74, 0x68, 0x52, 0x65, 0x22,
	0x31, 0x0a, 0x06, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x12, 0x27, 0x0a, 0x06, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x75, 0x74, 0x68,
	0x64, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x06, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x73, 0x22, 0x2d, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x42, 0x89, 0x01, 0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x64, 0x2e,
	0x76, 0x31, 0x42, 0x0a, 0x41, 0x75, 0x74, 0x68, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x2c, 0x67, 0x6f, 0x2e, 0x73, 0x65, 0x61, 0x6e, 0x6b, 0x68, 0x6c, 0x69, 0x61, 0x6f, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x75,
	0x74, 0x68, 0x64, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x75, 0x74, 0x68, 0x64, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x41, 0x58, 0x58, 0xaa, 0x02, 0x08, 0x41, 0x75, 0x74, 0x68, 0x64, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x08, 0x41, 0x75, 0x74, 0x68, 0x64, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x14, 0x41, 0x75, 0x74,
	0x68, 0x64, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x09, 0x41, 0x75, 0x74, 0x68, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_authd_v1_authd_proto_rawDescOnce sync.Once
	file_authd_v1_authd_proto_rawDescData = file_authd_v1_authd_proto_rawDesc
)

func file_authd_v1_authd_proto_rawDescGZIP() []byte {
	file_authd_v1_authd_proto_rawDescOnce.Do(func() {
		file_authd_v1_authd_proto_rawDescData = protoimpl.X.CompressGZIP(file_authd_v1_authd_proto_rawDescData)
	})
	return file_authd_v1_authd_proto_rawDescData
}

var file_authd_v1_authd_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_authd_v1_authd_proto_goTypes = []interface{}{
	(*Config)(nil),       // 0: authd.v1.Config
	(*AllowedPaths)(nil), // 1: authd.v1.AllowedPaths
	(*Tokens)(nil),       // 2: authd.v1.Tokens
	(*Token)(nil),        // 3: authd.v1.Token
	nil,                  // 4: authd.v1.Config.AllowlistEntry
	nil,                  // 5: authd.v1.Config.TokensEntry
	nil,                  // 6: authd.v1.Config.HtpasswdEntry
}
var file_authd_v1_authd_proto_depIdxs = []int32{
	4, // 0: authd.v1.Config.allowlist:type_name -> authd.v1.Config.AllowlistEntry
	5, // 1: authd.v1.Config.tokens:type_name -> authd.v1.Config.TokensEntry
	6, // 2: authd.v1.Config.htpasswd:type_name -> authd.v1.Config.HtpasswdEntry
	3, // 3: authd.v1.Tokens.tokens:type_name -> authd.v1.Token
	1, // 4: authd.v1.Config.AllowlistEntry.value:type_name -> authd.v1.AllowedPaths
	2, // 5: authd.v1.Config.TokensEntry.value:type_name -> authd.v1.Tokens
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_authd_v1_authd_proto_init() }
func file_authd_v1_authd_proto_init() {
	if File_authd_v1_authd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_authd_v1_authd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_authd_v1_authd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllowedPaths); i {
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
		file_authd_v1_authd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tokens); i {
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
		file_authd_v1_authd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Token); i {
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
			RawDescriptor: file_authd_v1_authd_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_authd_v1_authd_proto_goTypes,
		DependencyIndexes: file_authd_v1_authd_proto_depIdxs,
		MessageInfos:      file_authd_v1_authd_proto_msgTypes,
	}.Build()
	File_authd_v1_authd_proto = out.File
	file_authd_v1_authd_proto_rawDesc = nil
	file_authd_v1_authd_proto_goTypes = nil
	file_authd_v1_authd_proto_depIdxs = nil
}
