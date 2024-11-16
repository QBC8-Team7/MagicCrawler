package logger

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/net"
)

// LogMemoryUsage logs current memory usage
func (l *AppLogger) LogMemoryUsage() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	l.LogSystemInfof("Memory usage",
		"Alloc", memStats.Alloc,
		"TotalAlloc", memStats.TotalAlloc,
		"Sys", memStats.Sys,
		"NumGC", memStats.NumGC,
	)
}

// LogCPUUsage logs CPU usage percentage over a short interval
func (l *AppLogger) logCPUusage() {
	percentages, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		l.LogSystemInfof("Failed to get CPU usage", "error", err)
		return
	}

	l.LogSystemInfof("CPU usage",
		"UsagePercent", percentages[0],
	)
}

// LogNetworkUsage logs network I/O stats
func (l *AppLogger) LogNetworkUsage() {
	netStats, err := net.IOCounters(false)
	if err != nil {
		l.LogSystemInfof("Failed to get network stats", "error", err)
		return
	}

	for _, stat := range netStats {
		l.LogSystemInfof("Network usage",
			"Interface", stat.Name,
			"BytesSent", stat.BytesSent,
			"BytesRecv", stat.BytesRecv,
			"PacketsSent", stat.PacketsSent,
			"PacketsRecv", stat.PacketsRecv,
		)
	}
}

func (l *AppLogger) StartSystemMetricsLogging() {
	go func() {
		for {
			l.LogMemoryUsage()
			l.logCPUusage()
			l.LogNetworkUsage()
			time.Sleep(10 * time.Second)
		}
	}()
}
