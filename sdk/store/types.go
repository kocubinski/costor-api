package store

import "time"

type StoreKVPair struct {
	StoreKey string `protobuf:"bytes,1,opt,name=store_key,json=storeKey,proto3" json:"store_key,omitempty"`
	Delete   bool   `protobuf:"varint,2,opt,name=delete,proto3" json:"delete,omitempty"`
	Key      []byte `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

type StoreKVPairs struct {
	Pairs       []*StoreKVPair `protobuf:"bytes,1,rep,name=pairs,proto3" json:"pairs,omitempty"`
	BlockHeight int64          `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	StoreKey    string         `protobuf:"bytes,3,opt,name=store_key,json=storeKey,proto3" json:"store_key,omitempty"`
	KeyPrefix   []byte         `protobuf:"bytes,4,opt,name=key_prefix,json=keyPrefix,proto3" json:"key_prefix,omitempty"`
}

func (m *StoreKVPairs) Reset()         { *m = StoreKVPairs{} }
func (m *StoreKVPairs) String() string { return "" }
func (*StoreKVPairs) ProtoMessage()    {}

func (m *StoreKVPair) Reset()         { *m = StoreKVPair{} }
func (m *StoreKVPair) String() string { return "" }
func (*StoreKVPair) ProtoMessage()    {}

type CommitInfo struct {
	Version    int64        `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	StoreInfos []*StoreInfo `protobuf:"bytes,2,rep,name=store_infos,json=storeInfos,proto3" json:"store_infos"`
	Timestamp  *time.Time   `protobuf:"bytes,3,opt,name=timestamp,proto3,stdtime" json:"timestamp"`
}

func (m *CommitInfo) Reset()         { *m = CommitInfo{} }
func (m *CommitInfo) String() string { return "" }
func (*CommitInfo) ProtoMessage()    {}

type StoreInfo struct {
	Name     string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CommitId *CommitID `protobuf:"bytes,2,opt,name=commit_id,json=commitId,proto3" json:"commit_id"`
}

func (m *StoreInfo) Reset()         { *m = StoreInfo{} }
func (m *StoreInfo) String() string { return "" }
func (*StoreInfo) ProtoMessage()    {}

type CommitID struct {
	Version int64  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Hash    []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *CommitID) Reset() { *m = CommitID{} }

func (m *CommitID) String() string {
	return ""
}
func (*CommitID) ProtoMessage() {}
