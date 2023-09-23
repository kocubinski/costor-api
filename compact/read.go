package compact

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	api "github.com/kocubinski/costor-api"
	"github.com/kocubinski/costor-api/logz"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

var log = logz.Logger.With().Str("module", "compact").Logger()

type SequencedIterator[T Sequenced] struct {
	Node T
	Err  error

	valid     bool
	newNodeFn func() T
	log       zerolog.Logger
	nextFile  chan string
	file      *os.File
	zr        *gzip.Reader
	// debug
	idx        int
	totalNodes int
	totalBytes int64
}

func NewSequencedIterator[T Sequenced](dir string, newNode func() T) (*SequencedIterator[T], error) {
	ch, err := newIteratorChannel(dir)
	if err != nil {
		return nil, err
	}
	itr := &SequencedIterator[T]{
		nextFile:  ch,
		log:       log.With().Str("path", dir).Logger(),
		newNodeFn: newNode,
	}
	return itr, itr.Next()
}

func (it *SequencedIterator[T]) GetNode() T {
	return it.Node
}

func (it *SequencedIterator[T]) Valid() bool {
	return it.valid
}

func (it *SequencedIterator[T]) Next() error {
	var err error

	if it.file == nil {
		nextFile, ok := <-it.nextFile
		// end of iteration
		if !ok {
			it.valid = false
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
	node := it.newNodeFn()
	if err := proto.Unmarshal(nbz, node); err != nil {
		return err
	}
	it.totalBytes += int64(length)
	it.idx += length
	it.Node = node
	it.valid = true
	return nil
}

func (c *StreamingContext) NewIterator(dir string) (*SequencedIterator[*api.Node], error) {
	return NewSequencedIterator[*api.Node](dir, func() *api.Node { return &api.Node{} })
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

var _ api.NodeIterator = (*StoreKeyedIterator)(nil)

// StoreKeyedIterator iterates over all nodes in a directory, setting each node's StoreKey to the given value.
type StoreKeyedIterator struct {
	// if set, set each node's StoreKey to this
	StoreKey string

	itr     *SequencedIterator[*api.Node]
	paused  bool
	version int64
}

func (s *StoreKeyedIterator) Next() error {
	if s.paused {
		return nil
	}
	if err := s.itr.Next(); err != nil {
		return err
	}

	if s.itr.Node != nil && s.StoreKey != "" {
		s.itr.Node.StoreKey = s.StoreKey
	}
	if s.itr.Node.Block > s.version {
		s.paused = true
	}

	return nil
}

func (s *StoreKeyedIterator) Valid() bool {
	if s == nil {
		return false
	}
	if s.itr == nil {
		return false
	}
	return !s.paused && s.itr.Valid()
}

func (s *StoreKeyedIterator) GetNode() *api.Node {
	if s.paused {
		return nil
	}
	return s.itr.Node
}

type ChangesetIterator struct {
	nodes    *StoreKeyedIterator
	version  int64
	nextNode *api.Node
}

func NewChangesetIterator(dir string, storeKey ...string) (*ChangesetIterator, error) {
	streamCtx := &StreamingContext{}
	itr, err := streamCtx.NewIterator(dir)
	if err != nil {
		return nil, err
	}
	skItr := &StoreKeyedIterator{itr: itr, paused: true}
	if len(storeKey) > 0 {
		skItr.StoreKey = storeKey[0]
	}
	changeItr := &ChangesetIterator{
		nodes: skItr,
	}
	err = changeItr.Next()
	if err != nil {
		return nil, err
	}
	return changeItr, nil
}

func (it *ChangesetIterator) Next() error {
	if !it.nodes.Valid() && !it.nodes.paused {
		it.nodes = nil
		return nil
	}
	if it.nodes.paused {
		it.nodes.paused = false
		it.nodes.version = it.nodes.GetNode().Block
	} else {
		return fmt.Errorf("expected paused iterator")
	}
	return nil
}

func (it *ChangesetIterator) Valid() bool {
	return it.nodes.Valid()
}

func (it *ChangesetIterator) Nodes() api.NodeIterator {
	return it.nodes
}

func (it *ChangesetIterator) Version() int64 {
	return it.version
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
	//it.Changeset = &api.Changeset{Version: math.MaxInt64}
	//for _, itr := range it.iterators {
	//	if itr.Valid() {
	//		if itr.Version < it.Version {
	//			it.Version = itr.Version
	//		}
	//	}
	//}
	//for _, itr := range it.iterators {
	//	if itr.Version == it.Version {
	//		it.Nodes = append(it.Nodes, itr.Nodes...)
	//		err := itr.Next()
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}
	//return nil
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

func (it *MulitChangesetIterator) Nodes() api.NodeIterator {
	return nil
	//return it.Changeset
}

func (it *MulitChangesetIterator) Version() int64 {
	//return it.Version
	return 0
}
