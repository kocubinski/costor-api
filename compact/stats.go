package compact

import (
	"fmt"
	"strings"

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
