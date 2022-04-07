// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: user/service/v1/user_error.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type UserErrorReason int32

const (
	UserErrorReason_UNKNOWN_ERROR     UserErrorReason = 0
	UserErrorReason_LOGIN_FAILED      UserErrorReason = 1
	UserErrorReason_USERNAME_CONFLICT UserErrorReason = 2
	UserErrorReason_PHONE_CONFLICT    UserErrorReason = 3
	UserErrorReason_EMAIL_CONFLICT    UserErrorReason = 4
	UserErrorReason_REGISTER_FAILED   UserErrorReason = 5
)

// Enum value maps for UserErrorReason.
var (
	UserErrorReason_name = map[int32]string{
		0: "UNKNOWN_ERROR",
		1: "LOGIN_FAILED",
		2: "USERNAME_CONFLICT",
		3: "PHONE_CONFLICT",
		4: "EMAIL_CONFLICT",
		5: "REGISTER_FAILED",
	}
	UserErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR":     0,
		"LOGIN_FAILED":      1,
		"USERNAME_CONFLICT": 2,
		"PHONE_CONFLICT":    3,
		"EMAIL_CONFLICT":    4,
		"REGISTER_FAILED":   5,
	}
)

func (x UserErrorReason) Enum() *UserErrorReason {
	p := new(UserErrorReason)
	*p = x
	return p
}

func (x UserErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_user_service_v1_user_error_proto_enumTypes[0].Descriptor()
}

func (UserErrorReason) Type() protoreflect.EnumType {
	return &file_user_service_v1_user_error_proto_enumTypes[0]
}

func (x UserErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserErrorReason.Descriptor instead.
func (UserErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_user_service_v1_user_error_proto_rawDescGZIP(), []int{0}
}

var File_user_service_v1_user_error_proto protoreflect.FileDescriptor

var file_user_service_v1_user_error_proto_rawDesc = []byte{
	0x0a, 0x20, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76,
	0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2a, 0xae, 0x01, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x0c, 0x4c, 0x4f, 0x47, 0x49, 0x4e,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x01, 0x1a, 0x04, 0xa8, 0x45, 0x91, 0x03, 0x12,
	0x1b, 0x0a, 0x11, 0x55, 0x53, 0x45, 0x52, 0x4e, 0x41, 0x4d, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46,
	0x4c, 0x49, 0x43, 0x54, 0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x99, 0x03, 0x12, 0x18, 0x0a, 0x0e,
	0x50, 0x48, 0x4f, 0x4e, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x4c, 0x49, 0x43, 0x54, 0x10, 0x03,
	0x1a, 0x04, 0xa8, 0x45, 0x99, 0x03, 0x12, 0x18, 0x0a, 0x0e, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0x5f,
	0x43, 0x4f, 0x4e, 0x46, 0x4c, 0x49, 0x43, 0x54, 0x10, 0x04, 0x1a, 0x04, 0xa8, 0x45, 0x99, 0x03,
	0x12, 0x19, 0x0a, 0x0f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x05, 0x1a, 0x04, 0xa8, 0x45, 0x9d, 0x04, 0x1a, 0x04, 0xa0, 0x45, 0xf4,
	0x03, 0x42, 0x18, 0x5a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_user_service_v1_user_error_proto_rawDescOnce sync.Once
	file_user_service_v1_user_error_proto_rawDescData = file_user_service_v1_user_error_proto_rawDesc
)

func file_user_service_v1_user_error_proto_rawDescGZIP() []byte {
	file_user_service_v1_user_error_proto_rawDescOnce.Do(func() {
		file_user_service_v1_user_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_service_v1_user_error_proto_rawDescData)
	})
	return file_user_service_v1_user_error_proto_rawDescData
}

var file_user_service_v1_user_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_user_service_v1_user_error_proto_goTypes = []interface{}{
	(UserErrorReason)(0), // 0: user.v1.UserErrorReason
}
var file_user_service_v1_user_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_user_service_v1_user_error_proto_init() }
func file_user_service_v1_user_error_proto_init() {
	if File_user_service_v1_user_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_service_v1_user_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_service_v1_user_error_proto_goTypes,
		DependencyIndexes: file_user_service_v1_user_error_proto_depIdxs,
		EnumInfos:         file_user_service_v1_user_error_proto_enumTypes,
	}.Build()
	File_user_service_v1_user_error_proto = out.File
	file_user_service_v1_user_error_proto_rawDesc = nil
	file_user_service_v1_user_error_proto_goTypes = nil
	file_user_service_v1_user_error_proto_depIdxs = nil
}
