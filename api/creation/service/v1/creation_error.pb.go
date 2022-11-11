// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
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
	CreationErrorReason_GET_COLLECTIONS_LIST_FAILED      CreationErrorReason = 5
	CreationErrorReason_GET_TIMELINE_LIST_FAILED         CreationErrorReason = 6
	CreationErrorReason_CREATE_COLLECTIONS_FAILED        CreationErrorReason = 7
	CreationErrorReason_CREATE_TIMELINE_FAILED           CreationErrorReason = 8
	CreationErrorReason_EDIT_COLLECTIONS_FAILED          CreationErrorReason = 9
	CreationErrorReason_DELETE_COLLECTIONS_FAILED        CreationErrorReason = 10
	CreationErrorReason_GET_ARTICLE_FAILED               CreationErrorReason = 11
	CreationErrorReason_GET_ARTICLE_DRAFT_FAILED         CreationErrorReason = 12
	CreationErrorReason_GET_ARTICLE_LIST_FAILED          CreationErrorReason = 13
	CreationErrorReason_GET_ARTICLE_SEARCH_FAILED        CreationErrorReason = 14
	CreationErrorReason_GET_ARTICLE_AGREE_FAILED         CreationErrorReason = 15
	CreationErrorReason_GET_ARTICLE_COLLECT_FAILED       CreationErrorReason = 16
	CreationErrorReason_CREATE_ARTICLE_FAILED            CreationErrorReason = 17
	CreationErrorReason_EDIT_ARTICLE_FAILED              CreationErrorReason = 18
	CreationErrorReason_DELETE_ARTICLE_FAILED            CreationErrorReason = 19
	CreationErrorReason_GET_TALK_FAILED                  CreationErrorReason = 20
	CreationErrorReason_GET_TALK_AGREE_FAILED            CreationErrorReason = 21
	CreationErrorReason_GET_TALK_COLLECT_FAILED          CreationErrorReason = 22
	CreationErrorReason_GET_TALK_DRAFT_FAILED            CreationErrorReason = 23
	CreationErrorReason_GET_TALK_LIST_FAILED             CreationErrorReason = 24
	CreationErrorReason_GET_TALK_SEARCH_FAILED           CreationErrorReason = 25
	CreationErrorReason_CREATE_TALK_FAILED               CreationErrorReason = 26
	CreationErrorReason_EDIT_TALK_FAILED                 CreationErrorReason = 27
	CreationErrorReason_DELETE_TALK_FAILED               CreationErrorReason = 28
	CreationErrorReason_GET_DRAFT_LIST_FAILED            CreationErrorReason = 29
	CreationErrorReason_CREATE_DRAFT_FAILED              CreationErrorReason = 30
	CreationErrorReason_DRAFT_MARK_FAILED                CreationErrorReason = 31
	CreationErrorReason_GET_COLUMN_FAILED                CreationErrorReason = 32
	CreationErrorReason_GET_COLUMN_AGREE_FAILED          CreationErrorReason = 33
	CreationErrorReason_GET_COLUMN_COLLECT_FAILED        CreationErrorReason = 34
	CreationErrorReason_GET_COLUMN_DRAFT_FAILED          CreationErrorReason = 35
	CreationErrorReason_GET_COLUMN_SEARCH_FAILED         CreationErrorReason = 36
	CreationErrorReason_CREATE_COLUMN_FAILED             CreationErrorReason = 37
	CreationErrorReason_EDIT_COLUMN_FAILED               CreationErrorReason = 38
	CreationErrorReason_DELETE_COLUMN_FAILED             CreationErrorReason = 39
	CreationErrorReason_GET_COLUMN_LIST_FAILED           CreationErrorReason = 40
	CreationErrorReason_GET_SUBSCRIBE_COLUMN_LIST_FAILED CreationErrorReason = 41
	CreationErrorReason_GET_COLUMN_SUBSCRIBES_FAILED     CreationErrorReason = 42
	CreationErrorReason_GET_SUBSCRIBE_COLUMN_FAILED      CreationErrorReason = 43
	CreationErrorReason_ADD_COLUMN_INCLUDES_FAILED       CreationErrorReason = 44
	CreationErrorReason_DELETE_COLUMN_INCLUDES_FAILED    CreationErrorReason = 45
	CreationErrorReason_SUBSCRIBE_COLUMN_FAILED          CreationErrorReason = 46
	CreationErrorReason_SUBSCRIBE_COLUMN_JUDGE_FAILED    CreationErrorReason = 47
	CreationErrorReason_CANCEL_SUBSCRIBE_COLUMN_FAILED   CreationErrorReason = 48
	CreationErrorReason_GET_NEWS_FAILED                  CreationErrorReason = 49
	CreationErrorReason_SET_AGREE_FAILED                 CreationErrorReason = 50
	CreationErrorReason_SET_VIEW_FAILED                  CreationErrorReason = 51
	CreationErrorReason_SET_COLLECT_FAILED               CreationErrorReason = 52
	CreationErrorReason_SET_IMAGE_IRREGULAR_FAILED       CreationErrorReason = 53
	CreationErrorReason_SET_CONTENT_IRREGULAR_FAILED     CreationErrorReason = 54
	CreationErrorReason_CANCEL_AGREE_FAILED              CreationErrorReason = 55
	CreationErrorReason_CANCEL_VIEW_FAILED               CreationErrorReason = 56
	CreationErrorReason_CANCEL_COLLECT_FAILED            CreationErrorReason = 57
	CreationErrorReason_GET_STATISTIC_FAILED             CreationErrorReason = 58
	CreationErrorReason_GET_STATISTIC_JUDGE_FAILED       CreationErrorReason = 59
	CreationErrorReason_GET_COUNT_FAILED                 CreationErrorReason = 60
	CreationErrorReason_GET_CREATION_USER_FAILED         CreationErrorReason = 61
	CreationErrorReason_GET_IMAGE_REVIEW_FAILED          CreationErrorReason = 62
	CreationErrorReason_GET_CONTENT_REVIEW_FAILED        CreationErrorReason = 63
	CreationErrorReason_SET_RECORD_FAILED                CreationErrorReason = 64
	CreationErrorReason_NOT_EMPTY                        CreationErrorReason = 65
	CreationErrorReason_ADD_COMMENT_FAILED               CreationErrorReason = 66
	CreationErrorReason_REDUCE_COMMENT_FAILED            CreationErrorReason = 67
)

// Enum value maps for CreationErrorReason.
var (
	CreationErrorReason_name = map[int32]string{
		0:  "UNKNOWN_ERROR",
		1:  "RECORD_NOT_FOUND",
		2:  "GET_LEADER_BOARD_FAILED",
		3:  "GET_COLLECT_ARTICLE_FAILED",
		4:  "GET_COLLECTION_FAILED",
		5:  "GET_COLLECTIONS_LIST_FAILED",
		6:  "GET_TIMELINE_LIST_FAILED",
		7:  "CREATE_COLLECTIONS_FAILED",
		8:  "CREATE_TIMELINE_FAILED",
		9:  "EDIT_COLLECTIONS_FAILED",
		10: "DELETE_COLLECTIONS_FAILED",
		11: "GET_ARTICLE_FAILED",
		12: "GET_ARTICLE_DRAFT_FAILED",
		13: "GET_ARTICLE_LIST_FAILED",
		14: "GET_ARTICLE_SEARCH_FAILED",
		15: "GET_ARTICLE_AGREE_FAILED",
		16: "GET_ARTICLE_COLLECT_FAILED",
		17: "CREATE_ARTICLE_FAILED",
		18: "EDIT_ARTICLE_FAILED",
		19: "DELETE_ARTICLE_FAILED",
		20: "GET_TALK_FAILED",
		21: "GET_TALK_AGREE_FAILED",
		22: "GET_TALK_COLLECT_FAILED",
		23: "GET_TALK_DRAFT_FAILED",
		24: "GET_TALK_LIST_FAILED",
		25: "GET_TALK_SEARCH_FAILED",
		26: "CREATE_TALK_FAILED",
		27: "EDIT_TALK_FAILED",
		28: "DELETE_TALK_FAILED",
		29: "GET_DRAFT_LIST_FAILED",
		30: "CREATE_DRAFT_FAILED",
		31: "DRAFT_MARK_FAILED",
		32: "GET_COLUMN_FAILED",
		33: "GET_COLUMN_AGREE_FAILED",
		34: "GET_COLUMN_COLLECT_FAILED",
		35: "GET_COLUMN_DRAFT_FAILED",
		36: "GET_COLUMN_SEARCH_FAILED",
		37: "CREATE_COLUMN_FAILED",
		38: "EDIT_COLUMN_FAILED",
		39: "DELETE_COLUMN_FAILED",
		40: "GET_COLUMN_LIST_FAILED",
		41: "GET_SUBSCRIBE_COLUMN_LIST_FAILED",
		42: "GET_COLUMN_SUBSCRIBES_FAILED",
		43: "GET_SUBSCRIBE_COLUMN_FAILED",
		44: "ADD_COLUMN_INCLUDES_FAILED",
		45: "DELETE_COLUMN_INCLUDES_FAILED",
		46: "SUBSCRIBE_COLUMN_FAILED",
		47: "SUBSCRIBE_COLUMN_JUDGE_FAILED",
		48: "CANCEL_SUBSCRIBE_COLUMN_FAILED",
		49: "GET_NEWS_FAILED",
		50: "SET_AGREE_FAILED",
		51: "SET_VIEW_FAILED",
		52: "SET_COLLECT_FAILED",
		53: "SET_IMAGE_IRREGULAR_FAILED",
		54: "SET_CONTENT_IRREGULAR_FAILED",
		55: "CANCEL_AGREE_FAILED",
		56: "CANCEL_VIEW_FAILED",
		57: "CANCEL_COLLECT_FAILED",
		58: "GET_STATISTIC_FAILED",
		59: "GET_STATISTIC_JUDGE_FAILED",
		60: "GET_COUNT_FAILED",
		61: "GET_CREATION_USER_FAILED",
		62: "GET_IMAGE_REVIEW_FAILED",
		63: "GET_CONTENT_REVIEW_FAILED",
		64: "SET_RECORD_FAILED",
		65: "NOT_EMPTY",
		66: "ADD_COMMENT_FAILED",
		67: "REDUCE_COMMENT_FAILED",
	}
	CreationErrorReason_value = map[string]int32{
		"UNKNOWN_ERROR":                    0,
		"RECORD_NOT_FOUND":                 1,
		"GET_LEADER_BOARD_FAILED":          2,
		"GET_COLLECT_ARTICLE_FAILED":       3,
		"GET_COLLECTION_FAILED":            4,
		"GET_COLLECTIONS_LIST_FAILED":      5,
		"GET_TIMELINE_LIST_FAILED":         6,
		"CREATE_COLLECTIONS_FAILED":        7,
		"CREATE_TIMELINE_FAILED":           8,
		"EDIT_COLLECTIONS_FAILED":          9,
		"DELETE_COLLECTIONS_FAILED":        10,
		"GET_ARTICLE_FAILED":               11,
		"GET_ARTICLE_DRAFT_FAILED":         12,
		"GET_ARTICLE_LIST_FAILED":          13,
		"GET_ARTICLE_SEARCH_FAILED":        14,
		"GET_ARTICLE_AGREE_FAILED":         15,
		"GET_ARTICLE_COLLECT_FAILED":       16,
		"CREATE_ARTICLE_FAILED":            17,
		"EDIT_ARTICLE_FAILED":              18,
		"DELETE_ARTICLE_FAILED":            19,
		"GET_TALK_FAILED":                  20,
		"GET_TALK_AGREE_FAILED":            21,
		"GET_TALK_COLLECT_FAILED":          22,
		"GET_TALK_DRAFT_FAILED":            23,
		"GET_TALK_LIST_FAILED":             24,
		"GET_TALK_SEARCH_FAILED":           25,
		"CREATE_TALK_FAILED":               26,
		"EDIT_TALK_FAILED":                 27,
		"DELETE_TALK_FAILED":               28,
		"GET_DRAFT_LIST_FAILED":            29,
		"CREATE_DRAFT_FAILED":              30,
		"DRAFT_MARK_FAILED":                31,
		"GET_COLUMN_FAILED":                32,
		"GET_COLUMN_AGREE_FAILED":          33,
		"GET_COLUMN_COLLECT_FAILED":        34,
		"GET_COLUMN_DRAFT_FAILED":          35,
		"GET_COLUMN_SEARCH_FAILED":         36,
		"CREATE_COLUMN_FAILED":             37,
		"EDIT_COLUMN_FAILED":               38,
		"DELETE_COLUMN_FAILED":             39,
		"GET_COLUMN_LIST_FAILED":           40,
		"GET_SUBSCRIBE_COLUMN_LIST_FAILED": 41,
		"GET_COLUMN_SUBSCRIBES_FAILED":     42,
		"GET_SUBSCRIBE_COLUMN_FAILED":      43,
		"ADD_COLUMN_INCLUDES_FAILED":       44,
		"DELETE_COLUMN_INCLUDES_FAILED":    45,
		"SUBSCRIBE_COLUMN_FAILED":          46,
		"SUBSCRIBE_COLUMN_JUDGE_FAILED":    47,
		"CANCEL_SUBSCRIBE_COLUMN_FAILED":   48,
		"GET_NEWS_FAILED":                  49,
		"SET_AGREE_FAILED":                 50,
		"SET_VIEW_FAILED":                  51,
		"SET_COLLECT_FAILED":               52,
		"SET_IMAGE_IRREGULAR_FAILED":       53,
		"SET_CONTENT_IRREGULAR_FAILED":     54,
		"CANCEL_AGREE_FAILED":              55,
		"CANCEL_VIEW_FAILED":               56,
		"CANCEL_COLLECT_FAILED":            57,
		"GET_STATISTIC_FAILED":             58,
		"GET_STATISTIC_JUDGE_FAILED":       59,
		"GET_COUNT_FAILED":                 60,
		"GET_CREATION_USER_FAILED":         61,
		"GET_IMAGE_REVIEW_FAILED":          62,
		"GET_CONTENT_REVIEW_FAILED":        63,
		"SET_RECORD_FAILED":                64,
		"NOT_EMPTY":                        65,
		"ADD_COMMENT_FAILED":               66,
		"REDUCE_COMMENT_FAILED":            67,
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
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0xef, 0x0e, 0x0a,
	0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65,
	0x61, 0x73, 0x6f, 0x6e, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x52, 0x45, 0x43, 0x4f, 0x52,
	0x44, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x1b, 0x0a,
	0x17, 0x47, 0x45, 0x54, 0x5f, 0x4c, 0x45, 0x41, 0x44, 0x45, 0x52, 0x5f, 0x42, 0x4f, 0x41, 0x52,
	0x44, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x47, 0x45,
	0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c,
	0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x45,
	0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x04, 0x12, 0x1f, 0x0a, 0x1b, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c,
	0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x49,
	0x4d, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x06, 0x12, 0x1d, 0x0a, 0x19, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x43,
	0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x07, 0x12, 0x1a, 0x0a, 0x16, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x54, 0x49,
	0x4d, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x08, 0x12,
	0x1b, 0x0a, 0x17, 0x45, 0x44, 0x49, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49,
	0x4f, 0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x09, 0x12, 0x1d, 0x0a, 0x19,
	0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0a, 0x12, 0x16, 0x0a, 0x12, 0x47,
	0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x0b, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43,
	0x4c, 0x45, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x0c, 0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45,
	0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0d, 0x12, 0x1d,
	0x0a, 0x19, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x53, 0x45,
	0x41, 0x52, 0x43, 0x48, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0e, 0x12, 0x1c, 0x0a,
	0x18, 0x47, 0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x41, 0x47, 0x52,
	0x45, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x0f, 0x12, 0x1e, 0x0a, 0x1a, 0x47,
	0x45, 0x54, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45,
	0x43, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x10, 0x12, 0x19, 0x0a, 0x15, 0x43,
	0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x11, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x44, 0x49, 0x54, 0x5f, 0x41,
	0x52, 0x54, 0x49, 0x43, 0x4c, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x12, 0x12,
	0x19, 0x0a, 0x15, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x41, 0x52, 0x54, 0x49, 0x43, 0x4c,
	0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x13, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x45,
	0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x14, 0x12,
	0x19, 0x0a, 0x15, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x41, 0x47, 0x52, 0x45,
	0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x15, 0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45,
	0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x16, 0x12, 0x19, 0x0a, 0x15, 0x47, 0x45, 0x54, 0x5f, 0x54,
	0x41, 0x4c, 0x4b, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x17, 0x12, 0x18, 0x0a, 0x14, 0x47, 0x45, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x4c,
	0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x18, 0x12, 0x1a, 0x0a, 0x16,
	0x47, 0x45, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x53, 0x45, 0x41, 0x52, 0x43, 0x48, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x19, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x52, 0x45, 0x41,
	0x54, 0x45, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1a,
	0x12, 0x14, 0x0a, 0x10, 0x45, 0x44, 0x49, 0x54, 0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x1b, 0x12, 0x16, 0x0a, 0x12, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x5f, 0x54, 0x41, 0x4c, 0x4b, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1c, 0x12, 0x19,
	0x0a, 0x15, 0x47, 0x45, 0x54, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x4c, 0x49, 0x53, 0x54,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1d, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x52, 0x45,
	0x41, 0x54, 0x45, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x1e, 0x12, 0x15, 0x0a, 0x11, 0x44, 0x52, 0x41, 0x46, 0x54, 0x5f, 0x4d, 0x41, 0x52, 0x4b,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x1f, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x45, 0x54,
	0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x20,
	0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x41,
	0x47, 0x52, 0x45, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x21, 0x12, 0x1d, 0x0a,
	0x19, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x43, 0x4f, 0x4c, 0x4c,
	0x45, 0x43, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x22, 0x12, 0x1b, 0x0a, 0x17,
	0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x44, 0x52, 0x41, 0x46, 0x54,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x23, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x45, 0x54,
	0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x53, 0x45, 0x41, 0x52, 0x43, 0x48, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x24, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x52, 0x45, 0x41, 0x54,
	0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x25, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x44, 0x49, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x26, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x45, 0x4c,
	0x45, 0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x27, 0x12, 0x1a, 0x0a, 0x16, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d,
	0x4e, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x28, 0x12,
	0x24, 0x0a, 0x20, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45,
	0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x29, 0x12, 0x20, 0x0a, 0x1c, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c,
	0x55, 0x4d, 0x4e, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x53, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2a, 0x12, 0x1f, 0x0a, 0x1b, 0x47, 0x45, 0x54, 0x5f, 0x53,
	0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2b, 0x12, 0x1e, 0x0a, 0x1a, 0x41, 0x44, 0x44, 0x5f,
	0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x49, 0x4e, 0x43, 0x4c, 0x55, 0x44, 0x45, 0x53, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2c, 0x12, 0x21, 0x0a, 0x1d, 0x44, 0x45, 0x4c, 0x45,
	0x54, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x49, 0x4e, 0x43, 0x4c, 0x55, 0x44,
	0x45, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2d, 0x12, 0x1b, 0x0a, 0x17, 0x53,
	0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2e, 0x12, 0x21, 0x0a, 0x1d, 0x53, 0x55, 0x42, 0x53,
	0x43, 0x52, 0x49, 0x42, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x4a, 0x55, 0x44,
	0x47, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x2f, 0x12, 0x22, 0x0a, 0x1e, 0x43,
	0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x42, 0x45, 0x5f,
	0x43, 0x4f, 0x4c, 0x55, 0x4d, 0x4e, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x30, 0x12,
	0x13, 0x0a, 0x0f, 0x47, 0x45, 0x54, 0x5f, 0x4e, 0x45, 0x57, 0x53, 0x5f, 0x46, 0x41, 0x49, 0x4c,
	0x45, 0x44, 0x10, 0x31, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x45, 0x54, 0x5f, 0x41, 0x47, 0x52, 0x45,
	0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x32, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x45,
	0x54, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x33, 0x12,
	0x16, 0x0a, 0x12, 0x53, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x34, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x54, 0x5f, 0x49,
	0x4d, 0x41, 0x47, 0x45, 0x5f, 0x49, 0x52, 0x52, 0x45, 0x47, 0x55, 0x4c, 0x41, 0x52, 0x5f, 0x46,
	0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x35, 0x12, 0x20, 0x0a, 0x1c, 0x53, 0x45, 0x54, 0x5f, 0x43,
	0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x49, 0x52, 0x52, 0x45, 0x47, 0x55, 0x4c, 0x41, 0x52,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x36, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x41, 0x4e,
	0x43, 0x45, 0x4c, 0x5f, 0x41, 0x47, 0x52, 0x45, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44,
	0x10, 0x37, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x56, 0x49, 0x45,
	0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x38, 0x12, 0x19, 0x0a, 0x15, 0x43, 0x41,
	0x4e, 0x43, 0x45, 0x4c, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x39, 0x12, 0x18, 0x0a, 0x14, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x49, 0x53, 0x54, 0x49, 0x43, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x3a, 0x12,
	0x1e, 0x0a, 0x1a, 0x47, 0x45, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x49, 0x53, 0x54, 0x49, 0x43,
	0x5f, 0x4a, 0x55, 0x44, 0x47, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x3b, 0x12,
	0x14, 0x0a, 0x10, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x3c, 0x12, 0x1c, 0x0a, 0x18, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x52, 0x45,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45,
	0x44, 0x10, 0x3d, 0x12, 0x1b, 0x0a, 0x17, 0x47, 0x45, 0x54, 0x5f, 0x49, 0x4d, 0x41, 0x47, 0x45,
	0x5f, 0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x3e,
	0x12, 0x1d, 0x0a, 0x19, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f,
	0x52, 0x45, 0x56, 0x49, 0x45, 0x57, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x3f, 0x12,
	0x15, 0x0a, 0x11, 0x53, 0x45, 0x54, 0x5f, 0x52, 0x45, 0x43, 0x4f, 0x52, 0x44, 0x5f, 0x46, 0x41,
	0x49, 0x4c, 0x45, 0x44, 0x10, 0x40, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x4d,
	0x50, 0x54, 0x59, 0x10, 0x41, 0x12, 0x16, 0x0a, 0x12, 0x41, 0x44, 0x44, 0x5f, 0x43, 0x4f, 0x4d,
	0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x42, 0x12, 0x19, 0x0a,
	0x15, 0x52, 0x45, 0x44, 0x55, 0x43, 0x45, 0x5f, 0x43, 0x4f, 0x4d, 0x4d, 0x45, 0x4e, 0x54, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x43, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x1f,
	0x5a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x62, 0x3b, 0x76, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
