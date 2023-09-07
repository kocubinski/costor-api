package compact

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/proto"
	api "github.com/kocubinski/costor-api"
	"github.com/kocubinski/costor-api/core"
	"github.com/kocubinski/costor-api/logz"
)

func prettyByteSize(b int64) string {
	bf := float64(b)
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}

func (c *StreamingContext) PrintDebug() {
	log.Info().Msgf("minBlock: %d, maxBlock: %d", c.minBlock, c.maxBlock)
}

type StreamingContext struct {
	core.Context
	OutDir             string
	In                 chan *api.Node
	MaxFileSize        int
	WorkerId           int
	ExpectedEfficiency float64
	FileSeq            int

	// Assume that the input is ordered by block height
	OrderedInput bool

	minBlock int64
	maxBlock int64
}

func (c *StreamingContext) nextFilename() (string, error) {
	var filename string
	if c.OrderedInput {
		var base string
		if c.minBlock == c.maxBlock {
			base = fmt.Sprintf("%s/%08d", c.OutDir, c.minBlock)
		} else {
			base = fmt.Sprintf("%s/%08d-%08d", c.OutDir, c.minBlock, c.maxBlock)
		}
		filename = fmt.Sprintf("%s.pb.gz", base)
		if api.IsFileExistent(filename) {
			filename = fmt.Sprintf("%s-%08d.pb.gz", base, c.FileSeq)
		}
		if api.IsFileExistent(filename) {
			log.Error().Msg(fmt.Sprintf("file %s already exists", filename))
			return "", fmt.Errorf("file %s already exists", filename)
		}
	} else {
		filename = fmt.Sprintf("%s/%02d-%08d.pb.gz", c.OutDir, c.WorkerId, c.FileSeq)
	}

	return filename, nil
}

func (c *StreamingContext) Compact() (*Stats, error) {
	logger := logz.Logger.With().Str("module", "streaming").Logger()
	c.minBlock = math.MaxInt64
	c.maxBlock = 0
	stats := &Stats{}
	var (
		buf    bytes.Buffer
		gz     = gzip.NewWriter(&buf)
		uzSize int
	)

	flush := func() error {
		if buf.Len() == 0 {
			return nil
		}
		err := gz.Close()
		if err != nil {
			return err
		}
		stats.BytesWritten += int64(buf.Len())

		filename, err := c.nextFilename()
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
		logger.Info().Msg(fmt.Sprintf("wrote %s", filepath.Base(filename)))

		stats.FilesWritten = append(stats.FilesWritten, filename)
		buf.Reset()
		gz = gzip.NewWriter(&buf)
		uzSize = 0
		c.FileSeq++
		c.minBlock = math.MaxInt64
		c.maxBlock = 0
		return nil
	}

	var flushErr error
	defer func() {
		flushErr = flush()
	}()

	for node := range c.In {
		stats.NodeCount++
		if node.Block < c.minBlock {
			c.minBlock = node.Block
		}
		if node.Block > c.maxBlock {
			c.maxBlock = node.Block
		}

		protoBz, err := proto.Marshal(node)
		if err != nil {
			return nil, err
		}
		stats.BytesRead += int64(len(protoBz))
		err = binary.Write(gz, binary.LittleEndian, uint32(len(protoBz)))
		if err != nil {
			return nil, err
		}
		uzSize += 4

		_, err = gz.Write(protoBz)
		if err != nil {
			return nil, err
		}
		uzSize += len(protoBz)

		if buf.Len() > c.MaxFileSize {
			err := flush()
			if err != nil {
				return nil, err
			}
		}
	}

	return stats, flushErr
}
