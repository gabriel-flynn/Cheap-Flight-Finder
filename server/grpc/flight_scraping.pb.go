// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.6.1
// source: flight_scraping.proto

package grpc

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flight_scraping_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_flight_scraping_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_flight_scraping_proto_rawDescGZIP(), []int{0}
}

type SouthwestHeadersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Headers map[string]string `protobuf:"bytes,1,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *SouthwestHeadersResponse) Reset() {
	*x = SouthwestHeadersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_flight_scraping_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SouthwestHeadersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SouthwestHeadersResponse) ProtoMessage() {}

func (x *SouthwestHeadersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_flight_scraping_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SouthwestHeadersResponse.ProtoReflect.Descriptor instead.
func (*SouthwestHeadersResponse) Descriptor() ([]byte, []int) {
	return file_flight_scraping_proto_rawDescGZIP(), []int{1}
}

func (x *SouthwestHeadersResponse) GetHeaders() map[string]string {
	if x != nil {
		return x.Headers
	}
	return nil
}

var File_flight_scraping_proto protoreflect.FileDescriptor

var file_flight_scraping_proto_rawDesc = []byte{
	0x0a, 0x15, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x73, 0x63, 0x72, 0x61, 0x70, 0x69, 0x6e,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73,
	0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0xa0, 0x01, 0x0a, 0x18, 0x53, 0x6f,
	0x75, 0x74, 0x68, 0x77, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x73, 0x2e, 0x53, 0x6f, 0x75, 0x74, 0x68, 0x77, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73,
	0x1a, 0x3a, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x5b, 0x0a, 0x0d,
	0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x53, 0x63, 0x72, 0x61, 0x70, 0x65, 0x72, 0x12, 0x4a, 0x0a,
	0x13, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x75, 0x74, 0x68, 0x77, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x12, 0x0e, 0x2e, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x66, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x53,
	0x6f, 0x75, 0x74, 0x68, 0x77, 0x65, 0x73, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x61, 0x62, 0x72, 0x69, 0x65, 0x6c, 0x2d,
	0x66, 0x6c, 0x79, 0x6e, 0x6e, 0x2f, 0x43, 0x68, 0x65, 0x61, 0x70, 0x2d, 0x46, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x2d, 0x46, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_flight_scraping_proto_rawDescOnce sync.Once
	file_flight_scraping_proto_rawDescData = file_flight_scraping_proto_rawDesc
)

func file_flight_scraping_proto_rawDescGZIP() []byte {
	file_flight_scraping_proto_rawDescOnce.Do(func() {
		file_flight_scraping_proto_rawDescData = protoimpl.X.CompressGZIP(file_flight_scraping_proto_rawDescData)
	})
	return file_flight_scraping_proto_rawDescData
}

var file_flight_scraping_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_flight_scraping_proto_goTypes = []interface{}{
	(*Empty)(nil),                    // 0: flights.Empty
	(*SouthwestHeadersResponse)(nil), // 1: flights.SouthwestHeadersResponse
	nil,                              // 2: flights.SouthwestHeadersResponse.HeadersEntry
}
var file_flight_scraping_proto_depIdxs = []int32{
	2, // 0: flights.SouthwestHeadersResponse.headers:type_name -> flights.SouthwestHeadersResponse.HeadersEntry
	0, // 1: flights.FlightScraper.GetSouthwestHeaders:input_type -> flights.Empty
	1, // 2: flights.FlightScraper.GetSouthwestHeaders:output_type -> flights.SouthwestHeadersResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_flight_scraping_proto_init() }
func file_flight_scraping_proto_init() {
	if File_flight_scraping_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_flight_scraping_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_flight_scraping_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SouthwestHeadersResponse); i {
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
			RawDescriptor: file_flight_scraping_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_flight_scraping_proto_goTypes,
		DependencyIndexes: file_flight_scraping_proto_depIdxs,
		MessageInfos:      file_flight_scraping_proto_msgTypes,
	}.Build()
	File_flight_scraping_proto = out.File
	file_flight_scraping_proto_rawDesc = nil
	file_flight_scraping_proto_goTypes = nil
	file_flight_scraping_proto_depIdxs = nil
}
