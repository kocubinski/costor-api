package compact

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	api "github.com/kocubinski/costor-api"
	"github.com/kocubinski/costor-api/logz"
)

var log = logz.Logger.With().Str("module", "compact").Logger()

type StreamingIterator struct {
	*api.Node
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
		log.Info().Msgf("open file: %s", filepath.Base(nextFile))
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
