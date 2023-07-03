// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.3
// source: automation_function.proto

package proto

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

type AutomationFunctionMeta struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Short       string `protobuf:"bytes,1,opt,name=short,proto3" json:"short,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *AutomationFunctionMeta) Reset() {
	*x = AutomationFunctionMeta{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_function_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AutomationFunctionMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AutomationFunctionMeta) ProtoMessage() {}

func (x *AutomationFunctionMeta) ProtoReflect() protoreflect.Message {
	mi := &file_automation_function_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AutomationFunctionMeta.ProtoReflect.Descriptor instead.
func (*AutomationFunctionMeta) Descriptor() ([]byte, []int) {
	return file_automation_function_proto_rawDescGZIP(), []int{0}
}

func (x *AutomationFunctionMeta) GetShort() string {
	if x != nil {
		return x.Short
	}
	return ""
}

func (x *AutomationFunctionMeta) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type AutomationFunctionParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Types    []string `protobuf:"bytes,2,rep,name=types,proto3" json:"types,omitempty"`
	Required bool     `protobuf:"varint,3,opt,name=required,proto3" json:"required,omitempty"`
	Isarray  bool     `protobuf:"varint,4,opt,name=isarray,proto3" json:"isarray,omitempty"`
}

func (x *AutomationFunctionParams) Reset() {
	*x = AutomationFunctionParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_function_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AutomationFunctionParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AutomationFunctionParams) ProtoMessage() {}

func (x *AutomationFunctionParams) ProtoReflect() protoreflect.Message {
	mi := &file_automation_function_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AutomationFunctionParams.ProtoReflect.Descriptor instead.
func (*AutomationFunctionParams) Descriptor() ([]byte, []int) {
	return file_automation_function_proto_rawDescGZIP(), []int{1}
}

func (x *AutomationFunctionParams) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AutomationFunctionParams) GetTypes() []string {
	if x != nil {
		return x.Types
	}
	return nil
}

func (x *AutomationFunctionParams) GetRequired() bool {
	if x != nil {
		return x.Required
	}
	return false
}

func (x *AutomationFunctionParams) GetIsarray() bool {
	if x != nil {
		return x.Isarray
	}
	return false
}

type AutomationFunction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ref     string                      `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
	Kind    string                      `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Meta    *AutomationFunctionMeta     `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	Params  []*AutomationFunctionParams `protobuf:"bytes,4,rep,name=params,proto3" json:"params,omitempty"`
	Results []*AutomationFunctionParams `protobuf:"bytes,5,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *AutomationFunction) Reset() {
	*x = AutomationFunction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_function_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AutomationFunction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AutomationFunction) ProtoMessage() {}

func (x *AutomationFunction) ProtoReflect() protoreflect.Message {
	mi := &file_automation_function_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AutomationFunction.ProtoReflect.Descriptor instead.
func (*AutomationFunction) Descriptor() ([]byte, []int) {
	return file_automation_function_proto_rawDescGZIP(), []int{2}
}

func (x *AutomationFunction) GetRef() string {
	if x != nil {
		return x.Ref
	}
	return ""
}

func (x *AutomationFunction) GetKind() string {
	if x != nil {
		return x.Kind
	}
	return ""
}

func (x *AutomationFunction) GetMeta() *AutomationFunctionMeta {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *AutomationFunction) GetParams() []*AutomationFunctionParams {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *AutomationFunction) GetResults() []*AutomationFunctionParams {
	if x != nil {
		return x.Results
	}
	return nil
}

type ExecReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ExecReq) Reset() {
	*x = ExecReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_function_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecReq) ProtoMessage() {}

func (x *ExecReq) ProtoReflect() protoreflect.Message {
	mi := &file_automation_function_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecReq.ProtoReflect.Descriptor instead.
func (*ExecReq) Descriptor() ([]byte, []int) {
	return file_automation_function_proto_rawDescGZIP(), []int{3}
}

func (x *ExecReq) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type ExecResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ExecResp) Reset() {
	*x = ExecResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_function_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecResp) ProtoMessage() {}

func (x *ExecResp) ProtoReflect() protoreflect.Message {
	mi := &file_automation_function_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecResp.ProtoReflect.Descriptor instead.
func (*ExecResp) Descriptor() ([]byte, []int) {
	return file_automation_function_proto_rawDescGZIP(), []int{4}
}

func (x *ExecResp) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type MetaReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MetaReq) Reset() {
	*x = MetaReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_automation_function_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaReq) ProtoMessage() {}

func (x *MetaReq) ProtoReflect() protoreflect.Message {
	mi := &file_automation_function_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaReq.ProtoReflect.Descriptor instead.
func (*MetaReq) Descriptor() ([]byte, []int) {
	return file_automation_function_proto_rawDescGZIP(), []int{5}
}

var File_automation_function_proto protoreflect.FileDescriptor

var file_automation_function_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x50, 0x0a, 0x16, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x68, 0x6f,
	0x72, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x7a, 0x0a, 0x18, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x05, 0x74, 0x79, 0x70, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x61, 0x72, 0x72, 0x61,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x61, 0x72, 0x72, 0x61, 0x79,
	0x22, 0xe1, 0x01, 0x0a, 0x12, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46,
	0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x12, 0x31, 0x0a,
	0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x75,
	0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61,
	0x12, 0x37, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x39, 0x0a, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x52, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x22, 0x1f, 0x0a, 0x07, 0x45, 0x78, 0x65, 0x63, 0x52, 0x65, 0x71, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x20, 0x0a, 0x08, 0x45, 0x78, 0x65, 0x63, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x09, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x61, 0x52,
	0x65, 0x71, 0x32, 0x77, 0x0a, 0x19, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x27, 0x0a, 0x04, 0x45, 0x78, 0x65, 0x63, 0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x78, 0x65, 0x63, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x78, 0x65, 0x63, 0x52, 0x65, 0x73, 0x70, 0x12, 0x31, 0x0a, 0x04, 0x4d, 0x65, 0x74, 0x61,
	0x12, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x52, 0x65, 0x71,
	0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x6d, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x4b, 0x5a, 0x49, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x7a,
	0x61, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x7a, 0x61,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x6c, 0x75, 0x67,
	0x69, 0x6e, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_automation_function_proto_rawDescOnce sync.Once
	file_automation_function_proto_rawDescData = file_automation_function_proto_rawDesc
)

func file_automation_function_proto_rawDescGZIP() []byte {
	file_automation_function_proto_rawDescOnce.Do(func() {
		file_automation_function_proto_rawDescData = protoimpl.X.CompressGZIP(file_automation_function_proto_rawDescData)
	})
	return file_automation_function_proto_rawDescData
}

var file_automation_function_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_automation_function_proto_goTypes = []interface{}{
	(*AutomationFunctionMeta)(nil),   // 0: proto.AutomationFunctionMeta
	(*AutomationFunctionParams)(nil), // 1: proto.AutomationFunctionParams
	(*AutomationFunction)(nil),       // 2: proto.AutomationFunction
	(*ExecReq)(nil),                  // 3: proto.ExecReq
	(*ExecResp)(nil),                 // 4: proto.ExecResp
	(*MetaReq)(nil),                  // 5: proto.MetaReq
}
var file_automation_function_proto_depIdxs = []int32{
	0, // 0: proto.AutomationFunction.meta:type_name -> proto.AutomationFunctionMeta
	1, // 1: proto.AutomationFunction.params:type_name -> proto.AutomationFunctionParams
	1, // 2: proto.AutomationFunction.results:type_name -> proto.AutomationFunctionParams
	3, // 3: proto.AutomationFunctionService.Exec:input_type -> proto.ExecReq
	5, // 4: proto.AutomationFunctionService.Meta:input_type -> proto.MetaReq
	4, // 5: proto.AutomationFunctionService.Exec:output_type -> proto.ExecResp
	2, // 6: proto.AutomationFunctionService.Meta:output_type -> proto.AutomationFunction
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_automation_function_proto_init() }
func file_automation_function_proto_init() {
	if File_automation_function_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_automation_function_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AutomationFunctionMeta); i {
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
		file_automation_function_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AutomationFunctionParams); i {
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
		file_automation_function_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AutomationFunction); i {
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
		file_automation_function_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecReq); i {
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
		file_automation_function_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExecResp); i {
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
		file_automation_function_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaReq); i {
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
			RawDescriptor: file_automation_function_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_automation_function_proto_goTypes,
		DependencyIndexes: file_automation_function_proto_depIdxs,
		MessageInfos:      file_automation_function_proto_msgTypes,
	}.Build()
	File_automation_function_proto = out.File
	file_automation_function_proto_rawDesc = nil
	file_automation_function_proto_goTypes = nil
	file_automation_function_proto_depIdxs = nil
}