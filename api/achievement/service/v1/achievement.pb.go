// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.0
// source: achievement/service/v1/achievement.proto

package v1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SetAchievementAgreeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *SetAchievementAgreeReq) Reset() {
	*x = SetAchievementAgreeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_achievement_service_v1_achievement_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetAchievementAgreeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetAchievementAgreeReq) ProtoMessage() {}

func (x *SetAchievementAgreeReq) ProtoReflect() protoreflect.Message {
	mi := &file_achievement_service_v1_achievement_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetAchievementAgreeReq.ProtoReflect.Descriptor instead.
func (*SetAchievementAgreeReq) Descriptor() ([]byte, []int) {
	return file_achievement_service_v1_achievement_proto_rawDescGZIP(), []int{0}
}

func (x *SetAchievementAgreeReq) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type CancelAchievementAgreeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *CancelAchievementAgreeReq) Reset() {
	*x = CancelAchievementAgreeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_achievement_service_v1_achievement_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelAchievementAgreeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelAchievementAgreeReq) ProtoMessage() {}

func (x *CancelAchievementAgreeReq) ProtoReflect() protoreflect.Message {
	mi := &file_achievement_service_v1_achievement_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelAchievementAgreeReq.ProtoReflect.Descriptor instead.
func (*CancelAchievementAgreeReq) Descriptor() ([]byte, []int) {
	return file_achievement_service_v1_achievement_proto_rawDescGZIP(), []int{1}
}

func (x *CancelAchievementAgreeReq) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type SetAchievementViewReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *SetAchievementViewReq) Reset() {
	*x = SetAchievementViewReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_achievement_service_v1_achievement_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetAchievementViewReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetAchievementViewReq) ProtoMessage() {}

func (x *SetAchievementViewReq) ProtoReflect() protoreflect.Message {
	mi := &file_achievement_service_v1_achievement_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetAchievementViewReq.ProtoReflect.Descriptor instead.
func (*SetAchievementViewReq) Descriptor() ([]byte, []int) {
	return file_achievement_service_v1_achievement_proto_rawDescGZIP(), []int{2}
}

func (x *SetAchievementViewReq) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type SetAchievementCollectReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *SetAchievementCollectReq) Reset() {
	*x = SetAchievementCollectReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_achievement_service_v1_achievement_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetAchievementCollectReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetAchievementCollectReq) ProtoMessage() {}

func (x *SetAchievementCollectReq) ProtoReflect() protoreflect.Message {
	mi := &file_achievement_service_v1_achievement_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetAchievementCollectReq.ProtoReflect.Descriptor instead.
func (*SetAchievementCollectReq) Descriptor() ([]byte, []int) {
	return file_achievement_service_v1_achievement_proto_rawDescGZIP(), []int{3}
}

func (x *SetAchievementCollectReq) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

type CancelAchievementCollectReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *CancelAchievementCollectReq) Reset() {
	*x = CancelAchievementCollectReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_achievement_service_v1_achievement_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelAchievementCollectReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelAchievementCollectReq) ProtoMessage() {}

func (x *CancelAchievementCollectReq) ProtoReflect() protoreflect.Message {
	mi := &file_achievement_service_v1_achievement_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelAchievementCollectReq.ProtoReflect.Descriptor instead.
func (*CancelAchievementCollectReq) Descriptor() ([]byte, []int) {
	return file_achievement_service_v1_achievement_proto_rawDescGZIP(), []int{4}
}

func (x *CancelAchievementCollectReq) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

var File_achievement_service_v1_achievement_proto protoreflect.FileDescriptor

var file_achievement_service_v1_achievement_proto_rawDesc = []byte{
	0x0a, 0x28, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x61, 0x63, 0x68, 0x69,
	0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36,
	0x0a, 0x16, 0x53, 0x65, 0x74, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x41, 0x67, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01,
	0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0x39, 0x0a, 0x19, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c,
	0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x67, 0x72, 0x65, 0x65,
	0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x04, 0x75, 0x75, 0x69,
	0x64, 0x22, 0x35, 0x0a, 0x15, 0x53, 0x65, 0x74, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x56, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0,
	0x01, 0x01, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x22, 0x38, 0x0a, 0x18, 0x53, 0x65, 0x74, 0x41,
	0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x04, 0x75, 0x75,
	0x69, 0x64, 0x22, 0x3b, 0x0a, 0x1b, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x41, 0x63, 0x68, 0x69,
	0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x12, 0x1c, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x32,
	0xbd, 0x04, 0x0a, 0x0b, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x57, 0x0a, 0x13, 0x53, 0x65, 0x74, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x41, 0x67, 0x72, 0x65, 0x65, 0x12, 0x26, 0x2e, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x41, 0x63, 0x68, 0x69, 0x65,
	0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x67, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x5d, 0x0a, 0x16, 0x43, 0x61, 0x6e, 0x63,
	0x65, 0x6c, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x67, 0x72,
	0x65, 0x65, 0x12, 0x29, 0x2e, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x67, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x12, 0x53, 0x65, 0x74, 0x41, 0x63,
	0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x69, 0x65, 0x77, 0x12, 0x25, 0x2e,
	0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x65, 0x74, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x69, 0x65,
	0x77, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x5b,
	0x0a, 0x15, 0x53, 0x65, 0x74, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x12, 0x28, 0x2e, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x74, 0x41, 0x63, 0x68, 0x69,
	0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x61, 0x0a, 0x18, 0x43,
	0x61, 0x6e, 0x63, 0x65, 0x6c, 0x41, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x12, 0x2b, 0x2e, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x41,
	0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x5f,
	0x0a, 0x09, 0x47, 0x65, 0x74, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x22, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x1c, 0x12, 0x1a, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x74, 0x2f, 0x61, 0x63, 0x68,
	0x69, 0x65, 0x76, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x42,
	0x1f, 0x5a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x63, 0x68, 0x69, 0x65, 0x76, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_achievement_service_v1_achievement_proto_rawDescOnce sync.Once
	file_achievement_service_v1_achievement_proto_rawDescData = file_achievement_service_v1_achievement_proto_rawDesc
)

func file_achievement_service_v1_achievement_proto_rawDescGZIP() []byte {
	file_achievement_service_v1_achievement_proto_rawDescOnce.Do(func() {
		file_achievement_service_v1_achievement_proto_rawDescData = protoimpl.X.CompressGZIP(file_achievement_service_v1_achievement_proto_rawDescData)
	})
	return file_achievement_service_v1_achievement_proto_rawDescData
}

var file_achievement_service_v1_achievement_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_achievement_service_v1_achievement_proto_goTypes = []interface{}{
	(*SetAchievementAgreeReq)(nil),      // 0: achievement.v1.SetAchievementAgreeReq
	(*CancelAchievementAgreeReq)(nil),   // 1: achievement.v1.CancelAchievementAgreeReq
	(*SetAchievementViewReq)(nil),       // 2: achievement.v1.SetAchievementViewReq
	(*SetAchievementCollectReq)(nil),    // 3: achievement.v1.SetAchievementCollectReq
	(*CancelAchievementCollectReq)(nil), // 4: achievement.v1.CancelAchievementCollectReq
	(*emptypb.Empty)(nil),               // 5: google.protobuf.Empty
}
var file_achievement_service_v1_achievement_proto_depIdxs = []int32{
	0, // 0: achievement.v1.Achievement.SetAchievementAgree:input_type -> achievement.v1.SetAchievementAgreeReq
	1, // 1: achievement.v1.Achievement.CancelAchievementAgree:input_type -> achievement.v1.CancelAchievementAgreeReq
	2, // 2: achievement.v1.Achievement.SetAchievementView:input_type -> achievement.v1.SetAchievementViewReq
	3, // 3: achievement.v1.Achievement.SetAchievementCollect:input_type -> achievement.v1.SetAchievementCollectReq
	4, // 4: achievement.v1.Achievement.CancelAchievementCollect:input_type -> achievement.v1.CancelAchievementCollectReq
	5, // 5: achievement.v1.Achievement.GetHealth:input_type -> google.protobuf.Empty
	5, // 6: achievement.v1.Achievement.SetAchievementAgree:output_type -> google.protobuf.Empty
	5, // 7: achievement.v1.Achievement.CancelAchievementAgree:output_type -> google.protobuf.Empty
	5, // 8: achievement.v1.Achievement.SetAchievementView:output_type -> google.protobuf.Empty
	5, // 9: achievement.v1.Achievement.SetAchievementCollect:output_type -> google.protobuf.Empty
	5, // 10: achievement.v1.Achievement.CancelAchievementCollect:output_type -> google.protobuf.Empty
	5, // 11: achievement.v1.Achievement.GetHealth:output_type -> google.protobuf.Empty
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_achievement_service_v1_achievement_proto_init() }
func file_achievement_service_v1_achievement_proto_init() {
	if File_achievement_service_v1_achievement_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_achievement_service_v1_achievement_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetAchievementAgreeReq); i {
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
		file_achievement_service_v1_achievement_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelAchievementAgreeReq); i {
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
		file_achievement_service_v1_achievement_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetAchievementViewReq); i {
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
		file_achievement_service_v1_achievement_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetAchievementCollectReq); i {
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
		file_achievement_service_v1_achievement_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelAchievementCollectReq); i {
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
			RawDescriptor: file_achievement_service_v1_achievement_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_achievement_service_v1_achievement_proto_goTypes,
		DependencyIndexes: file_achievement_service_v1_achievement_proto_depIdxs,
		MessageInfos:      file_achievement_service_v1_achievement_proto_msgTypes,
	}.Build()
	File_achievement_service_v1_achievement_proto = out.File
	file_achievement_service_v1_achievement_proto_rawDesc = nil
	file_achievement_service_v1_achievement_proto_goTypes = nil
	file_achievement_service_v1_achievement_proto_depIdxs = nil
}
