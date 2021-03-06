// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: httpmid/v1/accesslog.proto

package httpmidv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AccessLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ts            *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=ts,proto3" json:"ts,omitempty"`
	TraceId       string                 `protobuf:"bytes,2,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	SpanId        string                 `protobuf:"bytes,3,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	HttpMethod    string                 `protobuf:"bytes,4,opt,name=http_method,json=httpMethod,proto3" json:"http_method,omitempty"`
	HttpUrl       string                 `protobuf:"bytes,5,opt,name=http_url,json=httpUrl,proto3" json:"http_url,omitempty"`
	HttpVersion   string                 `protobuf:"bytes,6,opt,name=http_version,json=httpVersion,proto3" json:"http_version,omitempty"`
	HttpHost      string                 `protobuf:"bytes,7,opt,name=http_host,json=httpHost,proto3" json:"http_host,omitempty"`
	HttpUseragent string                 `protobuf:"bytes,8,opt,name=http_useragent,json=httpUseragent,proto3" json:"http_useragent,omitempty"`
	HttpReferrer  string                 `protobuf:"bytes,9,opt,name=http_referrer,json=httpReferrer,proto3" json:"http_referrer,omitempty"`
	HandleTime    *durationpb.Duration   `protobuf:"bytes,10,opt,name=handle_time,json=handleTime,proto3" json:"handle_time,omitempty"`
	HttpStatus    int32                  `protobuf:"varint,11,opt,name=http_status,json=httpStatus,proto3" json:"http_status,omitempty"`
	BytesWritten  int64                  `protobuf:"varint,12,opt,name=bytes_written,json=bytesWritten,proto3" json:"bytes_written,omitempty"` // int64 bytes_read = 13;
}

func (x *AccessLog) Reset() {
	*x = AccessLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_httpmid_v1_accesslog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccessLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessLog) ProtoMessage() {}

func (x *AccessLog) ProtoReflect() protoreflect.Message {
	mi := &file_httpmid_v1_accesslog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessLog.ProtoReflect.Descriptor instead.
func (*AccessLog) Descriptor() ([]byte, []int) {
	return file_httpmid_v1_accesslog_proto_rawDescGZIP(), []int{0}
}

func (x *AccessLog) GetTs() *timestamppb.Timestamp {
	if x != nil {
		return x.Ts
	}
	return nil
}

func (x *AccessLog) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *AccessLog) GetSpanId() string {
	if x != nil {
		return x.SpanId
	}
	return ""
}

func (x *AccessLog) GetHttpMethod() string {
	if x != nil {
		return x.HttpMethod
	}
	return ""
}

func (x *AccessLog) GetHttpUrl() string {
	if x != nil {
		return x.HttpUrl
	}
	return ""
}

func (x *AccessLog) GetHttpVersion() string {
	if x != nil {
		return x.HttpVersion
	}
	return ""
}

func (x *AccessLog) GetHttpHost() string {
	if x != nil {
		return x.HttpHost
	}
	return ""
}

func (x *AccessLog) GetHttpUseragent() string {
	if x != nil {
		return x.HttpUseragent
	}
	return ""
}

func (x *AccessLog) GetHttpReferrer() string {
	if x != nil {
		return x.HttpReferrer
	}
	return ""
}

func (x *AccessLog) GetHandleTime() *durationpb.Duration {
	if x != nil {
		return x.HandleTime
	}
	return nil
}

func (x *AccessLog) GetHttpStatus() int32 {
	if x != nil {
		return x.HttpStatus
	}
	return 0
}

func (x *AccessLog) GetBytesWritten() int64 {
	if x != nil {
		return x.BytesWritten
	}
	return 0
}

var File_httpmid_v1_accesslog_proto protoreflect.FileDescriptor

var file_httpmid_v1_accesslog_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x68, 0x74, 0x74, 0x70, 0x6d, 0x69, 0x64, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x68, 0x74,
	0x74, 0x70, 0x6d, 0x69, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb5, 0x03, 0x0a, 0x09, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x12, 0x2a, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x02, 0x74, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x74, 0x74, 0x70, 0x5f,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x68, 0x74,
	0x74, 0x70, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x74, 0x74, 0x70,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x74, 0x74, 0x70,
	0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x68, 0x74, 0x74, 0x70, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x68,
	0x6f, 0x73, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x74, 0x74, 0x70, 0x48,
	0x6f, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x68, 0x74, 0x74, 0x70, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x61, 0x67, 0x65, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x68, 0x74, 0x74,
	0x70, 0x55, 0x73, 0x65, 0x72, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x68, 0x74,
	0x74, 0x70, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x52, 0x65, 0x66, 0x65, 0x72, 0x72, 0x65, 0x72, 0x12,
	0x3a, 0x0a, 0x0b, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0a, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x68,
	0x74, 0x74, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x68, 0x74, 0x74, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a, 0x0d,
	0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0c, 0x62, 0x79, 0x74, 0x65, 0x73, 0x57, 0x72, 0x69, 0x74, 0x74, 0x65,
	0x6e, 0x42, 0x9b, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x74, 0x74, 0x70, 0x6d, 0x69,
	0x64, 0x2e, 0x76, 0x31, 0x42, 0x0e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6c, 0x6f, 0x67, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x30, 0x67, 0x6f, 0x2e, 0x73, 0x65, 0x61, 0x6e, 0x6b,
	0x68, 0x6c, 0x69, 0x61, 0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x6f, 0x6e, 0x6f, 0x2f, 0x61,
	0x70, 0x69, 0x73, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x6d, 0x69, 0x64, 0x2f, 0x76, 0x31, 0x3b, 0x68,
	0x74, 0x74, 0x70, 0x6d, 0x69, 0x64, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x48, 0x58, 0x58, 0xaa, 0x02,
	0x0a, 0x48, 0x74, 0x74, 0x70, 0x6d, 0x69, 0x64, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a, 0x48, 0x74,
	0x74, 0x70, 0x6d, 0x69, 0x64, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16, 0x48, 0x74, 0x74, 0x70, 0x6d,
	0x69, 0x64, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0b, 0x48, 0x74, 0x74, 0x70, 0x6d, 0x69, 0x64, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_httpmid_v1_accesslog_proto_rawDescOnce sync.Once
	file_httpmid_v1_accesslog_proto_rawDescData = file_httpmid_v1_accesslog_proto_rawDesc
)

func file_httpmid_v1_accesslog_proto_rawDescGZIP() []byte {
	file_httpmid_v1_accesslog_proto_rawDescOnce.Do(func() {
		file_httpmid_v1_accesslog_proto_rawDescData = protoimpl.X.CompressGZIP(file_httpmid_v1_accesslog_proto_rawDescData)
	})
	return file_httpmid_v1_accesslog_proto_rawDescData
}

var file_httpmid_v1_accesslog_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_httpmid_v1_accesslog_proto_goTypes = []interface{}{
	(*AccessLog)(nil),             // 0: httpmid.v1.AccessLog
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),   // 2: google.protobuf.Duration
}
var file_httpmid_v1_accesslog_proto_depIdxs = []int32{
	1, // 0: httpmid.v1.AccessLog.ts:type_name -> google.protobuf.Timestamp
	2, // 1: httpmid.v1.AccessLog.handle_time:type_name -> google.protobuf.Duration
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_httpmid_v1_accesslog_proto_init() }
func file_httpmid_v1_accesslog_proto_init() {
	if File_httpmid_v1_accesslog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_httpmid_v1_accesslog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccessLog); i {
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
			RawDescriptor: file_httpmid_v1_accesslog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_httpmid_v1_accesslog_proto_goTypes,
		DependencyIndexes: file_httpmid_v1_accesslog_proto_depIdxs,
		MessageInfos:      file_httpmid_v1_accesslog_proto_msgTypes,
	}.Build()
	File_httpmid_v1_accesslog_proto = out.File
	file_httpmid_v1_accesslog_proto_rawDesc = nil
	file_httpmid_v1_accesslog_proto_goTypes = nil
	file_httpmid_v1_accesslog_proto_depIdxs = nil
}
