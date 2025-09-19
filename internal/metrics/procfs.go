package metrics

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/kitproj/kit/internal/types"
)

var (
	pageSize     int64
	pageSizeOnce sync.Once
)

// getSystemPageSize returns the system page size in bytes
func getSystemPageSize() int64 {
	pageSizeOnce.Do(func() {
		pageSize = getPageSizeFromSystem()
	})
	return pageSize
}

// getPageSizeFromSystem determines the system page size using multiple methods
func getPageSizeFromSystem() int64 {
	// Method 1: Try using getconf command
	if cmd := exec.Command("getconf", "PAGESIZE"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			if size, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64); err == nil && size > 0 {
				return size
			}
		}
	}

	// Method 2: Try reading /proc/meminfo (some systems show page size info)
	if content, err := os.ReadFile("/proc/meminfo"); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Hugepagesize:") {
				// This won't give us regular page size, but let's check other approaches
				continue
			}
		}
	}

	// Method 3: Read from /sys/kernel/mm/transparent_hugepage/hpage_pmd_size (if available)
	// This is for huge pages, not regular pages

	// Method 4: Use common page sizes based on architecture detection
	// Most x86/x86_64 systems use 4KB pages
	// Some ARM systems use 16KB or 64KB pages
	// For now, we'll try to detect through /proc/cpuinfo
	if content, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		cpuInfo := string(content)
		if strings.Contains(cpuInfo, "aarch64") || strings.Contains(cpuInfo, "arm64") {
			// ARM64 systems often use larger pages, but 4KB is still common
			// Without more specific detection, we'll fall back to 4KB
		}
	}

	// Fallback: Use standard 4KB page size (most common)
	return 4096
}

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
	systemPageSize := getSystemPageSize()
	memoryBytes := uint64(rssPages * systemPageSize)

	return &types.Metrics{Mem: memoryBytes}, nil
}
