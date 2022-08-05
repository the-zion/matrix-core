// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: comment/service/v1/comment_error.proto

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

type CommentErrorReason int32

const (
	CommentErrorReason_UNKNOWN_ERROR            CommentErrorReason = 0
	CommentErrorReason_GET_COMMENT_DRAFT_FAILED CommentErrorReason = 1
	CommentErrorReason_GET_COMMENT_LIST_FAILED  CommentErrorReason = 2
	CommentErrorReason_CREATE_DRAFT_FAILED      CommentErrorReason = 3
	CommentErrorReason_CREATE_COMMENT_FAILED    CommentErrorReason = 4
	CommentErrorReason_SET_RECORD_FAILED        CommentErrorReason = 5
	CommentErrorReason_RECORD_NOT_FOUND         CommentErrorReason = 6
)

// Enum value maps for CommentErrorReason.
var (
	CommentErrorReason_name = map[int32]string{
		0: "UNKNOWN_ERROR",
		1: "GET_COMMENT_DRAFT_FAILED",
		2: "GET_COMMENT_LIST_FAILED",
		3: "CREATE_DRAFT_FAILED",
		4: "CREATE_COMMENT_FAILED",
		5: "SET_RECORD_FAILED",
		6: "RECORD_NOT_FOUND",
	}
	CommentErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR":            0,
		"GET_COMMENT_DRAFT_FAILED": 1,
		"GET_COMMENT_LIST_FAILED":  2,
		"CREATE_DRAFT_FAILED":      3,
		"CREATE_COMMENT_FAILED":    4,
		"SET_RECORD_FAILED":        5,
		"RECORD_NOT_FOUND":         6,
	}
)

func (x CommentErrorReason) Enum() *CommentErrorReason {
	p := new(CommentErrorReason)
	*p = x
	return p
}

func (x CommentErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CommentErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_comment_service_v1_comment_error_proto_enumTypes[0].Descriptor()
}

func (CommentErrorReason) Type() protoreflect.EnumType {
	return &file_comment_service_v1_comment_error_proto_enumTypes[0]
}

func (x CommentErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CommentErrorReason.Descriptor instead.
func (CommentErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_comment_service_v1_comment_error_proto_rawDescGZIP(), []int{0}
}

var File_comment_service_v1_comment_error_proto protoreflect.FileDescriptor

var file_comment_service_v1_comment_error_proto_rawDesc = []byte{
	0x0a, 0x26, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xc9, 0x01, 0x0a, 0x12, 0x43, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x01, 0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x45, 0x4e, 0x54,
	0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x17,
	0x0a, 0x13, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x52, 0x45, 0x41, 0x54,
	0x45, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x45, 0x54, 0x5f, 0x52, 0x45, 0x43, 0x4f, 0x52, 0x44,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x12, 0x14, 0x0a, 0x10, 0x52, 0x45, 0x43,
	0x4f, 0x52, 0x44, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x06, 0x1a,
	0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x1e, 0x5a, 0x1c, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x62, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comment_service_v1_comment_error_proto_rawDescOnce sync.Once
	file_comment_service_v1_comment_error_proto_rawDescData = file_comment_service_v1_comment_error_proto_rawDesc
)

func file_comment_service_v1_comment_error_proto_rawDescGZIP() []byte {
	file_comment_service_v1_comment_error_proto_rawDescOnce.Do(func() {
		file_comment_service_v1_comment_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_comment_service_v1_comment_error_proto_rawDescData)
	})
	return file_comment_service_v1_comment_error_proto_rawDescData
}

var file_comment_service_v1_comment_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_comment_service_v1_comment_error_proto_goTypes = []interface{}{
	(CommentErrorReason)(0), // 0: comment.v1.CommentErrorReason
}
var file_comment_service_v1_comment_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_comment_service_v1_comment_error_proto_init() }
func file_comment_service_v1_comment_error_proto_init() {
	if File_comment_service_v1_comment_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_comment_service_v1_comment_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_comment_service_v1_comment_error_proto_goTypes,
		DependencyIndexes: file_comment_service_v1_comment_error_proto_depIdxs,
		EnumInfos:         file_comment_service_v1_comment_error_proto_enumTypes,
	}.Build()
	File_comment_service_v1_comment_error_proto = out.File
	file_comment_service_v1_comment_error_proto_rawDesc = nil
	file_comment_service_v1_comment_error_proto_goTypes = nil
	file_comment_service_v1_comment_error_proto_depIdxs = nil
}