// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/environment.proto

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

type Environment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source *Source `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"` // wip: Multiple environments (gateway needed)
}

func (x *Environment) Reset() {
	*x = Environment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_environment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Environment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Environment) ProtoMessage() {}

func (x *Environment) ProtoReflect() protoreflect.Message {
	mi := &file_proto_environment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Environment.ProtoReflect.Descriptor instead.
func (*Environment) Descriptor() ([]byte, []int) {
	return file_proto_environment_proto_rawDescGZIP(), []int{0}
}

func (x *Environment) GetSource() *Source {
	if x != nil {
		return x.Source
	}
	return nil
}

type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Provider:
	//
	//	*Source_Github
	Provider isSource_Provider `protobuf_oneof:"provider"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_environment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_proto_environment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_proto_environment_proto_rawDescGZIP(), []int{1}
}

func (m *Source) GetProvider() isSource_Provider {
	if m != nil {
		return m.Provider
	}
	return nil
}

func (x *Source) GetGithub() *GitHubSource {
	if x, ok := x.GetProvider().(*Source_Github); ok {
		return x.Github
	}
	return nil
}

type isSource_Provider interface {
	isSource_Provider()
}

type Source_Github struct {
	Github *GitHubSource `protobuf:"bytes,1,opt,name=github,proto3,oneof"`
}

func (*Source_Github) isSource_Provider() {}

type GitHubSource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Owner       string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Repo        string `protobuf:"bytes,2,opt,name=repo,proto3" json:"repo,omitempty"`
	Branch      string `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`
	AccessToken string `protobuf:"bytes,4,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *GitHubSource) Reset() {
	*x = GitHubSource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_environment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GitHubSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GitHubSource) ProtoMessage() {}

func (x *GitHubSource) ProtoReflect() protoreflect.Message {
	mi := &file_proto_environment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GitHubSource.ProtoReflect.Descriptor instead.
func (*GitHubSource) Descriptor() ([]byte, []int) {
	return file_proto_environment_proto_rawDescGZIP(), []int{2}
}

func (x *GitHubSource) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *GitHubSource) GetRepo() string {
	if x != nil {
		return x.Repo
	}
	return ""
}

func (x *GitHubSource) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *GitHubSource) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

var File_proto_environment_proto protoreflect.FileDescriptor

var file_proto_environment_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72,
	0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x3a, 0x0a, 0x0b, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x22, 0x49, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x33, 0x0a, 0x06,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x69, 0x74, 0x48, 0x75,
	0x62, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x48, 0x00, 0x52, 0x06, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x42, 0x0a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x22, 0x73, 0x0a,
	0x0c, 0x47, 0x69, 0x74, 0x48, 0x75, 0x62, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12,
	0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x7a, 0x65, 0x61, 0x62, 0x75, 0x72, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_environment_proto_rawDescOnce sync.Once
	file_proto_environment_proto_rawDescData = file_proto_environment_proto_rawDesc
)

func file_proto_environment_proto_rawDescGZIP() []byte {
	file_proto_environment_proto_rawDescOnce.Do(func() {
		file_proto_environment_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_environment_proto_rawDescData)
	})
	return file_proto_environment_proto_rawDescData
}

var file_proto_environment_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_environment_proto_goTypes = []interface{}{
	(*Environment)(nil),  // 0: environment.Environment
	(*Source)(nil),       // 1: environment.Source
	(*GitHubSource)(nil), // 2: environment.GitHubSource
}
var file_proto_environment_proto_depIdxs = []int32{
	1, // 0: environment.Environment.source:type_name -> environment.Source
	2, // 1: environment.Source.github:type_name -> environment.GitHubSource
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_environment_proto_init() }
func file_proto_environment_proto_init() {
	if File_proto_environment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_environment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Environment); i {
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
		file_proto_environment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Source); i {
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
		file_proto_environment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GitHubSource); i {
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
	file_proto_environment_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Source_Github)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_environment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_environment_proto_goTypes,
		DependencyIndexes: file_proto_environment_proto_depIdxs,
		MessageInfos:      file_proto_environment_proto_msgTypes,
	}.Build()
	File_proto_environment_proto = out.File
	file_proto_environment_proto_rawDesc = nil
	file_proto_environment_proto_goTypes = nil
	file_proto_environment_proto_depIdxs = nil
}
