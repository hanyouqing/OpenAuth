package utils

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var startTime = time.Now()

type SystemResources struct {
	CPU     CPUInfo     `json:"cpu"`
	Memory  MemoryInfo  `json:"memory"`
	Disk    DiskInfo    `json:"disk"`
	Network NetworkInfo `json:"network"`
	Uptime  int64       `json:"uptime"` // seconds
}

type CPUInfo struct {
	Count        int     `json:"count"`
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`        // bytes
	Used        uint64  `json:"used"`         // bytes
	Available   uint64  `json:"available"`    // bytes
	UsagePercent float64 `json:"usage_percent"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`        // bytes
	Used        uint64  `json:"used"`         // bytes
	Free        uint64  `json:"free"`         // bytes
	UsagePercent float64 `json:"usage_percent"`
}

type NetworkInfo struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

func GetSystemResources() (*SystemResources, error) {
	resources := &SystemResources{
		Uptime: int64(time.Since(startTime).Seconds()),
	}

	// CPU info
	cpuCount := runtime.NumCPU()
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		// Fallback: set CPU count only
		resources.CPU = CPUInfo{
			Count:        cpuCount,
			UsagePercent: 0,
		}
	} else {
		resources.CPU = CPUInfo{
			Count:        cpuCount,
			UsagePercent: 0,
		}
		if len(cpuPercent) > 0 {
			resources.CPU.UsagePercent = cpuPercent[0]
		}
	}

	// Memory info
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		// Fallback: set defaults
		resources.Memory = MemoryInfo{
			Total:        0,
			Used:         0,
			Available:    0,
			UsagePercent: 0,
		}
	} else {
		resources.Memory = MemoryInfo{
			Total:        memInfo.Total,
			Used:         memInfo.Used,
			Available:    memInfo.Available,
			UsagePercent: memInfo.UsedPercent,
		}
	}

	// Disk info (root partition)
	diskInfo, err := disk.Usage("/")
	if err != nil {
		// If we can't get disk info, set defaults
		resources.Disk = DiskInfo{
			Total:        0,
			Used:         0,
			Free:         0,
			UsagePercent: 0,
		}
	} else {
		resources.Disk = DiskInfo{
			Total:        diskInfo.Total,
			Used:         diskInfo.Used,
			Free:         diskInfo.Free,
			UsagePercent: diskInfo.UsedPercent,
		}
	}

	// Network info
	netIO, err := net.IOCounters(false)
	if err != nil {
		// Fallback: set defaults
		resources.Network = NetworkInfo{
			BytesSent:   0,
			BytesRecv:   0,
			PacketsSent: 0,
			PacketsRecv: 0,
		}
	} else {
		if len(netIO) > 0 {
			resources.Network = NetworkInfo{
				BytesSent:   netIO[0].BytesSent,
				BytesRecv:   netIO[0].BytesRecv,
				PacketsSent: netIO[0].PacketsSent,
				PacketsRecv: netIO[0].PacketsRecv,
			}
		} else {
			resources.Network = NetworkInfo{
				BytesSent:   0,
				BytesRecv:   0,
				PacketsSent: 0,
				PacketsRecv: 0,
			}
		}
	}

	return resources, nil
}
