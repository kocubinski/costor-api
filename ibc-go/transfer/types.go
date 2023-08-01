package transfer

import (
	"encoding/hex"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	StoreKey      = "transfer"
	DenomTraceKey = []byte{0x02}
)

type DenomTrace struct {
	// path defines the chain of port/channel identifiers used for tracing the
	// source of the fungible token.
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// base denomination of the relayed fungible token.
	BaseDenom string `protobuf:"bytes,2,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
}

func (m *DenomTrace) Reset()         { *m = DenomTrace{} }
func (m *DenomTrace) String() string { return "" }
func (*DenomTrace) ProtoMessage()    {}

// ParseHexHash parses a hex hash in string format to bytes and validates its correctness.
func ParseHexHash(hexHash string) (tmbytes.HexBytes, error) {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, err
	}

	if err := tmtypes.ValidateHash(hash); err != nil {
		return nil, err
	}

	return hash, nil
}
