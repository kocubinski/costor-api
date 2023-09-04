package compact

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	api "github.com/kocubinski/costor-api"
	"github.com/kocubinski/costor-api/logz"
	"github.com/rs/zerolog"
)

var log = logz.Logger.With().Str("module", "compact").Logger()

type StreamingIterator struct {
	*api.Node
	log      zerolog.Logger
	nextFile chan string
	file     *os.File
	zr       *gzip.Reader
	Err      error
	// debug
	idx        int
	totalNodes int
	totalBytes int64
}

func (it *StreamingIterator) GetNode() *api.Node {
	return it.Node
}

func (it *StreamingIterator) Valid() bool {
	return it.Node != nil
}

func (it *StreamingIterator) Next() error {
	var err error

	if it.file == nil {
		nextFile, ok := <-it.nextFile
		// end of iteration
		if !ok {
			it.Node = nil
			return nil
		}
		it.log.Info().Msgf("open file: %s", filepath.Base(nextFile))
		it.file, err = os.Open(nextFile)
		it.zr, err = gzip.NewReader(it.file)
		if err != nil {
			return err
		}
	}

	lengthBz := make([]byte, 4)
	err = binary.Read(it.zr, binary.LittleEndian, &lengthBz)
	if err != nil {
		if err == io.EOF {
			if err := it.file.Close(); err != nil {
				return err
			}
			it.file = nil
			it.idx = 0
			return it.Next()
		}
		return err
	}
	it.totalBytes += 4
	it.idx += 4
	length := int(binary.LittleEndian.Uint32(lengthBz))

	nbz := make([]byte, length)
	n, err := it.zr.Read(nbz)
	if err != nil && err != io.EOF {
		return err
	}
	// fill rest of the bytes if the decompressor returned less than length
	if n != length {
		rest := length - n
		for {
			bz := make([]byte, rest)
			m, err := it.zr.Read(bz)
			if err != nil && err != io.EOF {
				return err
			}
			if m == rest {
				copy(nbz[length-rest:], bz)
				break
			}
			if m == 0 {
				return fmt.Errorf("read 0 bytes, needed %d", rest)
			}
			copy(nbz[length-rest:], bz[:m])
			rest -= m
		}
	}

	it.totalNodes++
	node := &api.Node{}
	if err := proto.Unmarshal(nbz, node); err != nil {
		return err
	}
	it.totalBytes += int64(length)
	it.idx += length
	it.Node = node
	return nil
}

func (c *StreamingContext) NewIterator(dir string) (*StreamingIterator, error) {
	ch, err := newIteratorChannel(dir)
	if err != nil {
		return nil, err
	}
	itr := &StreamingIterator{
		nextFile: ch,
		log:      log.With().Str("path", dir).Logger(),
	}
	return itr, itr.Next()
}

// newIteratorChannel returns a channel that enumerates over all files in a directory.
// the go thread parks between each file read, so it's not very efficient.
func newIteratorChannel(dir string) (chan string, error) {
	ch := make(chan string)
	go func() {
		err := filepath.WalkDir(dir, func(path string, info os.DirEntry, readErr error) error {
			if readErr != nil {
				return readErr
			}
			if info.IsDir() {
				return nil
			}
			ch <- fmt.Sprintf("%s/%s", dir, info.Name())
			return nil
		})
		if err != nil {
			panic(err)
		}
		close(ch)
	}()
	return ch, nil
}

type ChangesetIterator struct {
	*api.Changeset
	// if set, set each node's StoreKey to this
	StoreKey string

	nodeItr *StreamingIterator
}

func NewChangesetIterator(dir string, storeKey ...string) (*ChangesetIterator, error) {
	streamCtx := &StreamingContext{}
	itr, err := streamCtx.NewIterator(dir)
	if err != nil {
		return nil, err
	}
	changeItr := &ChangesetIterator{
		nodeItr: itr,
	}
	if len(storeKey) > 0 {
		changeItr.StoreKey = storeKey[0]
	}
	err = changeItr.Next()
	if err != nil {
		return nil, err
	}
	return changeItr, nil
}

func (it *ChangesetIterator) Next() error {
	it.Changeset = &api.Changeset{}
	var err error
	for ; it.nodeItr.Valid(); err = it.nodeItr.Next() {
		if err != nil {
			return err
		}
		node := it.nodeItr.Node
		if it.StoreKey != "" {
			node.StoreKey = it.StoreKey
		}
		if it.Version == 0 {
			it.Version = node.Block
		}
		if node.Block > it.Version {
			break
		}
		it.Nodes = append(it.Nodes, node)
	}

	return nil
}

func (it *ChangesetIterator) Valid() bool {
	return it.nodeItr.Valid()
}

func (it *ChangesetIterator) GetChangeset() *api.Changeset {
	return it.Changeset
}

type MulitChangesetIterator struct {
	*api.Changeset
	iterators []*ChangesetIterator
}

func NewMultiChangesetIterator(dir string) (*MulitChangesetIterator, error) {
	multiItr := &MulitChangesetIterator{
		Changeset: &api.Changeset{},
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			return nil, fmt.Errorf("expected directory, got file: %s", file.Name())
		}
		itr, err := NewChangesetIterator(fmt.Sprintf("%s/%s", dir, file.Name()), file.Name())
		if err != nil {
			return nil, err
		}
		multiItr.iterators = append(multiItr.iterators, itr)
	}
	err = multiItr.Next()
	if err != nil {
		return nil, err
	}
	return multiItr, nil
}

func (it *MulitChangesetIterator) Next() error {
	it.Changeset = &api.Changeset{Version: math.MaxInt64}
	for _, itr := range it.iterators {
		if itr.Valid() {
			if itr.Version < it.Version {
				it.Version = itr.Version
			}
		}
	}
	for _, itr := range it.iterators {
		if itr.Version == it.Version {
			it.Nodes = append(it.Nodes, itr.Nodes...)
			err := itr.Next()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (it *MulitChangesetIterator) Valid() bool {
	for _, itr := range it.iterators {
		if !itr.Valid() {
			return false
		}
	}
	return true
}

func (it *MulitChangesetIterator) GetChangeset() *api.Changeset {
	return it.Changeset
}
