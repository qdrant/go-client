// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: snapshots_service.proto

package qdrant

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateFullSnapshotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateFullSnapshotRequest) Reset() {
	*x = CreateFullSnapshotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFullSnapshotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFullSnapshotRequest) ProtoMessage() {}

func (x *CreateFullSnapshotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFullSnapshotRequest.ProtoReflect.Descriptor instead.
func (*CreateFullSnapshotRequest) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{0}
}

type ListFullSnapshotsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListFullSnapshotsRequest) Reset() {
	*x = ListFullSnapshotsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListFullSnapshotsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListFullSnapshotsRequest) ProtoMessage() {}

func (x *ListFullSnapshotsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListFullSnapshotsRequest.ProtoReflect.Descriptor instead.
func (*ListFullSnapshotsRequest) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{1}
}

type DeleteFullSnapshotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SnapshotName string `protobuf:"bytes,1,opt,name=snapshot_name,json=snapshotName,proto3" json:"snapshot_name,omitempty"` // Name of the full snapshot
}

func (x *DeleteFullSnapshotRequest) Reset() {
	*x = DeleteFullSnapshotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFullSnapshotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFullSnapshotRequest) ProtoMessage() {}

func (x *DeleteFullSnapshotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFullSnapshotRequest.ProtoReflect.Descriptor instead.
func (*DeleteFullSnapshotRequest) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteFullSnapshotRequest) GetSnapshotName() string {
	if x != nil {
		return x.SnapshotName
	}
	return ""
}

type CreateSnapshotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionName string `protobuf:"bytes,1,opt,name=collection_name,json=collectionName,proto3" json:"collection_name,omitempty"` // Name of the collection
}

func (x *CreateSnapshotRequest) Reset() {
	*x = CreateSnapshotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSnapshotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSnapshotRequest) ProtoMessage() {}

func (x *CreateSnapshotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSnapshotRequest.ProtoReflect.Descriptor instead.
func (*CreateSnapshotRequest) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{3}
}

func (x *CreateSnapshotRequest) GetCollectionName() string {
	if x != nil {
		return x.CollectionName
	}
	return ""
}

type ListSnapshotsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionName string `protobuf:"bytes,1,opt,name=collection_name,json=collectionName,proto3" json:"collection_name,omitempty"` // Name of the collection
}

func (x *ListSnapshotsRequest) Reset() {
	*x = ListSnapshotsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSnapshotsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSnapshotsRequest) ProtoMessage() {}

func (x *ListSnapshotsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSnapshotsRequest.ProtoReflect.Descriptor instead.
func (*ListSnapshotsRequest) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{4}
}

func (x *ListSnapshotsRequest) GetCollectionName() string {
	if x != nil {
		return x.CollectionName
	}
	return ""
}

type DeleteSnapshotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionName string `protobuf:"bytes,1,opt,name=collection_name,json=collectionName,proto3" json:"collection_name,omitempty"` // Name of the collection
	SnapshotName   string `protobuf:"bytes,2,opt,name=snapshot_name,json=snapshotName,proto3" json:"snapshot_name,omitempty"`       // Name of the collection snapshot
}

func (x *DeleteSnapshotRequest) Reset() {
	*x = DeleteSnapshotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSnapshotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSnapshotRequest) ProtoMessage() {}

func (x *DeleteSnapshotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSnapshotRequest.ProtoReflect.Descriptor instead.
func (*DeleteSnapshotRequest) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteSnapshotRequest) GetCollectionName() string {
	if x != nil {
		return x.CollectionName
	}
	return ""
}

func (x *DeleteSnapshotRequest) GetSnapshotName() string {
	if x != nil {
		return x.SnapshotName
	}
	return ""
}

type SnapshotDescription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                                     // Name of the snapshot
	CreationTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=creation_time,json=creationTime,proto3" json:"creation_time,omitempty"` // Creation time of the snapshot
	Size         int64                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`                                    // Size of the snapshot in bytes
	Checksum     *string                `protobuf:"bytes,4,opt,name=checksum,proto3,oneof" json:"checksum,omitempty"`                       // SHA256 digest of the snapshot file
}

func (x *SnapshotDescription) Reset() {
	*x = SnapshotDescription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SnapshotDescription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SnapshotDescription) ProtoMessage() {}

func (x *SnapshotDescription) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SnapshotDescription.ProtoReflect.Descriptor instead.
func (*SnapshotDescription) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{6}
}

func (x *SnapshotDescription) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SnapshotDescription) GetCreationTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreationTime
	}
	return nil
}

func (x *SnapshotDescription) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *SnapshotDescription) GetChecksum() string {
	if x != nil && x.Checksum != nil {
		return *x.Checksum
	}
	return ""
}

type CreateSnapshotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SnapshotDescription *SnapshotDescription `protobuf:"bytes,1,opt,name=snapshot_description,json=snapshotDescription,proto3" json:"snapshot_description,omitempty"`
	Time                float64              `protobuf:"fixed64,2,opt,name=time,proto3" json:"time,omitempty"` // Time spent to process
}

func (x *CreateSnapshotResponse) Reset() {
	*x = CreateSnapshotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSnapshotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSnapshotResponse) ProtoMessage() {}

func (x *CreateSnapshotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSnapshotResponse.ProtoReflect.Descriptor instead.
func (*CreateSnapshotResponse) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{7}
}

func (x *CreateSnapshotResponse) GetSnapshotDescription() *SnapshotDescription {
	if x != nil {
		return x.SnapshotDescription
	}
	return nil
}

func (x *CreateSnapshotResponse) GetTime() float64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type ListSnapshotsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SnapshotDescriptions []*SnapshotDescription `protobuf:"bytes,1,rep,name=snapshot_descriptions,json=snapshotDescriptions,proto3" json:"snapshot_descriptions,omitempty"`
	Time                 float64                `protobuf:"fixed64,2,opt,name=time,proto3" json:"time,omitempty"` // Time spent to process
}

func (x *ListSnapshotsResponse) Reset() {
	*x = ListSnapshotsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListSnapshotsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSnapshotsResponse) ProtoMessage() {}

func (x *ListSnapshotsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListSnapshotsResponse.ProtoReflect.Descriptor instead.
func (*ListSnapshotsResponse) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{8}
}

func (x *ListSnapshotsResponse) GetSnapshotDescriptions() []*SnapshotDescription {
	if x != nil {
		return x.SnapshotDescriptions
	}
	return nil
}

func (x *ListSnapshotsResponse) GetTime() float64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type DeleteSnapshotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time float64 `protobuf:"fixed64,1,opt,name=time,proto3" json:"time,omitempty"` // Time spent to process
}

func (x *DeleteSnapshotResponse) Reset() {
	*x = DeleteSnapshotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_snapshots_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteSnapshotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSnapshotResponse) ProtoMessage() {}

func (x *DeleteSnapshotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_snapshots_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteSnapshotResponse.ProtoReflect.Descriptor instead.
func (*DeleteSnapshotResponse) Descriptor() ([]byte, []int) {
	return file_snapshots_service_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteSnapshotResponse) GetTime() float64 {
	if x != nil {
		return x.Time
	}
	return 0
}

var File_snapshots_service_proto protoreflect.FileDescriptor

var file_snapshots_service_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x71, 0x64, 0x72, 0x61, 0x6e,
	0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x1b, 0x0a, 0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x75, 0x6c, 0x6c,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x1a, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x75, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x40, 0x0a, 0x19, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x40, 0x0a,
	0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x3f, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65,
	0x22, 0x65, 0x0a, 0x15, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68,
	0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0xac, 0x01, 0x0a, 0x13, 0x53, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x3f, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1f, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x73, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x22, 0x7c, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x4e, 0x0a, 0x14, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x5f, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x71, 0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x13, 0x73, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04,
	0x74, 0x69, 0x6d, 0x65, 0x22, 0x7d, 0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6e, 0x61, 0x70,
	0x73, 0x68, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a,
	0x15, 0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x71,
	0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x44, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x14, 0x73, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x22, 0x2c, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x32, 0xdd, 0x03, 0x0a, 0x09, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x12,
	0x49, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x71, 0x64, 0x72, 0x61,
	0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x71, 0x64, 0x72, 0x61, 0x6e,
	0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x04, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x1c, 0x2e, 0x71, 0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x71, 0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x49, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x71, 0x64,
	0x72, 0x61, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x71, 0x64, 0x72,
	0x61, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68,
	0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0a,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x12, 0x21, 0x2e, 0x71, 0x64, 0x72,
	0x61, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e,
	0x71, 0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x4d, 0x0a, 0x08, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x75, 0x6c, 0x6c, 0x12, 0x20, 0x2e, 0x71, 0x64,
	0x72, 0x61, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x46, 0x75, 0x6c, 0x6c, 0x53, 0x6e, 0x61,
	0x70, 0x73, 0x68, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e,
	0x71, 0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51,
	0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x75, 0x6c, 0x6c, 0x12, 0x21, 0x2e, 0x71,
	0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x75, 0x6c, 0x6c,
	0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1e, 0x2e, 0x71, 0x64, 0x72, 0x61, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_snapshots_service_proto_rawDescOnce sync.Once
	file_snapshots_service_proto_rawDescData = file_snapshots_service_proto_rawDesc
)

func file_snapshots_service_proto_rawDescGZIP() []byte {
	file_snapshots_service_proto_rawDescOnce.Do(func() {
		file_snapshots_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_snapshots_service_proto_rawDescData)
	})
	return file_snapshots_service_proto_rawDescData
}

var file_snapshots_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_snapshots_service_proto_goTypes = []any{
	(*CreateFullSnapshotRequest)(nil), // 0: qdrant.CreateFullSnapshotRequest
	(*ListFullSnapshotsRequest)(nil),  // 1: qdrant.ListFullSnapshotsRequest
	(*DeleteFullSnapshotRequest)(nil), // 2: qdrant.DeleteFullSnapshotRequest
	(*CreateSnapshotRequest)(nil),     // 3: qdrant.CreateSnapshotRequest
	(*ListSnapshotsRequest)(nil),      // 4: qdrant.ListSnapshotsRequest
	(*DeleteSnapshotRequest)(nil),     // 5: qdrant.DeleteSnapshotRequest
	(*SnapshotDescription)(nil),       // 6: qdrant.SnapshotDescription
	(*CreateSnapshotResponse)(nil),    // 7: qdrant.CreateSnapshotResponse
	(*ListSnapshotsResponse)(nil),     // 8: qdrant.ListSnapshotsResponse
	(*DeleteSnapshotResponse)(nil),    // 9: qdrant.DeleteSnapshotResponse
	(*timestamppb.Timestamp)(nil),     // 10: google.protobuf.Timestamp
}
var file_snapshots_service_proto_depIdxs = []int32{
	10, // 0: qdrant.SnapshotDescription.creation_time:type_name -> google.protobuf.Timestamp
	6,  // 1: qdrant.CreateSnapshotResponse.snapshot_description:type_name -> qdrant.SnapshotDescription
	6,  // 2: qdrant.ListSnapshotsResponse.snapshot_descriptions:type_name -> qdrant.SnapshotDescription
	3,  // 3: qdrant.Snapshots.Create:input_type -> qdrant.CreateSnapshotRequest
	4,  // 4: qdrant.Snapshots.List:input_type -> qdrant.ListSnapshotsRequest
	5,  // 5: qdrant.Snapshots.Delete:input_type -> qdrant.DeleteSnapshotRequest
	0,  // 6: qdrant.Snapshots.CreateFull:input_type -> qdrant.CreateFullSnapshotRequest
	1,  // 7: qdrant.Snapshots.ListFull:input_type -> qdrant.ListFullSnapshotsRequest
	2,  // 8: qdrant.Snapshots.DeleteFull:input_type -> qdrant.DeleteFullSnapshotRequest
	7,  // 9: qdrant.Snapshots.Create:output_type -> qdrant.CreateSnapshotResponse
	8,  // 10: qdrant.Snapshots.List:output_type -> qdrant.ListSnapshotsResponse
	9,  // 11: qdrant.Snapshots.Delete:output_type -> qdrant.DeleteSnapshotResponse
	7,  // 12: qdrant.Snapshots.CreateFull:output_type -> qdrant.CreateSnapshotResponse
	8,  // 13: qdrant.Snapshots.ListFull:output_type -> qdrant.ListSnapshotsResponse
	9,  // 14: qdrant.Snapshots.DeleteFull:output_type -> qdrant.DeleteSnapshotResponse
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_snapshots_service_proto_init() }
func file_snapshots_service_proto_init() {
	if File_snapshots_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_snapshots_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateFullSnapshotRequest); i {
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
		file_snapshots_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ListFullSnapshotsRequest); i {
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
		file_snapshots_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteFullSnapshotRequest); i {
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
		file_snapshots_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*CreateSnapshotRequest); i {
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
		file_snapshots_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ListSnapshotsRequest); i {
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
		file_snapshots_service_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteSnapshotRequest); i {
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
		file_snapshots_service_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*SnapshotDescription); i {
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
		file_snapshots_service_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*CreateSnapshotResponse); i {
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
		file_snapshots_service_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*ListSnapshotsResponse); i {
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
		file_snapshots_service_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteSnapshotResponse); i {
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
	file_snapshots_service_proto_msgTypes[6].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_snapshots_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_snapshots_service_proto_goTypes,
		DependencyIndexes: file_snapshots_service_proto_depIdxs,
		MessageInfos:      file_snapshots_service_proto_msgTypes,
	}.Build()
	File_snapshots_service_proto = out.File
	file_snapshots_service_proto_rawDesc = nil
	file_snapshots_service_proto_goTypes = nil
	file_snapshots_service_proto_depIdxs = nil
}
