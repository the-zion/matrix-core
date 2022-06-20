// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: message/service/v1/message.proto

package v1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type AvatarReviewReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobsDetail *AvatarReviewReq_JobsDetailStruct `protobuf:"bytes,1,opt,name=JobsDetail,proto3" json:"JobsDetail,omitempty"`
	EventName  string                            `protobuf:"bytes,2,opt,name=EventName,proto3" json:"EventName,omitempty"`
}

func (x *AvatarReviewReq) Reset() {
	*x = AvatarReviewReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AvatarReviewReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AvatarReviewReq) ProtoMessage() {}

func (x *AvatarReviewReq) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AvatarReviewReq.ProtoReflect.Descriptor instead.
func (*AvatarReviewReq) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{0}
}

func (x *AvatarReviewReq) GetJobsDetail() *AvatarReviewReq_JobsDetailStruct {
	if x != nil {
		return x.JobsDetail
	}
	return nil
}

func (x *AvatarReviewReq) GetEventName() string {
	if x != nil {
		return x.EventName
	}
	return ""
}

type AvatarReviewReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AvatarReviewReply) Reset() {
	*x = AvatarReviewReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AvatarReviewReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AvatarReviewReply) ProtoMessage() {}

func (x *AvatarReviewReply) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AvatarReviewReply.ProtoReflect.Descriptor instead.
func (*AvatarReviewReply) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{1}
}

type ProfileReviewReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data    *ProfileReviewReq_Data `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ProfileReviewReq) Reset() {
	*x = ProfileReviewReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileReviewReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileReviewReq) ProtoMessage() {}

func (x *ProfileReviewReq) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileReviewReq.ProtoReflect.Descriptor instead.
func (*ProfileReviewReq) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{2}
}

func (x *ProfileReviewReq) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ProfileReviewReq) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *ProfileReviewReq) GetData() *ProfileReviewReq_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type ProfileReviewReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ProfileReviewReply) Reset() {
	*x = ProfileReviewReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileReviewReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileReviewReply) ProtoMessage() {}

func (x *ProfileReviewReply) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileReviewReply.ProtoReflect.Descriptor instead.
func (*ProfileReviewReply) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{3}
}

type AvatarReviewReq_JobsDetailStruct struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code       string            `protobuf:"bytes,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Message    string            `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
	JobId      string            `protobuf:"bytes,3,opt,name=JobId,proto3" json:"JobId,omitempty"`
	State      string            `protobuf:"bytes,4,opt,name=State,proto3" json:"State,omitempty"`
	Object     string            `protobuf:"bytes,5,opt,name=Object,proto3" json:"Object,omitempty"`
	Label      string            `protobuf:"bytes,6,opt,name=Label,proto3" json:"Label,omitempty"`
	Result     int32             `protobuf:"varint,7,opt,name=Result,proto3" json:"Result,omitempty"`
	Category   string            `protobuf:"bytes,8,opt,name=Category,proto3" json:"Category,omitempty"`
	BucketId   string            `protobuf:"bytes,9,opt,name=BucketId,proto3" json:"BucketId,omitempty"`
	Region     string            `protobuf:"bytes,10,opt,name=Region,proto3" json:"Region,omitempty"`
	CosHeaders map[string]string `protobuf:"bytes,11,rep,name=CosHeaders,proto3" json:"CosHeaders,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *AvatarReviewReq_JobsDetailStruct) Reset() {
	*x = AvatarReviewReq_JobsDetailStruct{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AvatarReviewReq_JobsDetailStruct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AvatarReviewReq_JobsDetailStruct) ProtoMessage() {}

func (x *AvatarReviewReq_JobsDetailStruct) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AvatarReviewReq_JobsDetailStruct.ProtoReflect.Descriptor instead.
func (*AvatarReviewReq_JobsDetailStruct) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{0, 0}
}

func (x *AvatarReviewReq_JobsDetailStruct) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *AvatarReviewReq_JobsDetailStruct) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetBucketId() string {
	if x != nil {
		return x.BucketId
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *AvatarReviewReq_JobsDetailStruct) GetCosHeaders() map[string]string {
	if x != nil {
		return x.CosHeaders
	}
	return nil
}

type ProfileReviewReq_PornInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HitFlag int32  `protobuf:"varint,1,opt,name=hit_flag,json=hitFlag,proto3" json:"hit_flag,omitempty"`
	Label   string `protobuf:"bytes,2,opt,name=label,proto3" json:"label,omitempty"`
	Count   int32  `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ProfileReviewReq_PornInfo) Reset() {
	*x = ProfileReviewReq_PornInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileReviewReq_PornInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileReviewReq_PornInfo) ProtoMessage() {}

func (x *ProfileReviewReq_PornInfo) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileReviewReq_PornInfo.ProtoReflect.Descriptor instead.
func (*ProfileReviewReq_PornInfo) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{2, 0}
}

func (x *ProfileReviewReq_PornInfo) GetHitFlag() int32 {
	if x != nil {
		return x.HitFlag
	}
	return 0
}

func (x *ProfileReviewReq_PornInfo) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *ProfileReviewReq_PornInfo) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ProfileReviewReq_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ForbiddenStatus string                     `protobuf:"bytes,1,opt,name=forbidden_status,json=forbiddenStatus,proto3" json:"forbidden_status,omitempty"`
	PornInfo        *ProfileReviewReq_PornInfo `protobuf:"bytes,2,opt,name=porn_info,json=pornInfo,proto3" json:"porn_info,omitempty"`
	Result          int32                      `protobuf:"varint,3,opt,name=result,proto3" json:"result,omitempty"`
	TraceId         string                     `protobuf:"bytes,4,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	Url             string                     `protobuf:"bytes,5,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ProfileReviewReq_Data) Reset() {
	*x = ProfileReviewReq_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_service_v1_message_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProfileReviewReq_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProfileReviewReq_Data) ProtoMessage() {}

func (x *ProfileReviewReq_Data) ProtoReflect() protoreflect.Message {
	mi := &file_message_service_v1_message_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProfileReviewReq_Data.ProtoReflect.Descriptor instead.
func (*ProfileReviewReq_Data) Descriptor() ([]byte, []int) {
	return file_message_service_v1_message_proto_rawDescGZIP(), []int{2, 1}
}

func (x *ProfileReviewReq_Data) GetForbiddenStatus() string {
	if x != nil {
		return x.ForbiddenStatus
	}
	return ""
}

func (x *ProfileReviewReq_Data) GetPornInfo() *ProfileReviewReq_PornInfo {
	if x != nil {
		return x.PornInfo
	}
	return nil
}

func (x *ProfileReviewReq_Data) GetResult() int32 {
	if x != nil {
		return x.Result
	}
	return 0
}

func (x *ProfileReviewReq_Data) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *ProfileReviewReq_Data) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

var File_message_service_v1_message_proto protoreflect.FileDescriptor

var file_message_service_v1_message_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x92, 0x05, 0x0a, 0x0f, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72,
	0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x12, 0x4c, 0x0a, 0x0a, 0x4a, 0x6f, 0x62,
	0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x2e, 0x4a, 0x6f, 0x62, 0x73, 0x44,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x0a, 0x4a, 0x6f, 0x62,
	0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x26, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0x18, 0xe8, 0x07, 0x52, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x1a,
	0x88, 0x04, 0x0a, 0x10, 0x4a, 0x6f, 0x62, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x53, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x12, 0x1c, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52, 0x04, 0x43, 0x6f,
	0x64, 0x65, 0x12, 0x22, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52, 0x07, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1e, 0x0a, 0x05, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52,
	0x05, 0x4a, 0x6f, 0x62, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52,
	0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07,
	0x52, 0x06, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x1e, 0x0a, 0x05, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8,
	0x07, 0x52, 0x05, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x24, 0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52, 0x08, 0x43, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x24, 0x0a, 0x08, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x49, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x18,
	0xe8, 0x07, 0x52, 0x08, 0x42, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x06,
	0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42,
	0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52, 0x06, 0x52, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x6b,
	0x0a, 0x0a, 0x43, 0x6f, 0x73, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x0b, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x3c, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x2e,
	0x4a, 0x6f, 0x62, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x43, 0x6f, 0x73, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x42, 0x0d, 0xfa, 0x42, 0x0a, 0x9a, 0x01, 0x07, 0x2a, 0x05, 0x72, 0x03, 0x18, 0xe8, 0x07, 0x52,
	0x0a, 0x43, 0x6f, 0x73, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x3d, 0x0a, 0x0f, 0x43,
	0x6f, 0x73, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x13, 0x0a, 0x11, 0x41, 0x76,
	0x61, 0x74, 0x61, 0x72, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x87, 0x03, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65,
	0x77, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x35, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x51, 0x0a, 0x08, 0x50, 0x6f, 0x72,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x69, 0x74, 0x5f, 0x66, 0x6c, 0x61,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x68, 0x69, 0x74, 0x46, 0x6c, 0x61, 0x67,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x1a, 0xba, 0x01, 0x0a,
	0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x29, 0x0a, 0x10, 0x66, 0x6f, 0x72, 0x62, 0x69, 0x64, 0x64,
	0x65, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x66, 0x6f, 0x72, 0x62, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x42, 0x0a, 0x09, 0x70, 0x6f, 0x72, 0x6e, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65,
	0x71, 0x2e, 0x50, 0x6f, 0x72, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x08, 0x70, 0x6f, 0x72, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x14, 0x0a, 0x12, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32,
	0xc5, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x4c, 0x0a, 0x0c, 0x41,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x1b, 0x2e, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x52,
	0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x52, 0x65, 0x76, 0x69,
	0x65, 0x77, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x6c, 0x0a, 0x0d, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x76, 0x69, 0x65, 0x77, 0x12, 0x1c, 0x2e, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x1a, 0x1e, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x76,
	0x69, 0x65, 0x77, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17,
	0x22, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x72, 0x65,
	0x76, 0x69, 0x65, 0x77, 0x3a, 0x01, 0x2a, 0x42, 0x1b, 0x5a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76,
	0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_service_v1_message_proto_rawDescOnce sync.Once
	file_message_service_v1_message_proto_rawDescData = file_message_service_v1_message_proto_rawDesc
)

func file_message_service_v1_message_proto_rawDescGZIP() []byte {
	file_message_service_v1_message_proto_rawDescOnce.Do(func() {
		file_message_service_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_service_v1_message_proto_rawDescData)
	})
	return file_message_service_v1_message_proto_rawDescData
}

var file_message_service_v1_message_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_message_service_v1_message_proto_goTypes = []interface{}{
	(*AvatarReviewReq)(nil),                  // 0: message.v1.AvatarReviewReq
	(*AvatarReviewReply)(nil),                // 1: message.v1.AvatarReviewReply
	(*ProfileReviewReq)(nil),                 // 2: message.v1.ProfileReviewReq
	(*ProfileReviewReply)(nil),               // 3: message.v1.ProfileReviewReply
	(*AvatarReviewReq_JobsDetailStruct)(nil), // 4: message.v1.AvatarReviewReq.JobsDetailStruct
	nil,                                      // 5: message.v1.AvatarReviewReq.JobsDetailStruct.CosHeadersEntry
	(*ProfileReviewReq_PornInfo)(nil),        // 6: message.v1.ProfileReviewReq.PornInfo
	(*ProfileReviewReq_Data)(nil),            // 7: message.v1.ProfileReviewReq.Data
}
var file_message_service_v1_message_proto_depIdxs = []int32{
	4, // 0: message.v1.AvatarReviewReq.JobsDetail:type_name -> message.v1.AvatarReviewReq.JobsDetailStruct
	7, // 1: message.v1.ProfileReviewReq.data:type_name -> message.v1.ProfileReviewReq.Data
	5, // 2: message.v1.AvatarReviewReq.JobsDetailStruct.CosHeaders:type_name -> message.v1.AvatarReviewReq.JobsDetailStruct.CosHeadersEntry
	6, // 3: message.v1.ProfileReviewReq.Data.porn_info:type_name -> message.v1.ProfileReviewReq.PornInfo
	0, // 4: message.v1.Message.AvatarReview:input_type -> message.v1.AvatarReviewReq
	2, // 5: message.v1.Message.ProfileReview:input_type -> message.v1.ProfileReviewReq
	1, // 6: message.v1.Message.AvatarReview:output_type -> message.v1.AvatarReviewReply
	3, // 7: message.v1.Message.ProfileReview:output_type -> message.v1.ProfileReviewReply
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_message_service_v1_message_proto_init() }
func file_message_service_v1_message_proto_init() {
	if File_message_service_v1_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_service_v1_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AvatarReviewReq); i {
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
		file_message_service_v1_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AvatarReviewReply); i {
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
		file_message_service_v1_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileReviewReq); i {
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
		file_message_service_v1_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileReviewReply); i {
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
		file_message_service_v1_message_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AvatarReviewReq_JobsDetailStruct); i {
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
		file_message_service_v1_message_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileReviewReq_PornInfo); i {
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
		file_message_service_v1_message_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProfileReviewReq_Data); i {
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
			RawDescriptor: file_message_service_v1_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_message_service_v1_message_proto_goTypes,
		DependencyIndexes: file_message_service_v1_message_proto_depIdxs,
		MessageInfos:      file_message_service_v1_message_proto_msgTypes,
	}.Build()
	File_message_service_v1_message_proto = out.File
	file_message_service_v1_message_proto_rawDesc = nil
	file_message_service_v1_message_proto_goTypes = nil
	file_message_service_v1_message_proto_depIdxs = nil
}
