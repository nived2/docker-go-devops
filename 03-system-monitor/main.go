package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type SystemStats struct {
	Timestamp     string         `json:"timestamp"`
	HostInfo      HostInfo       `json:"host_info"`
	CPU           CPUStats       `json:"cpu"`
	Memory        MemoryStats    `json:"memory"`
	Disk          []DiskStats    `json:"disk"`
	ProcessCount  int            `json:"process_count"`
	TopProcesses  []ProcessStats `json:"top_processes"`
}

type HostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	Uptime         uint64 `json:"uptime"`
}

type CPUStats struct {
	Usage       float64   `json:"usage"`
	CoreCount   int       `json:"core_count"`
	LoadAverage []float64 `json:"load_average"`
	PerCPU      []float64 `json:"per_cpu"`
}

type MemoryStats struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePerc    float64 `json:"usage_percentage"`
	SwapTotal    uint64  `json:"swap_total"`
	SwapUsed     uint64  `json:"swap_used"`
	SwapFree     uint64  `json:"swap_free"`
	SwapUsagePerc float64 `json:"swap_usage_percentage"`
}

type DiskStats struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsagePerc   float64 `json:"usage_percentage"`
	InodesTotal uint64  `json:"inodes_total"`
	InodesUsed  uint64  `json:"inodes_used"`
	InodesFree  uint64  `json:"inodes_free"`
}

type ProcessStats struct {
	PID         int32   `json:"pid"`
	PPID        int32   `json:"ppid"`
	Name        string  `json:"name"`
	Username    string  `json:"username"`
	CPUPercent  float64 `json:"cpu_percent"`
	MemoryPerc  float32 `json:"memory_percentage"`
	MemoryUsage uint64  `json:"memory_usage"`
	Status      string  `json:"status"`
	CreateTime  int64   `json:"create_time"`
}

func init() {
	// Set host paths for containerized environment
	if hostProc := os.Getenv("HOST_PROC"); hostProc != "" {
		if err := os.Setenv("HOST_PROC", hostProc); err != nil {
			log.Printf("Warning: Failed to set HOST_PROC: %v", err)
		}
	}
	if hostSys := os.Getenv("HOST_SYS"); hostSys != "" {
		if err := os.Setenv("HOST_SYS", hostSys); err != nil {
			log.Printf("Warning: Failed to set HOST_SYS: %v", err)
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting System Monitor on port %s", port)
	log.Printf("Running with CPU cores: %d", runtime.NumCPU())

	// Create routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/processes", handleProcesses)
	http.HandleFunc("/health", handleHealth)

	log.Printf("Server is ready to handle requests at :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "System Monitor is running. Available endpoints:\n"+
		"- /metrics - System metrics\n"+
		"- /processes - Process information\n"+
		"- /health - Health check")
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	stats, err := getSystemStats()
	if err != nil {
		log.Printf("Error getting system stats: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleProcesses(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Default limit
	processes, err := getTopProcesses(limit)
	if err != nil {
		log.Printf("Error getting processes: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(processes); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "UP",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func getSystemStats() (*SystemStats, error) {
	stats := &SystemStats{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Host Info
	hostInfo, err := host.Info()
	if err != nil {
		log.Printf("Warning: Could not get host info: %v", err)
	} else {
		stats.HostInfo = HostInfo{
			Hostname:        hostInfo.Hostname,
			OS:             hostInfo.OS,
			Platform:       hostInfo.Platform,
			PlatformVersion: hostInfo.PlatformVersion,
			KernelVersion:  hostInfo.KernelVersion,
			Uptime:        hostInfo.Uptime,
		}
	}

	// CPU Stats
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("error getting CPU stats: %v", err)
	}

	perCPU, err := cpu.Percent(0, true)
	if err != nil {
		log.Printf("Warning: Could not get per-CPU stats: %v", err)
	}

	loadAvg, err := load.Avg()
	if err != nil {
		log.Printf("Warning: Could not get load average: %v", err)
	}

	stats.CPU = CPUStats{
		Usage:     cpuPercent[0],
		CoreCount: runtime.NumCPU(),
		PerCPU:    perCPU,
	}

	if loadAvg != nil {
		stats.CPU.LoadAverage = []float64{loadAvg.Load1, loadAvg.Load5, loadAvg.Load15}
	}

	// Memory Stats
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting memory stats: %v", err)
	}

	swapMem, err := mem.SwapMemory()
	if err != nil {
		log.Printf("Warning: Could not get swap memory stats: %v", err)
	}

	stats.Memory = MemoryStats{
		Total:     virtualMem.Total,
		Used:      virtualMem.Used,
		Free:      virtualMem.Free,
		UsagePerc: virtualMem.UsedPercent,
	}

	if swapMem != nil {
		stats.Memory.SwapTotal = swapMem.Total
		stats.Memory.SwapUsed = swapMem.Used
		stats.Memory.SwapFree = swapMem.Free
		stats.Memory.SwapUsagePerc = swapMem.UsedPercent
	}

	// Disk Stats
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("error getting disk partitions: %v", err)
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			log.Printf("Warning: Could not get disk usage for %s: %v", partition.Mountpoint, err)
			continue
		}
		stats.Disk = append(stats.Disk, DiskStats{
			Device:      partition.Device,
			MountPoint:  partition.Mountpoint,
			FSType:      partition.Fstype,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsagePerc:   usage.UsedPercent,
			InodesTotal: usage.InodesTotal,
			InodesUsed:  usage.InodesUsed,
			InodesFree:  usage.InodesFree,
		})
	}

	// Process Count and Top Processes
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error getting processes: %v", err)
	}
	stats.ProcessCount = len(processes)

	topProcesses, err := getTopProcesses(5)
	if err == nil {
		stats.TopProcesses = topProcesses
	} else {
		log.Printf("Warning: Could not get top processes: %v", err)
	}

	return stats, nil
}

func getTopProcesses(limit int) ([]ProcessStats, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var processStats []ProcessStats
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			continue
		}

		mem, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		memPercent, err := p.MemoryPercent()
		if err != nil {
			log.Printf("Warning: Could not get memory percent for process %s: %v", name, err)
		}

		username, err := p.Username()
		if err != nil {
			log.Printf("Warning: Could not get username for process %s: %v", name, err)
		}

		status, err := p.Status()
		if err != nil {
			log.Printf("Warning: Could not get status for process %s: %v", name, err)
		}

		createTime, err := p.CreateTime()
		if err != nil {
			log.Printf("Warning: Could not get create time for process %s: %v", name, err)
		}

		ppid, err := p.Ppid()
		if err != nil {
			log.Printf("Warning: Could not get parent PID for process %s: %v", name, err)
		}

		processStats = append(processStats, ProcessStats{
			PID:         p.Pid,
			PPID:        ppid,
			Name:        name,
			Username:    username,
			CPUPercent:  cpu,
			MemoryPerc:  memPercent,
			MemoryUsage: mem.RSS,
			Status:      strings.Join(status, ", "),
			CreateTime:  createTime,
		})
	}

	// Sort by CPU usage (simple bubble sort for demonstration)
	for i := 0; i < len(processStats)-1; i++ {
		for j := 0; j < len(processStats)-i-1; j++ {
			if processStats[j].CPUPercent < processStats[j+1].CPUPercent {
				processStats[j], processStats[j+1] = processStats[j+1], processStats[j]
			}
		}
	}

	if len(processStats) > limit {
		processStats = processStats[:limit]
	}

	return processStats, nil
}
