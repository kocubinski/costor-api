package api

import (
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/proto"
)

var (
	_ proto.Message = &Node{}
	_ proto.Message = &Nodes{}
	_ proto.Message = &DecodeError{}
)

type Node struct {
	Key      []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Delete   bool   `protobuf:"varint,3,opt,name=delete,proto3" json:"delete"`
	Block    int64  `protobuf:"varint,4,opt,name=block,proto3" json:"block,omitempty"`
	StoreKey string `protobuf:"bytes,5,opt,name=store_key,json=storeKey,proto3" json:"store_key,omitempty"`

	// v2 //

	Hash []byte `protobuf:"bytes,6,opt,name=hash,proto3" json:"hash,omitempty"`

	// orphaned nodes only
	FirstVersion int64 `protobuf:"varint,7,opt,name=first_version,json=firstVersion,proto3" json:"first_version,omitempty"`
	LastVersion  int64 `protobuf:"varint,8,opt,name=last_version,json=lastVersion,proto3" json:"last_version,omitempty"`

	// not marshalled
	Height int8
	Size   int64
}

type Nodes struct {
	Nodes []*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (n *Node) Reset() {
	*n = Node{}
}

func (n *Node) String() string {
	return fmt.Sprintf("key: %s, value: %s, delete: %t, block: %d, store_key: %s, hash: %s, first_version %d, last_version %d",
		hex.EncodeToString(n.Key), hex.EncodeToString(n.Value), n.Delete, n.Block, n.StoreKey,
		hex.EncodeToString(n.Hash), n.FirstVersion, n.LastVersion)
}

func (n *Node) ProtoMessage() {}

func (n *Nodes) Reset() {
	*n = Nodes{}
}

func (n *Nodes) String() string {
	return fmt.Sprintf("nodes: %v", n.Nodes)
}

func (n *Nodes) ProtoMessage() {}

type DecodeError struct {
	Err  string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	Node *Node  `protobuf:"bytes,2,opt,name=node,proto3" json:"node,omitempty"`
}

func (e DecodeError) Error() string {
	return e.Err
}
func (e *DecodeError) ProtoMessage() {}

func (e *DecodeError) Reset() {
	*e = DecodeError{}
}

func (e *DecodeError) String() string {
	return fmt.Sprintf("err: %s, node: %v", e.Err, e.Node)
}

func DecodeErr(err string, node *Node) error {
	return DecodeError{
		Err:  err,
		Node: node,
	}
}

func IsDecodeError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(DecodeError)
	return ok
}
