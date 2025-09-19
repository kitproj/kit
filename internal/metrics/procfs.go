package metrics

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kitproj/kit/internal/types"
)

// GetProcFSCommand returns the cat command to read /proc/pid/statm
func GetProcFSCommand(pid int) []string {
	return []string{"cat", fmt.Sprintf("/proc/%d/statm", pid)}
}

// ParseProcFSOutput parses /proc/pid/statm output for memory usage only
func ParseProcFSOutput(output string) (*types.Metrics, error) {
	fields := strings.Fields(strings.TrimSpace(output))
	if len(fields) < 2 {
		return nil, fmt.Errorf("unexpected /proc/pid/statm output: insufficient fields")
	}

	// Field 1 (0-indexed): RSS (resident pages)
	rssPages, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS from /proc/pid/statm: %w", err)
	}

	// Convert pages to bytes using actual system page size
	systemPageSize := int64(os.Getpagesize())
	memoryBytes := uint64(rssPages * systemPageSize)

	return &types.Metrics{Mem: memoryBytes}, nil
}
