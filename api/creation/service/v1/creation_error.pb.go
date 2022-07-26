// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: creation/service/v1/creation_error.proto

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

type CreationErrorReason int32

const (
	CreationErrorReason_UNKNOWN_ERROR                    CreationErrorReason = 0
	CreationErrorReason_RECORD_NOT_FOUND                 CreationErrorReason = 1
	CreationErrorReason_GET_LEADER_BOARD_FAILED          CreationErrorReason = 2
	CreationErrorReason_GET_COLLECT_ARTICLE_FAILED       CreationErrorReason = 3
	CreationErrorReason_GET_COLLECTION_FAILED            CreationErrorReason = 4
	CreationErrorReason_GET_COLLECTIONS_FAILED           CreationErrorReason = 5
	CreationErrorReason_CREATE_COLLECTIONS_FAILED        CreationErrorReason = 6
	CreationErrorReason_EDIT_COLLECTIONS_FAILED          CreationErrorReason = 7
	CreationErrorReason_DELETE_COLLECTIONS_FAILED        CreationErrorReason = 8
	CreationErrorReason_GET_ARTICLE_FAILED               CreationErrorReason = 9
	CreationErrorReason_GET_ARTICLE_DRAFT_FAILED         CreationErrorReason = 10
	CreationErrorReason_GET_ARTICLE_LIST_FAILED          CreationErrorReason = 11
	CreationErrorReason_CREATE_ARTICLE_FAILED            CreationErrorReason = 12
	CreationErrorReason_EDIT_ARTICLE_FAILED              CreationErrorReason = 13
	CreationErrorReason_DELETE_ARTICLE_FAILED            CreationErrorReason = 14
	CreationErrorReason_GET_TALK_FAILED                  CreationErrorReason = 15
	CreationErrorReason_GET_TALK_DRAFT_FAILED            CreationErrorReason = 16
	CreationErrorReason_GET_TALK_LIST_FAILED             CreationErrorReason = 17
	CreationErrorReason_CREATE_TALK_FAILED               CreationErrorReason = 18
	CreationErrorReason_EDIT_TALK_FAILED                 CreationErrorReason = 19
	CreationErrorReason_DELETE_TALK_FAILED               CreationErrorReason = 20
	CreationErrorReason_GET_DRAFT_LIST_FAILED            CreationErrorReason = 21
	CreationErrorReason_CREATE_DRAFT_FAILED              CreationErrorReason = 22
	CreationErrorReason_DRAFT_MARK_FAILED                CreationErrorReason = 23
	CreationErrorReason_GET_COLUMN_FAILED                CreationErrorReason = 24
	CreationErrorReason_GET_COLUMN_DRAFT_FAILED          CreationErrorReason = 25
	CreationErrorReason_CREATE_COLUMN_FAILED             CreationErrorReason = 26
	CreationErrorReason_EDIT_COLUMN_FAILED               CreationErrorReason = 27
	CreationErrorReason_DELETE_COLUMN_FAILED             CreationErrorReason = 28
	CreationErrorReason_GET_COLUMN_LIST_FAILED           CreationErrorReason = 29
	CreationErrorReason_GET_SUBSCRIBE_COLUMN_LIST_FAILED CreationErrorReason = 30
	CreationErrorReason_GET_COLUMN_SUBSCRIBES_FAILED     CreationErrorReason = 31
	CreationErrorReason_ADD_COLUMN_INCLUDES_FAILED       CreationErrorReason = 32
	CreationErrorReason_DELETE_COLUMN_INCLUDES_FAILED    CreationErrorReason = 33
	CreationErrorReason_SUBSCRIBE_COLUMN_FAILED          CreationErrorReason = 34
	CreationErrorReason_SUBSCRIBE_COLUMN_JUDGE_FAILED    CreationErrorReason = 35
	CreationErrorReason_CANCEL_SUBSCRIBE_COLUMN_FAILED   CreationErrorReason = 36
	CreationErrorReason_GET_NEWS_FAILED                  CreationErrorReason = 37
	CreationErrorReason_SET_AGREE_FAILED                 CreationErrorReason = 38
	CreationErrorReason_SET_VIEW_FAILED                  CreationErrorReason = 39
	CreationErrorReason_SET_COLLECT_FAILED               CreationErrorReason = 40
	CreationErrorReason_CANCEL_AGREE_FAILED              CreationErrorReason = 41
	CreationErrorReason_CANCEL_VIEW_FAILED               CreationErrorReason = 42
	CreationErrorReason_CANCEL_COLLECT_FAILED            CreationErrorReason = 43
	CreationErrorReason_GET_STATISTIC_FAILED             CreationErrorReason = 44
	CreationErrorReason_GET_STATISTIC_JUDGE_FAILED       CreationErrorReason = 45
	CreationErrorReason_GET_COUNT_FAILED                 CreationErrorReason = 46
	CreationErrorReason_NOT_EMPTY                        CreationErrorReason = 47
)

// Enum value maps for CreationErrorReason.
var (
	CreationErrorReason_name = map[int32]string{
		0:  "UNKNOWN_ERROR",
		1:  "RECORD_NOT_FOUND",
		2:  "GET_LEADER_BOARD_FAILED",
		3:  "GET_COLLECT_ARTICLE_FAILED",
		4:  "GET_COLLECTION_FAILED",
		5:  "GET_COLLECTIONS_FAILED",
		6:  "CREATE_COLLECTIONS_FAILED",
		7:  "EDIT_COLLECTIONS_FAILED",
		8:  "DELETE_COLLECTIONS_FAILED",
		9:  "GET_ARTICLE_FAILED",
		10: "GET_ARTICLE_DRAFT_FAILED",
		11: "GET_ARTICLE_LIST_FAILED",
		12: "CREATE_ARTICLE_FAILED",
		13: "EDIT_ARTICLE_FAILED",
		14: "DELETE_ARTICLE_FAILED",
		15: "GET_TALK_FAILED",
		16: "GET_TALK_DRAFT_FAILED",
		17: "GET_TALK_LIST_FAILED",
		18: "CREATE_TALK_FAILED",
		19: "EDIT_TALK_FAILED",
		20: "DELETE_TALK_FAILED",
		21: "GET_DRAFT_LIST_FAILED",
		22: "CREATE_DRAFT_FAILED",
		23: "DRAFT_MARK_FAILED",
		24: "GET_COLUMN_FAILED",
		25: "GET_COLUMN_DRAFT_FAILED",
		26: "CREATE_COLUMN_FAILED",
		27: "EDIT_COLUMN_FAILED",
		28: "DELETE_COLUMN_FAILED",
		29: "GET_COLUMN_LIST_FAILED",
		30: "GET_SUBSCRIBE_COLUMN_LIST_FAILED",
		31: "GET_COLUMN_SUBSCRIBES_FAILED",
		32: "ADD_COLUMN_INCLUDES_FAILED",
		33: "DELETE_COLUMN_INCLUDES_FAILED",
		34: "SUBSCRIBE_COLUMN_FAILED",
		35: "SUBSCRIBE_COLUMN_JUDGE_FAILED",
		36: "CANCEL_SUBSCRIBE_COLUMN_FAILED",
		37: "GET_NEWS_FAILED",
		38: "SET_AGREE_FAILED",
		39: "SET_VIEW_FAILED",
		40: "SET_COLLECT_FAILED",
		41: "CANCEL_AGREE_FAILED",
		42: "CANCEL_VIEW_FAILED",
		43: "CANCEL_COLLECT_FAILED",
		44: "GET_STATISTIC_FAILED",
		45: "GET_STATISTIC_JUDGE_FAILED",
		46: "GET_COUNT_FAILED",
		47: "NOT_EMPTY",
	}
	CreationErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR":                    0,
		"RECORD_NOT_FOUND":                 1,
		"GET_LEADER_BOARD_FAILED":          2,
		"GET_COLLECT_ARTICLE_FAILED":       3,
		"GET_COLLECTION_FAILED":            4,
		"GET_COLLECTIONS_FAILED":           5,
		"CREATE_COLLECTIONS_FAILED":        6,
		"EDIT_COLLECTIONS_FAILED":          7,
		"DELETE_COLLECTIONS_FAILED":        8,
		"GET_ARTICLE_FAILED":               9,
		"GET_ARTICLE_DRAFT_FAILED":         10,
		"GET_ARTICLE_LIST_FAILED":          11,
		"CREATE_ARTICLE_FAILED":            12,
		"EDIT_ARTICLE_FAILED":              13,
		"DELETE_ARTICLE_FAILED":            14,
		"GET_TALK_FAILED":                  15,
		"GET_TALK_DRAFT_FAILED":            16,
		"GET_TALK_LIST_FAILED":             17,
		"CREATE_TALK_FAILED":               18,
		"EDIT_TALK_FAILED":                 19,
		"DELETE_TALK_FAILED":               20,
		"GET_DRAFT_LIST_FAILED":            21,
		"CREATE_DRAFT_FAILED":              22,
		"DRAFT_MARK_FAILED":                23,
		"GET_COLUMN_FAILED":                24,
		"GET_COLUMN_DRAFT_FAILED":          25,
		"CREATE_COLUMN_FAILED":             26,
		"EDIT_COLUMN_FAILED":               27,
		"DELETE_COLUMN_FAILED":             28,
		"GET_COLUMN_LIST_FAILED":           29,
		"GET_SUBSCRIBE_COLUMN_LIST_FAILED": 30,
		"GET_COLUMN_SUBSCRIBES_FAILED":     31,
		"ADD_COLUMN_INCLUDES_FAILED":       32,
		"DELETE_COLUMN_INCLUDES_FAILED":    33,
		"SUBSCRIBE_COLUMN_FAILED":          34,
		"SUBSCRIBE_COLUMN_JUDGE_FAILED":    35,
		"CANCEL_SUBSCRIBE_COLUMN_FAILED":   36,
		"GET_NEWS_FAILED":                  37,
		"SET_AGREE_FAILED":                 38,
		"SET_VIEW_FAILED":                  39,
		"SET_COLLECT_FAILED":               40,
		"CANCEL_AGREE_FAILED":              41,
		"CANCEL_VIEW_FAILED":               42,
		"CANCEL_COLLECT_FAILED":            43,
		"GET_STATISTIC_FAILED":             44,
		"GET_STATISTIC_JUDGE_FAILED":       45,
		"GET_COUNT_FAILED":                 46,
		"NOT_EMPTY":                        47,
	}
)

func (x CreationErrorReason) Enum() *CreationErrorReason {
	p := new(CreationErrorReason)
	*p = x
	return p
}

func (x CreationErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CreationErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_creation_service_v1_creation_error_proto_enumTypes[0].Descriptor()
}

func (CreationErrorReason) Type() protoreflect.EnumType {
	return &file_creation_service_v1_creation_error_proto_enumTypes[0]
}

func (x CreationErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CreationErrorReason.Descriptor instead.
func (CreationErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_creation_service_v1_creation_error_proto_rawDescGZIP(), []int{0}
}

var File_creation_service_v1_creation_error_proto protoreflect.FileDescriptor

var file_creation_service_v1_creation_error_proto_rawDesc = []byte{
	0x0a, 0x28, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x9e, 0x0a, 0x0a,
	0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x52, 0x45, 0x43, 0x4f, 0x52,
	0x44, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x1b, 0x0a,
	0x17, 0x47, 0x45, 0x54, 0x5f, 0x4c, 0x45, 0x41, 0x44, 0x45, 0x52, 0x5f, 0x42, 0x4f, 0x41, 0x52,
	0x44, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x47, 0x45,
	0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c,
	0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x45,
	0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x04, 0x12, 0x1a, 0x0a, 0x16, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c,
	0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x05, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x4c,
	0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x06,
	0x12, 0x1b, 0x0a, 0x17, 0x45, 0x44, 0x49, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54,
	0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x07, 0x12, 0x1d, 0x0a,
	0x19, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49,
	0x4f, 0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x08, 0x12, 0x16, 0x0a, 0x12,
	0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x09, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49,
	0x43, 0x4c, 0x45, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x0a, 0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c,
	0x45, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0b, 0x12,
	0x19, 0x0a, 0x15, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c,
	0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0c, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x44,
	0x49, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x0d, 0x12, 0x19, 0x0a, 0x15, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x41, 0x52,
	0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0e, 0x12, 0x13,
	0x0a, 0x0f, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x0f, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f,
	0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x10, 0x12, 0x18,
	0x0a, 0x14, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x11, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x52, 0x45, 0x41,
	0x54, 0x45, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x12,
	0x12, 0x14, 0x0a, 0x10, 0x45, 0x44, 0x49, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x13, 0x12, 0x16, 0x0a, 0x12, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x14, 0x12, 0x19,
	0x0a, 0x15, 0x47, 0x45, 0x54, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x4c, 0x49, 0x53, 0x54,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x15, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x52, 0x45,
	0x41, 0x54, 0x45, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x16, 0x12, 0x15, 0x0a, 0x11, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x4d, 0x41, 0x52, 0x4b,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x17, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x45, 0x54,
	0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x18,
	0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x44,
	0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x19, 0x12, 0x18, 0x0a,
	0x14, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1a, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x44, 0x49, 0x54, 0x5f,
	0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1b, 0x12,
	0x18, 0x0a, 0x14, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1c, 0x12, 0x1a, 0x0a, 0x16, 0x47, 0x45, 0x54,
	0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x1d, 0x12, 0x24, 0x0a, 0x20, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x55, 0x42,
	0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x4c, 0x49,
	0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1e, 0x12, 0x20, 0x0a, 0x1c, 0x47,
	0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52,
	0x49, 0x42, 0x45, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1f, 0x12, 0x1e, 0x0a,
	0x1a, 0x41, 0x44, 0x44, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x49, 0x4e, 0x43, 0x4c,
	0x55, 0x44, 0x45, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x20, 0x12, 0x21, 0x0a,
	0x1d, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x49,
	0x4e, 0x43, 0x4c, 0x55, 0x44, 0x45, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x21,
	0x12, 0x1b, 0x0a, 0x17, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f,
	0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x22, 0x12, 0x21, 0x0a,
	0x1d, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d,
	0x4e, 0x5f, 0x4a, 0x55, 0x44, 0x47, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x23,
	0x12, 0x22, 0x0a, 0x1e, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43,
	0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x24, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x45, 0x54, 0x5f, 0x4e, 0x45, 0x57, 0x53,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x25, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x45, 0x54,
	0x5f, 0x41, 0x47, 0x52, 0x45, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x26, 0x12,
	0x13, 0x0a, 0x0f, 0x53, 0x45, 0x54, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x27, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c,
	0x45, 0x43, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x28, 0x12, 0x17, 0x0a, 0x13,
	0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x41, 0x47, 0x52, 0x45, 0x45, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x29, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f,
	0x56, 0x49, 0x45, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2a, 0x12, 0x19, 0x0a,
	0x15, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2b, 0x12, 0x18, 0x0a, 0x14, 0x47, 0x45, 0x54, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x49, 0x53, 0x54, 0x49, 0x43, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x2c, 0x12, 0x1e, 0x0a, 0x1a, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x49, 0x53,
	0x54, 0x49, 0x43, 0x5f, 0x4a, 0x55, 0x44, 0x47, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x2d, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2e, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f,
	0x45, 0x4d, 0x50, 0x54, 0x59, 0x10, 0x2f, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x1f, 0x5a,
	0x1d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x62, 0x3b, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_creation_service_v1_creation_error_proto_rawDescOnce sync.Once
	file_creation_service_v1_creation_error_proto_rawDescData = file_creation_service_v1_creation_error_proto_rawDesc
)

func file_creation_service_v1_creation_error_proto_rawDescGZIP() []byte {
	file_creation_service_v1_creation_error_proto_rawDescOnce.Do(func() {
		file_creation_service_v1_creation_error_proto_rawDescData = protoimpl.X.CompressGZIP(file_creation_service_v1_creation_error_proto_rawDescData)
	})
	return file_creation_service_v1_creation_error_proto_rawDescData
}

var file_creation_service_v1_creation_error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_creation_service_v1_creation_error_proto_goTypes = []interface{}{
	(CreationErrorReason)(0), // 0: creation.v1.CreationErrorReason
}
var file_creation_service_v1_creation_error_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_creation_service_v1_creation_error_proto_init() }
func file_creation_service_v1_creation_error_proto_init() {
	if File_creation_service_v1_creation_error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_creation_service_v1_creation_error_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_creation_service_v1_creation_error_proto_goTypes,
		DependencyIndexes: file_creation_service_v1_creation_error_proto_depIdxs,
		EnumInfos:         file_creation_service_v1_creation_error_proto_enumTypes,
	}.Build()
	File_creation_service_v1_creation_error_proto = out.File
	file_creation_service_v1_creation_error_proto_rawDesc = nil
	file_creation_service_v1_creation_error_proto_goTypes = nil
	file_creation_service_v1_creation_error_proto_depIdxs = nil
}
