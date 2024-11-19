package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type Usage struct {
	CPU_Before float64
	CPU_After  float64
	RAM        float64
	Time       time.Duration
}

func RunAndMeasureUsage[T any](mainLogger *logger.AppLogger, f func() T) (T, Usage) {
	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	startTime := time.Now()

	cpuBefore, err := getCPUUsage()
	if err != nil {
		mainLogger.Error("error measuring CPU usage:", err)
	}
	result := f()

	cpuAfter, err := getCPUUsage()
	if err != nil {
		mainLogger.Error("error measuring CPU usage:", err)
	}

	duration := time.Since(startTime)

	runtime.ReadMemStats(&memAfter)
	ramUsage := float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / 1024 / 1024

	return result, Usage{
		CPU_Before: cpuBefore,
		CPU_After:  cpuAfter,
		RAM:        ramUsage,
		Time:       duration,
	}
}

func getCPUUsage() (float64, error) {
	cmd := exec.Command("top", "-bn1")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return 0, err
	}

	output := out.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Cpu(s)") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "id," || part == "id" {
					idle, err := parseFloat(parts[i-1])
					if err != nil {
						return 0, err
					}
					return 100 - idle, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("CPU usage not found")
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
}
