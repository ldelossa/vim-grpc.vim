// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: env/env_service.proto

package env

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

var File_env_env_service_proto protoreflect.FileDescriptor

var file_env_env_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x65, 0x6e, 0x76, 0x2f, 0x65, 0x6e, 0x76, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x65, 0x6e, 0x76, 0x1a, 0x0d, 0x65, 0x6e,
	0x76, 0x2f, 0x65, 0x6e, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x3a, 0x0a, 0x03, 0x45,
	0x6e, 0x76, 0x12, 0x33, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x76, 0x12, 0x12, 0x2e, 0x65,
	0x6e, 0x76, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x76, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x13, 0x2e, 0x65, 0x6e, 0x76, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x76, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x64, 0x65, 0x6c, 0x6f, 0x73, 0x61, 0x2f, 0x76, 0x69,
	0x6d, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x69, 0x6d, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x65, 0x6e, 0x76, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_env_env_service_proto_goTypes = []interface{}{
	(*GetEnvRequest)(nil),  // 0: env.GetEnvRequest
	(*GetEnvResponse)(nil), // 1: env.GetEnvResponse
}
var file_env_env_service_proto_depIdxs = []int32{
	0, // 0: env.Env.GetEnv:input_type -> env.GetEnvRequest
	1, // 1: env.Env.GetEnv:output_type -> env.GetEnvResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_env_env_service_proto_init() }
func file_env_env_service_proto_init() {
	if File_env_env_service_proto != nil {
		return
	}
	file_env_env_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_env_env_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_env_env_service_proto_goTypes,
		DependencyIndexes: file_env_env_service_proto_depIdxs,
	}.Build()
	File_env_env_service_proto = out.File
	file_env_env_service_proto_rawDesc = nil
	file_env_env_service_proto_goTypes = nil
	file_env_env_service_proto_depIdxs = nil
}
