package api

func (n *Node) IsOrphaned() bool {
	return n.FirstVersion > 0
}

func (n *Node) Sequence() int64 {
	return n.Block
}

func (e *DecodeError) Sequence() int64 {
	if e.Node == nil {
		return 0
	}
	return e.Node.Block
}

type NodeIterator interface {
	Next() error
	Valid() bool
	GetNode() *Node
}
