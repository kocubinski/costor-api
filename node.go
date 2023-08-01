package api

func (n *Node) IsOrphaned() bool {
	return n.FirstVersion > 0
}
