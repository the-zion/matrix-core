// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: achievement/service/v1/achievement_error.proto

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

type AchievementErrorReason int32

const (
	AchievementErrorReason_UNKNOWN_ERROR                     AchievementErrorReason = 0
	AchievementErrorReason_SET_ACHIEVEMENT_AGREE_FAILED      AchievementErrorReason = 1
	AchievementErrorReason_SET_ACHIEVEMENT_VIEW_FAILED       AchievementErrorReason = 2
	AchievementErrorReason_SET_ACHIEVEMENT_COLLECT_FAILED    AchievementErrorReason = 3
	AchievementErrorReason_SET_ACHIEVEMENT_FOLLOW_FAILED     AchievementErrorReason = 4
	AchievementErrorReason_SET_MEDAL_FAILED                  AchievementErrorReason = 5
	AchievementErrorReason_ACCESS_MEDAL_FAILED               AchievementErrorReason = 6
	AchievementErrorReason_CANCEL_ACHIEVEMENT_AGREE_FAILED   AchievementErrorReason = 7
	AchievementErrorReason_CANCEL_ACHIEVEMENT_COLLECT_FAILED AchievementErrorReason = 8
	AchievementErrorReason_CANCEL_ACHIEVEMENT_FOLLOW_FAILED  AchievementErrorReason = 9
	AchievementErrorReason_CANCEL_MEDAL_SET_FAILED           AchievementErrorReason = 10
	AchievementErrorReason_GET_ACHIEVEMENT_LIST_FAILED       AchievementErrorReason = 11
	AchievementErrorReason_GET_ACHIEVEMENT_FAILED            AchievementErrorReason = 12
	AchievementErrorReason_GET_MEDAL_FAILED                  AchievementErrorReason = 13
	AchievementErrorReason_GET_ACTIVE_FAILED                 AchievementErrorReason = 14
	AchievementErrorReason_ADD_ACHIEVEMENT_SCORE_FAILED      AchievementErrorReason = 15
	AchievementErrorReason_REDUCE_ACHIEVEMENT_SCORE_FAILED   AchievementErrorReason = 16
)

// Enum value maps for AchievementErrorReason.
var (
	AchievementErrorReason_name = map[int32]string{
		0:  "UNKNOWN_ERROR",
		1:  "SET_ACHIEVEMENT_AGREE_FAILED",
		2:  "SET_ACHIEVEMENT_VIEW_FAILED",
		3:  "SET_ACHIEVEMENT_COLLECT_FAILED",
		4:  "SET_ACHIEVEMENT_FOLLOW_FAILED",
		5:  "SET_MEDAL_FAILED",
		6:  "ACCESS_MEDAL_FAILED",
		7:  "CANCEL_ACHIEVEMENT_AGREE_FAILED",
		8:  "CANCEL_ACHIEVEMENT_COLLECT_FAILED",
		9:  "CANCEL_ACHIEVEMENT_FOLLOW_FAILED",
		10: "CANCEL_MEDAL_SET_FAILED",
		11: "GET_ACHIEVEMENT_LIST_FAILED",
		12: "GET_ACHIEVEMENT_FAILED",
		13: "GET_MEDAL_FAILED",
		14: "GET_ACTIVE_FAILED",
		15: "ADD_ACHIEVEMENT_SCORE_FAILED",
		16: "REDUCE_ACHIEVEMENT_SCORE_FAILED",
	}
	AchievementErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR":                     0,
		"SET_ACHIEVEMENT_AGREE_FAILED":      1,
		"SET_ACHIEVEMENT_VIEW_FAILED":       2,
		"SET_ACHIEVEMENT_COLLECT_FAILED":    3,
		"SET_ACHIEVEMENT_FOLLOW_FAILED":     4,
		"SET_MEDAL_FAILED":                  5,
		"ACCESS_MEDAL_FAILED":               6,
		"CANCEL_ACHIEVEMENT_AGREE_FAILED":   7,
		"CANCEL_ACHIEVEMENT_COLLECT_FAILED": 8,
		"CANCEL_ACHIEVEMENT_FOLLOW_FAILED":  9,
		"CANCEL_MEDAL_SET_FAILED":           10,
		"GET_ACHIEVEMENT_LIST_FAILED":       11,
		"GET_ACHIEVEMENT_FAILED":            12,
		"GET_MEDAL_FAILED":                  13,
		"GET_ACTIVE_FAILED":                 14,
		"ADD_ACHIEVEMENT_SCORE_FAILED":      15,
		"REDUCE_ACHIEVEMENT_SCORE_FAILED":   16,
	}
)

func (x AchievementErrorReason) Enum() *AchievementErrorReason {
	p := new(AchievementErrorReason)
	*p = x
	return p
}

func (x AchievementErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AchievementErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_achievement_service_v1_achievement_error_proto_enumTypes[0].Descriptor()
}

func (AchievementErrorReason) Type() protoreflect.EnumType {
	return &file_achievement_service_v1_achievement_error_proto_enumTypes[0]
}

func (x AchievementErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AchievementErrorReason.Descriptor instead.
func (AchievementErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_achievement_service_v1_achievement_error_proto_rawDescGZIP(), []int{0}
}

var File_achievement_service_v1_achievement_error_proto protoreflect.FileDescriptor

var file_achievement_service_v1_achievement_error_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0e, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31,
	0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xaa, 0x04, 0x0a, 0x16, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0x00, 0x12, 0x20, 0x0a, 0x1c, 0x53, 0x45, 0x54, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45,
	0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x41, 0x47, 0x52, 0x45, 0x45, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x45, 0x54, 0x5f, 0x41, 0x43, 0x48,
	0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x22, 0x0a, 0x1e, 0x53, 0x45, 0x54, 0x5f, 0x41, 0x43,
	0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43,
	0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x21, 0x0a, 0x1d, 0x53, 0x45,
	0x54, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x46, 0x4f,
	0x4c, 0x4c, 0x4f, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x12, 0x14, 0x0a,
	0x10, 0x53, 0x45, 0x54, 0x5f, 0x4d, 0x45, 0x44, 0x41, 0x4c, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x05, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x4d, 0x45,
	0x44, 0x41, 0x4c, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x06, 0x12, 0x23, 0x0a, 0x1f,
	0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45,
	0x4e, 0x54, 0x5f, 0x41, 0x47, 0x52, 0x45, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x07, 0x12, 0x25, 0x0a, 0x21, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x41, 0x43, 0x48, 0x49,
	0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x08, 0x12, 0x24, 0x0a, 0x20, 0x43, 0x41, 0x4e, 0x43,
	0x45, 0x4c, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x46,
	0x4f, 0x4c, 0x4c, 0x4f, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x09, 0x12, 0x1b,
	0x0a, 0x17, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x4d, 0x45, 0x44, 0x41, 0x4c, 0x5f, 0x53,
	0x45, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0a, 0x12, 0x1f, 0x0a, 0x1b, 0x47,
	0x45, 0x54, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4c,
	0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0b, 0x12, 0x1a, 0x0a, 0x16,
	0x47, 0x45, 0x54, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0c, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x45, 0x54, 0x5f,
	0x4d, 0x45, 0x44, 0x41, 0x4c, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0d, 0x12, 0x15,
	0x0a, 0x11, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x0e, 0x12, 0x20, 0x0a, 0x1c, 0x41, 0x44, 0x44, 0x5f, 0x41, 0x43, 0x48,
	0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x43, 0x4f, 0x52, 0x45, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0f, 0x12, 0x23, 0x0a, 0x1f, 0x52, 0x45, 0x44, 0x55, 0x43,
	0x45, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x43,
	0x4f, 0x52, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x10, 0x1a, 0x04, 0xa0, 0x45,
	0xf4, 0x03, 0x42, 0x22, 0x5a, 0x20, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x70, 0x62, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_achievement_service_v1_achievement_error_proto_rawDescOnce sync.Once
	file_achievement_service_v1_achievement_error_proto_rawDescData = file_achievement_service_v1_achievement_error_proto_rawDesc
)

func file_achievement_service_v1_achievement_error_proto_rawDescGZIP() []byte {
	file_achievement_service_v1_achievement_error_proto_rawDescOnce.Do(func() {
		file_achievement_service_v1_achievement_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_achievement_service_v1_achievement_error_proto_rawDescData)
	})
	return file_achievement_service_v1_achievement_error_proto_rawDescData
}

var file_achievement_service_v1_achievement_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_achievement_service_v1_achievement_error_proto_goTypes = []interface{}{
	(AchievementErrorReason)(0), // 0: achievement.v1.AchievementErrorReason
}
var file_achievement_service_v1_achievement_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_achievement_service_v1_achievement_error_proto_init() }
func file_achievement_service_v1_achievement_error_proto_init() {
	if File_achievement_service_v1_achievement_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_achievement_service_v1_achievement_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_achievement_service_v1_achievement_error_proto_goTypes,
		DependencyIndexes: file_achievement_service_v1_achievement_error_proto_depIdxs,
		EnumInfos:         file_achievement_service_v1_achievement_error_proto_enumTypes,
	}.Build()
	File_achievement_service_v1_achievement_error_proto = out.File
	file_achievement_service_v1_achievement_error_proto_rawDesc = nil
	file_achievement_service_v1_achievement_error_proto_goTypes = nil
	file_achievement_service_v1_achievement_error_proto_depIdxs = nil
}
