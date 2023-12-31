// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.2
// source: store.proto

package api

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

// StoreKVPair duplicates the Cosmos SDK proto spec at:
// https://github.com/cosmos/cosmos-sdk/blob/33eead4d42049bf8bf3f5ceacd0af256a284d342/proto/cosmos/store/v1beta1/listening.proto#L13-L18
//
// This is done to prevent forming a dependency on the SDK's proto files in the API.  Osmosis streaming writes using a
// (now removed) StoreKVPairs struct (defined below) so it is needed to here unmarhsal the data.
type StoreKVPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreKey string `protobuf:"bytes,1,opt,name=store_key,json=storeKey,proto3" json:"store_key,omitempty"` // the store key for the KVStore this pair originates from
	Delete   bool   `protobuf:"varint,2,opt,name=delete,proto3" json:"delete,omitempty"`                    // true indicates a delete operation, false indicates a set operation
	Key      []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *StoreKVPair) Reset() {
	*x = StoreKVPair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreKVPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreKVPair) ProtoMessage() {}

func (x *StoreKVPair) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreKVPair.ProtoReflect.Descriptor instead.
func (*StoreKVPair) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{0}
}

func (x *StoreKVPair) GetStoreKey() string {
	if x != nil {
		return x.StoreKey
	}
	return ""
}

func (x *StoreKVPair) GetDelete() bool {
	if x != nil {
		return x.Delete
	}
	return false
}

func (x *StoreKVPair) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *StoreKVPair) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// StoreKVPairs is the format which pre 0.45x state streaming writes as (osmosis)
type StoreKVPairs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pairs       []*StoreKVPair `protobuf:"bytes,1,rep,name=pairs,proto3" json:"pairs,omitempty"`
	BlockHeight int64          `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	StoreKey    string         `protobuf:"bytes,3,opt,name=store_key,json=storeKey,proto3" json:"store_key,omitempty"`
	KeyPrefix   []byte         `protobuf:"bytes,4,opt,name=key_prefix,json=keyPrefix,proto3" json:"key_prefix,omitempty"`
}

func (x *StoreKVPairs) Reset() {
	*x = StoreKVPairs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreKVPairs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreKVPairs) ProtoMessage() {}

func (x *StoreKVPairs) ProtoReflect() protoreflect.Message {
	mi := &file_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreKVPairs.ProtoReflect.Descriptor instead.
func (*StoreKVPairs) Descriptor() ([]byte, []int) {
	return file_store_proto_rawDescGZIP(), []int{1}
}

func (x *StoreKVPairs) GetPairs() []*StoreKVPair {
	if x != nil {
		return x.Pairs
	}
	return nil
}

func (x *StoreKVPairs) GetBlockHeight() int64 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *StoreKVPairs) GetStoreKey() string {
	if x != nil {
		return x.StoreKey
	}
	return ""
}

func (x *StoreKVPairs) GetKeyPrefix() []byte {
	if x != nil {
		return x.KeyPrefix
	}
	return nil
}

var File_store_proto protoreflect.FileDescriptor

var file_store_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6a, 0x0a,
	0x0b, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4b, 0x56, 0x50, 0x61, 0x69, 0x72, 0x12, 0x1b, 0x0a, 0x09,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x91, 0x01, 0x0a, 0x0c, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x4b, 0x56, 0x50, 0x61, 0x69, 0x72, 0x73, 0x12, 0x22, 0x0a, 0x05, 0x70, 0x61,
	0x69, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x4b, 0x56, 0x50, 0x61, 0x69, 0x72, 0x52, 0x05, 0x70, 0x61, 0x69, 0x72, 0x73, 0x12, 0x21,
	0x0a, 0x0c, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x4b, 0x65, 0x79, 0x12, 0x1d,
	0x0a, 0x0a, 0x6b, 0x65, 0x79, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x6b, 0x65, 0x79, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x42, 0x26, 0x5a,
	0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x6f, 0x63, 0x75,
	0x62, 0x69, 0x6e, 0x73, 0x6b, 0x69, 0x2f, 0x63, 0x6f, 0x73, 0x74, 0x6f, 0x72, 0x2d, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_proto_rawDescOnce sync.Once
	file_store_proto_rawDescData = file_store_proto_rawDesc
)

func file_store_proto_rawDescGZIP() []byte {
	file_store_proto_rawDescOnce.Do(func() {
		file_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_proto_rawDescData)
	})
	return file_store_proto_rawDescData
}

var file_store_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_store_proto_goTypes = []interface{}{
	(*StoreKVPair)(nil),  // 0: StoreKVPair
	(*StoreKVPairs)(nil), // 1: StoreKVPairs
}
var file_store_proto_depIdxs = []int32{
	0, // 0: StoreKVPairs.pairs:type_name -> StoreKVPair
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_store_proto_init() }
func file_store_proto_init() {
	if File_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreKVPair); i {
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
		file_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreKVPairs); i {
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
			RawDescriptor: file_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_proto_goTypes,
		DependencyIndexes: file_store_proto_depIdxs,
		MessageInfos:      file_store_proto_msgTypes,
	}.Build()
	File_store_proto = out.File
	file_store_proto_rawDesc = nil
	file_store_proto_goTypes = nil
	file_store_proto_depIdxs = nil
}
