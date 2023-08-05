// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.2
// source: mini_daemon.proto

package daemon

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

type ListenAddressForAppRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppId string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
}

func (x *ListenAddressForAppRequest) Reset() {
	*x = ListenAddressForAppRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mini_daemon_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenAddressForAppRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenAddressForAppRequest) ProtoMessage() {}

func (x *ListenAddressForAppRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mini_daemon_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenAddressForAppRequest.ProtoReflect.Descriptor instead.
func (*ListenAddressForAppRequest) Descriptor() ([]byte, []int) {
	return file_mini_daemon_proto_rawDescGZIP(), []int{0}
}

func (x *ListenAddressForAppRequest) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

type ListenAddressForAppResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ListenAddress string `protobuf:"bytes,1,opt,name=listen_address,json=listenAddress,proto3" json:"listen_address,omitempty"` // the hostname and port the app is listening on
}

func (x *ListenAddressForAppResponse) Reset() {
	*x = ListenAddressForAppResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mini_daemon_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListenAddressForAppResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListenAddressForAppResponse) ProtoMessage() {}

func (x *ListenAddressForAppResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mini_daemon_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListenAddressForAppResponse.ProtoReflect.Descriptor instead.
func (*ListenAddressForAppResponse) Descriptor() ([]byte, []int) {
	return file_mini_daemon_proto_rawDescGZIP(), []int{1}
}

func (x *ListenAddressForAppResponse) GetListenAddress() string {
	if x != nil {
		return x.ListenAddress
	}
	return ""
}

var File_mini_daemon_proto protoreflect.FileDescriptor

var file_mini_daemon_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6d, 0x69, 0x6e, 0x69, 0x5f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x65, 0x6e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x64, 0x61, 0x65, 0x6d,
	0x6f, 0x6e, 0x22, 0x33, 0x0a, 0x1a, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x46, 0x6f, 0x72, 0x41, 0x70, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x15, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x22, 0x44, 0x0a, 0x1b, 0x4c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x46, 0x6f, 0x72, 0x41, 0x70, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x32, 0x76, 0x0a,
	0x06, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x12, 0x6c, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x46, 0x6f, 0x72, 0x41, 0x70, 0x70, 0x12, 0x29,
	0x2e, 0x65, 0x6e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x65, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x46, 0x6f, 0x72, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x65, 0x6e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x46, 0x6f, 0x72, 0x41, 0x70, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2f, 0x5a, 0x2d, 0x65, 0x6e, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x64, 0x65, 0x76, 0x2f, 0x73, 0x68, 0x65, 0x6c, 0x6c, 0x2f, 0x73, 0x68, 0x65, 0x6c, 0x6c, 0x72,
	0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mini_daemon_proto_rawDescOnce sync.Once
	file_mini_daemon_proto_rawDescData = file_mini_daemon_proto_rawDesc
)

func file_mini_daemon_proto_rawDescGZIP() []byte {
	file_mini_daemon_proto_rawDescOnce.Do(func() {
		file_mini_daemon_proto_rawDescData = protoimpl.X.CompressGZIP(file_mini_daemon_proto_rawDescData)
	})
	return file_mini_daemon_proto_rawDescData
}

var file_mini_daemon_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_mini_daemon_proto_goTypes = []interface{}{
	(*ListenAddressForAppRequest)(nil),  // 0: encore.daemon.ListenAddressForAppRequest
	(*ListenAddressForAppResponse)(nil), // 1: encore.daemon.ListenAddressForAppResponse
}
var file_mini_daemon_proto_depIdxs = []int32{
	0, // 0: encore.daemon.Daemon.ListenAddressForApp:input_type -> encore.daemon.ListenAddressForAppRequest
	1, // 1: encore.daemon.Daemon.ListenAddressForApp:output_type -> encore.daemon.ListenAddressForAppResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mini_daemon_proto_init() }
func file_mini_daemon_proto_init() {
	if File_mini_daemon_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mini_daemon_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenAddressForAppRequest); i {
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
		file_mini_daemon_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListenAddressForAppResponse); i {
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
			RawDescriptor: file_mini_daemon_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mini_daemon_proto_goTypes,
		DependencyIndexes: file_mini_daemon_proto_depIdxs,
		MessageInfos:      file_mini_daemon_proto_msgTypes,
	}.Build()
	File_mini_daemon_proto = out.File
	file_mini_daemon_proto_rawDesc = nil
	file_mini_daemon_proto_goTypes = nil
	file_mini_daemon_proto_depIdxs = nil
}
