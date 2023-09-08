package compact_test

import (
	"crypto/rand"
	"sync"
	"testing"

	api "github.com/kocubinski/costor-api"
	"github.com/kocubinski/costor-api/compact"
	"github.com/stretchr/testify/require"
)

func Test_WriteAndRead(t *testing.T) {
	const iterations = 1_000

	nodeDir := t.TempDir()
	nodeCtx := compact.StreamingContext{
		OutDir:       nodeDir,
		MaxFileSize:  1024 * 1024 * 100,
		OrderedInput: true,
		In:           make(chan compact.Sequenced),
	}

	errDir := t.TempDir()
	errCtx := compact.StreamingContext{
		OutDir:       errDir,
		MaxFileSize:  1024 * 1024 * 100,
		OrderedInput: true,
		In:           make(chan compact.Sequenced),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		stats, err := nodeCtx.Compact()
		require.NoError(t, err)
		stats.Report()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		stats, err := errCtx.Compact()
		require.NoError(t, err)
		stats.Report()
	}()

	for i := 0; i < iterations; i++ {
		k := make([]byte, 20)
		v := make([]byte, 1000)
		_, err := rand.Read(k)
		require.NoError(t, err)
		_, err = rand.Read(v)
		require.NoError(t, err)
		node := &api.Node{
			Key:      k,
			Value:    v,
			Block:    int64(i),
			StoreKey: "test",
		}
		nodeCtx.In <- node
		errCtx.In <- &api.DecodeError{Node: node, Reason: "testing", StoreKey: "test"}
	}
	close(nodeCtx.In)
	close(errCtx.In)
	wg.Wait()

	nodeItr, err := nodeCtx.NewIterator(nodeDir)
	require.NoError(t, err)
	cnt := 0
	for ; nodeItr.Valid(); err = nodeItr.Next() {
		require.NoError(t, err)
		node := nodeItr.Node
		require.Equal(t, "test", node.StoreKey)
		require.Equal(t, int64(cnt), node.Block)
		cnt++
	}
	require.Equal(t, iterations, cnt)

	cnt = 0
	errItr, err := compact.NewSequencedIterator(errDir, func() *api.DecodeError { return &api.DecodeError{} })
	require.NoError(t, err)
	for ; errItr.Valid(); err = errItr.Next() {
		require.NoError(t, err)
		node := errItr.Node
		require.Equal(t, "test", node.StoreKey)
		require.Equal(t, int64(cnt), node.Node.Block)
		require.Equal(t, "testing", node.Reason)
		cnt++
	}
	require.Equal(t, iterations, cnt)
}
