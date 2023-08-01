package bank

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// KVStore keys
var (
	StoreKey = "bank"
	// BalancesPrefix is the prefix for the account balances store. We use a byte
	// (instead of `[]byte("balances")` to save some disk space).
	BalancesPrefix             byte = 0x02
	SupplyKey                  byte = 0x00
	DenomMetadataPrefix        byte = 0x1
	DenomMetadataReversePrefix byte = 0x03
	SupplyOffsetKey            byte = 0x88
)

var (
	_ proto.Message = &Balance{}
	_ proto.Message = &Supply{}
	_ proto.Message = &SupplyOffset{}
	_ proto.Message = &DenomMetadata{}
	_ proto.Message = &Metadata{}
	_ proto.Message = &DenomUnit{}
	_ proto.Message = &ReverseDenomMetadata{}
)

var BalanceShortName = reflect.TypeOf((*Balance)(nil)).Elem().String()

type Balance struct {
	Address  []byte `protobuf:"bytes,1,opt,name=address,proto3" json:"-"`
	Denom    string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	Path     string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
	Amount   string `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	ICS20    string `protobuf:"bytes,5,opt,name=ics20,proto3" json:"ics20,omitempty"`
	Migrated bool   `protobuf:"varint,6,opt,name=migrated,proto3" json:"-"`

	Bech32Address string `json:"address"`
}

func (b *Balance) Reset() { *b = Balance{} }

func (b *Balance) String() string {
	return fmt.Sprintf("address: %s, denom: %s, path: %s, amount: %s",
		hex.EncodeToString(b.Address), b.Denom, b.Path, b.Amount)
}

func (b *Balance) ProtoMessage() {}

type Supply struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (s *Supply) Reset() {}

func (s *Supply) String() string {
	return ""
}

func (s *Supply) ProtoMessage() {}

type SupplyOffset struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (s SupplyOffset) Reset() {}

func (s SupplyOffset) String() string { return "" }

func (s SupplyOffset) ProtoMessage() {}

type DenomMetadata struct {
	Base     string    `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Metadata *Metadata `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (d DenomMetadata) Reset() {}

func (d DenomMetadata) String() string { return "" }

func (d DenomMetadata) ProtoMessage() {}

type Metadata struct {
	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
	// denom_units represents the list of DenomUnit's for a given coin
	DenomUnits []*DenomUnit `protobuf:"bytes,2,rep,name=denom_units,json=denomUnits,proto3" json:"denom_units,omitempty"`
	// base represents the base denom (should be the DenomUnit with exponent = 0).
	Base string `protobuf:"bytes,3,opt,name=base,proto3" json:"base,omitempty"`
	// display indicates the suggested denom that should be
	// displayed in clients.
	Display string `protobuf:"bytes,4,opt,name=display,proto3" json:"display,omitempty"`
	// name defines the name of the token (eg: Cosmos Atom)
	//
	// Since: cosmos-sdk 0.43
	Name string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	// symbol is the token symbol usually shown on exchanges (eg: ATOM). This can
	// be the same as the display.
	//
	// Since: cosmos-sdk 0.43
	Symbol string `protobuf:"bytes,6,opt,name=symbol,proto3" json:"symbol,omitempty"`
}

func (m Metadata) Reset() {}

func (m Metadata) String() string { return "" }

func (m Metadata) ProtoMessage() {}

type DenomUnit struct {
	// denom represents the string name of the given denom unit (e.g uatom).
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	// exponent represents power of 10 exponent that one must
	// raise the base_denom to in order to equal the given DenomUnit's denom
	// 1 denom = 1^exponent base_denom
	// (e.g. with a base_denom of uatom, one can create a DenomUnit of 'atom' with
	// exponent = 6, thus: 1 atom = 10^6 uatom).
	Exponent uint32 `protobuf:"varint,2,opt,name=exponent,proto3" json:"exponent,omitempty"`
	// aliases is a list of string aliases for the given denom
	Aliases []string `protobuf:"bytes,3,rep,name=aliases,proto3" json:"aliases,omitempty"`
}

func (d DenomUnit) Reset() {}

func (d DenomUnit) String() string { return "" }

func (d DenomUnit) ProtoMessage() {}

type ReverseDenomMetadata struct {
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Base  string `protobuf:"bytes,2,opt,name=base,proto3" json:"base,omitempty"`
}

func (r ReverseDenomMetadata) Reset() {}

func (r ReverseDenomMetadata) String() string {
	return fmt.Sprintf("denom: %s, base: %s", r.Denom, r.Base)
}

func (r ReverseDenomMetadata) ProtoMessage() {}
