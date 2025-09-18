package metrics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kitproj/kit/internal/types"
)

func GetCommand(pid int) []string {
	return []string{"ps", "-o", "%cpu,rss", "-p", strconv.Itoa(pid)}
}

func ParseOutput(output string) (*types.Metrics, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected ps output: %s", output)
	}

	var totalCPU float64
	var totalMemory uint64

	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) < 2 {
			continue
		}

		cpuPercentage, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			continue
		}

		rssMemoryKB, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}

		totalCPU += cpuPercentage
		totalMemory += uint64(rssMemoryKB * 1024)
	}

	cpuMillicores := totalCPU * 10

	return &types.Metrics{
		CPU: uint64(cpuMillicores),
		Mem: totalMemory,
	}, nil
}