package api_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	api "github.com/kocubinski/costor-api"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalNode(t *testing.T) {
	node := &api.Node{
		Key:      []byte("key"),
		Value:    []byte("value"),
		Delete:   true,
		Block:    1,
		StoreKey: "foo",
	}
	bz, err := proto.Marshal(node)
	require.NoError(t, err)
	require.NotEqual(t, 0, len(bz))

	var node2 api.Node
	err = proto.Unmarshal(bz, &node2)
	require.Equal(t, node, &node2)

	nodes := &api.Nodes{
		Nodes: []*api.Node{node, &node2},
	}
	bz, err = proto.Marshal(nodes)
	require.NoError(t, err)
	require.NotEqual(t, 0, len(bz))

	var nodes2 api.Nodes
	err = proto.Unmarshal(bz, &nodes2)
	require.Equal(t, nodes, &nodes2)
}

const testDbPath = "/Users/mattk/src/scratch/cosmosdb"

//
//func TestNodeLeafHash(t *testing.T) {
//	rs, err := store.NewReadonlyStore(testDbPath)
//	require.NoError(t, err)
//	prefixDb, _, err := rs.LatestTree("bank")
//	require.NoError(t, err)
//	itr, err := prefixDb.Iterator([]byte{'n'}, nil)
//	require.NoError(t, err)
//	i := 0
//	for ; itr.Valid(); itr.Next() {
//		k := itr.Key()
//		v := itr.Value()
//
//		// this is an anti-pattern. I shouldn't need to decode part of this nodes key to decode it.
//		node, err := iavl.MakeNode(nil, v)
//		require.NoError(t, err)
//		if !node.IsLeaf() {
//			continue
//		}
//		nodeHash, err := node.ComputeHash()
//		require.NoError(t, err)
//		costorHash, err := iavl.LeafNodeHash(node.GetLogicalKey(), node.GetValue(), node.GetVersion())
//		require.NoError(t, err)
//		require.Equal(t, nodeHash, costorHash)
//		require.Equal(t, k[1:], costorHash)
//
//		if i > 10 {
//			break
//		}
//		i++
//	}
//}
