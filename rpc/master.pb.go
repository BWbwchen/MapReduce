// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.0
// source: rpc/master.proto

package rpc

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

type WorkerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid string `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Ip   string `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
}

func (x *WorkerInfo) Reset() {
	*x = WorkerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_master_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerInfo) ProtoMessage() {}

func (x *WorkerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_master_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerInfo.ProtoReflect.Descriptor instead.
func (*WorkerInfo) Descriptor() ([]byte, []int) {
	return file_rpc_master_proto_rawDescGZIP(), []int{0}
}

func (x *WorkerInfo) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *WorkerInfo) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

type RegisterResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result bool  `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Id     int64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RegisterResult) Reset() {
	*x = RegisterResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_master_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResult) ProtoMessage() {}

func (x *RegisterResult) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_master_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResult.ProtoReflect.Descriptor instead.
func (*RegisterResult) Descriptor() ([]byte, []int) {
	return file_rpc_master_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterResult) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

func (x *RegisterResult) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type IMDInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid      string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Filenames []string `protobuf:"bytes,2,rep,name=filenames,proto3" json:"filenames,omitempty"`
}

func (x *IMDInfo) Reset() {
	*x = IMDInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_master_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IMDInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IMDInfo) ProtoMessage() {}

func (x *IMDInfo) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_master_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IMDInfo.ProtoReflect.Descriptor instead.
func (*IMDInfo) Descriptor() ([]byte, []int) {
	return file_rpc_master_proto_rawDescGZIP(), []int{2}
}

func (x *IMDInfo) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *IMDInfo) GetFilenames() []string {
	if x != nil {
		return x.Filenames
	}
	return nil
}

type UpdateResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result bool `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *UpdateResult) Reset() {
	*x = UpdateResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_master_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateResult) ProtoMessage() {}

func (x *UpdateResult) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_master_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateResult.ProtoReflect.Descriptor instead.
func (*UpdateResult) Descriptor() ([]byte, []int) {
	return file_rpc_master_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateResult) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

var File_rpc_master_proto protoreflect.FileDescriptor

var file_rpc_master_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x30, 0x0a, 0x0a, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x75, 0x75, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x70, 0x22, 0x38, 0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3b,
	0x0a, 0x07, 0x49, 0x4d, 0x44, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x26, 0x0a, 0x0c, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x32, 0x62, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x2e, 0x0a,
	0x0e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x0b, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0f, 0x2e, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x28, 0x0a,
	0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x4d, 0x44, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x08,
	0x2e, 0x49, 0x4d, 0x44, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0d, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x3b, 0x72, 0x70,
	0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_master_proto_rawDescOnce sync.Once
	file_rpc_master_proto_rawDescData = file_rpc_master_proto_rawDesc
)

func file_rpc_master_proto_rawDescGZIP() []byte {
	file_rpc_master_proto_rawDescOnce.Do(func() {
		file_rpc_master_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_master_proto_rawDescData)
	})
	return file_rpc_master_proto_rawDescData
}

var file_rpc_master_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_master_proto_goTypes = []interface{}{
	(*WorkerInfo)(nil),     // 0: WorkerInfo
	(*RegisterResult)(nil), // 1: RegisterResult
	(*IMDInfo)(nil),        // 2: IMDInfo
	(*UpdateResult)(nil),   // 3: UpdateResult
}
var file_rpc_master_proto_depIdxs = []int32{
	0, // 0: Master.WorkerRegister:input_type -> WorkerInfo
	2, // 1: Master.UpdateIMDInfo:input_type -> IMDInfo
	1, // 2: Master.WorkerRegister:output_type -> RegisterResult
	3, // 3: Master.UpdateIMDInfo:output_type -> UpdateResult
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_master_proto_init() }
func file_rpc_master_proto_init() {
	if File_rpc_master_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_master_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerInfo); i {
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
		file_rpc_master_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterResult); i {
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
		file_rpc_master_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IMDInfo); i {
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
		file_rpc_master_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateResult); i {
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
			RawDescriptor: file_rpc_master_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_master_proto_goTypes,
		DependencyIndexes: file_rpc_master_proto_depIdxs,
		MessageInfos:      file_rpc_master_proto_msgTypes,
	}.Build()
	File_rpc_master_proto = out.File
	file_rpc_master_proto_rawDesc = nil
	file_rpc_master_proto_goTypes = nil
	file_rpc_master_proto_depIdxs = nil
}
