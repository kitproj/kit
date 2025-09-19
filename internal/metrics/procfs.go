package metrics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kitproj/kit/internal/types"
)

// GetProcFSCommand returns the cat command to read /proc/pid/stat
func GetProcFSCommand(pid int) []string {
	return []string{"cat", fmt.Sprintf("/proc/%d/stat", pid)}
}

// ParseProcFSOutput parses /proc/pid/stat output for memory usage only
func ParseProcFSOutput(output string) (*types.Metrics, error) {
	fields := strings.Fields(strings.TrimSpace(output))
	if len(fields) < 24 {
		return nil, nil, fmt.Errorf("unexpected /proc/pid/stat output: insufficient fields")
	}

	// Field 23 (0-indexed): RSS in pages
	rssPages, err := strconv.ParseInt(fields[23], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse RSS from /proc/pid/stat: %w", err)
	}

	// Convert pages to bytes (assuming 4KB page size)
	memoryBytes := uint64(rssPages * 4096)

	return &types.Metrics{
		Mem: memoryBytes,
	}, currentSnapshot, nil
}
