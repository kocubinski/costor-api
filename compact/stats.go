package compact

import (
	"fmt"
	"sync"
)

type Stats struct {
	FilesWritten []string
	BytesRead    int64
	BytesWritten int64

	//
	WastedFlushes int
	NodeCount     int
}

func (stats *Stats) Report() {
	fmt.Printf("file count: %d\n", len(stats.FilesWritten))
	fmt.Printf("read: %s\n", prettyByteSize(stats.BytesRead))
	fmt.Printf("wrote: %s\n", prettyByteSize(stats.BytesWritten))
	fmt.Printf("node count: %d\n", stats.NodeCount)
	if len(stats.FilesWritten) == 0 {
		return
	}
	fmt.Printf("last file written: %s\n", stats.FilesWritten[len(stats.FilesWritten)-1])
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
