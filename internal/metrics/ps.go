package metrics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kitproj/kit/internal/types"
)

func GetPSCommand(pid int) []string {
	return []string{"ps", "-o", "rss", "-p", strconv.Itoa(pid)}
}

func ParsePSOutput(output string) (*types.Metrics, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected ps output: %s", output)
	}

	var totalMemory uint64

	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) < 1 {
			continue
		}

		rssMemoryKB, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			continue
		}

		totalMemory += uint64(rssMemoryKB * 1024)
	}

	return &types.Metrics{
		Mem: totalMemory,
	}, nil
}
