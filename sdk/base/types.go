package base

type Coin struct {
	Denom  string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Amount string `protobuf:"bytes,2,opt,name=amount,proto3,customtype=Int" json:"amount"`
}

func (m *Coin) Reset()      { *m = Coin{} }
func (*Coin) ProtoMessage() {}
func (m *Coin) String() string {
	return ""
}
