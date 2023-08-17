package compact

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
)

type Stats struct {
	FilesWritten []string
	BytesRead    int64
	BytesWritten int64

	//
	WastedFlushes int
	NodeCount     int
}

func (stats *Stats) Report() string {
	var sb strings.Builder
	sb.WriteString("compaction stats:\n")
	sb.WriteString(fmt.Sprintf("file count: %d\n", len(stats.FilesWritten)))
	sb.WriteString(fmt.Sprintf("read: %s\n", prettyByteSize(stats.BytesRead)))
	sb.WriteString(fmt.Sprintf("wrote: %s\n", prettyByteSize(stats.BytesWritten)))
	sb.WriteString(fmt.Sprintf("node count: %s\n", humanize.Comma(int64(stats.NodeCount))))
	if len(stats.FilesWritten) == 0 {
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("last file written: %s\n", stats.FilesWritten[len(stats.FilesWritten)-1]))
	return sb.String()
}

func (stats *Stats) Validate() error {
	if len(stats.FilesWritten) == 0 {
		return fmt.Errorf("no files written")
	}
	ch := make(chan string)
	wg := &sync.WaitGroup{}
	itr := &StreamingIterator{
		nextFile: ch,
	}

	var itrErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		for itrErr = itr.Next(); itr.Valid(); itrErr = itr.Next() {
			if itrErr != nil {
				break
			}
		}
	}()

	for _, filename := range stats.FilesWritten {
		ch <- filename
	}
	close(ch)

	wg.Wait()
	if itrErr != nil {
		return itrErr
	}

	if stats.NodeCount != itr.totalNodes {
		return fmt.Errorf("node count mismatch: %d != %d", stats.NodeCount, itr.totalNodes)
	}
	return nil
}
