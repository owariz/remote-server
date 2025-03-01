package handlers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// ServerStatus โครงสร้างข้อมูลสถานะเซิร์ฟเวอร์
type ServerStatus struct {
	Hostname    string    `json:"hostname"`
	OS          string    `json:"os"`
	Arch        string    `json:"architecture"`   // เช่น "amd64"
	Uptime      uint64    `json:"uptime_seconds"` // วินาที
	UptimeHuman string    `json:"uptime_human"`   // เช่น "2 days, 3 hours"
	LoadAverage []float64 `json:"load_average"`   // Load 1, 5, 15 นาที
	CPU         CPUInfo   `json:"cpu"`            // ข้อมูล CPU
	Memory      Memory    `json:"memory"`         // ข้อมูลหน่วยความจำ
	Time        time.Time `json:"time"`           // เวลาปัจจุบัน
}

type Memory struct {
	TotalGB   float64 `json:"total_gb"`      // GB
	UsedGB    float64 `json:"used_gb"`       // GB
	FreeGB    float64 `json:"free_gb"`       // GB
	Total     uint64  `json:"total_bytes"`   // bytes
	Used      uint64  `json:"used_bytes"`    // bytes
	Free      uint64  `json:"free_bytes"`    // bytes
	UsagePerc float64 `json:"usage_percent"` // เปอร์เซ็นต์
}

type CPUInfo struct {
	Cores     int     `json:"cores"`         // จำนวน cores
	UsagePerc float64 `json:"usage_percent"` // เปอร์เซ็นต์การใช้งาน
}

func GetStatus(c *fiber.Ctx) error {
	// ใช้ WaitGroup เพื่อรอข้อมูลจาก goroutines
	var wg sync.WaitGroup
	var mu sync.Mutex // ป้องกัน race condition

	// ตัวแปรสำหรับเก็บข้อมูล
	var hostname string
	var osDetail string
	var arch string
	var uptime uint64
	var loadAvg []float64
	var memInfo Memory
	var cpuInfo CPUInfo
	var err error

	// ดึง hostname
	wg.Add(1)
	go func() {
		defer wg.Done()
		hostname, err = os.Hostname()
		if err != nil {
			hostname = "unknown"
		}
	}()

	// ดึงข้อมูล OS, architecture, และ uptime จาก host.Info()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if info, err := host.Info(); err == nil {
			mu.Lock()
			osDetail = info.OS + " " + info.PlatformVersion
			arch = info.KernelArch
			uptime = info.Uptime
			mu.Unlock()
		} else {
			mu.Lock()
			osDetail = runtime.GOOS
			arch = runtime.GOARCH
			uptime = 0
			mu.Unlock()
		}
	}()

	// ดึง load average
	wg.Add(1)
	go func() {
		defer wg.Done()
		if avg, err := load.Avg(); err == nil && runtime.GOOS != "windows" {
			mu.Lock()
			loadAvg = []float64{avg.Load1, avg.Load5, avg.Load15}
			mu.Unlock()
		} else {
			mu.Lock()
			loadAvg = []float64{0.0, 0.0, 0.0}
			mu.Unlock()
		}
	}()

	// ดึงข้อมูล memory
	wg.Add(1)
	go func() {
		defer wg.Done()
		if v, err := mem.VirtualMemory(); err == nil {
			mu.Lock()
			memInfo = Memory{
				TotalGB:   float64(v.Total) / (1024 * 1024 * 1024),
				UsedGB:    float64(v.Used) / (1024 * 1024 * 1024),
				FreeGB:    float64(v.Free) / (1024 * 1024 * 1024),
				Total:     v.Total,
				Used:      v.Used,
				Free:      v.Free,
				UsagePerc: v.UsedPercent,
			}
			mu.Unlock()
		} else {
			mu.Lock()
			memInfo = Memory{
				TotalGB:   0,
				UsedGB:    0,
				FreeGB:    0,
				Total:     0,
				Used:      0,
				Free:      0,
				UsagePerc: 0.0,
			}
			mu.Unlock()
		}
	}()

	// ดึงข้อมูล CPU
	wg.Add(1)
	go func() {
		defer wg.Done()
		if counts, err := cpu.Counts(true); err == nil {
			if percent, err := cpu.Percent(1*time.Second, false); err == nil && len(percent) > 0 {
				mu.Lock()
				cpuInfo = CPUInfo{
					Cores:     counts,
					UsagePerc: percent[0], // เปอร์เซ็นต์รวมของทุก core
				}
				mu.Unlock()
			}
		}
	}()

	// รอให้ทุก goroutine เสร็จสิ้น
	wg.Wait()

	// แปลง uptime เป็น human-readable
	uptimeHuman := formatUptime(uptime)

	// สร้าง status
	status := ServerStatus{
		Hostname:    hostname,
		OS:          osDetail,
		Arch:        arch,
		Uptime:      uptime,
		UptimeHuman: uptimeHuman,
		LoadAverage: loadAvg,
		CPU:         cpuInfo,
		Memory:      memInfo,
		Time:        time.Now(),
	}

	return c.JSON(status)
}

func formatUptime(uptime uint64) string {
	days := uptime / (24 * 3600)
	hours := (uptime % (24 * 3600)) / 3600
	minutes := (uptime % 3600) / 60
	seconds := uptime % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", seconds))
	}
	return strings.Join(parts, ", ")
}

func GetMetrics(c *fiber.Ctx) error {
	metrics := fiber.Map{
		"cpu": fiber.Map{
			"usage": 35.5,
			"cores": runtime.NumCPU(),
		},
		"memory": fiber.Map{
			"usage":      50.0,
			"total":      8 * 1024 * 1024 * 1024,
			"used":       4 * 1024 * 1024 * 1024,
			"free":       4 * 1024 * 1024 * 1024,
			"swap_total": 2 * 1024 * 1024 * 1024,
			"swap_used":  512 * 1024 * 1024,
		},
		"disk": fiber.Map{
			"usage": 40.0,
			"total": 500 * 1024 * 1024 * 1024,
			"used":  200 * 1024 * 1024 * 1024,
			"free":  300 * 1024 * 1024 * 1024,
		},
		"network": fiber.Map{
			"in_traffic":  1024 * 1024, // 1MB
			"out_traffic": 512 * 1024,  // 512KB
			"in_speed":    100 * 1024,  // 100KB/s
			"out_speed":   50 * 1024,   // 50KB/s
			"connections": 42,
		},
	}

	return c.JSON(metrics)
}

func RestartService(c *fiber.Ctx) error {
	serviceName := c.Params("name")

	if strings.ContainsAny(serviceName, ";&|") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid service name")
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("net", "stop", serviceName, "&&", "net", "start", serviceName)
	} else {
		cmd = exec.Command("systemctl", "restart", serviceName)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to restart service: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Service restarted successfully",
		"output":  string(output),
	})
}

func GetServiceLogs(c *fiber.Ctx) error {
	serviceName := c.Params("name")

	if strings.ContainsAny(serviceName, ";&|") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid service name")
	}

	lines := c.Query("lines", "100")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command",
			"Get-EventLog -LogName Application -Source "+serviceName+" -Newest "+lines)
	} else {
		cmd = exec.Command("journalctl", "-u", serviceName, "-n", lines)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get logs: "+err.Error())
	}

	logs := strings.Split(string(output), "\n")

	return c.JSON(fiber.Map{
		"service": serviceName,
		"lines":   len(logs),
		"logs":    logs,
	})
}

var AllowedCommands = map[string]bool{
	"dir":  true, // Windows built-in
	"echo": true, // Windows/Unix built-in
	"ls":   true, // Unix
	"ping": true, // Windows/Unix executable
}

func ExecuteCommand(c *fiber.Ctx) error {
	type CommandRequest struct {
		Command string `json:"command" validate:"required"`
	}

	var req CommandRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid request body format",
			"detail": err.Error(),
		})
	}

	req.Command = strings.TrimSpace(req.Command)
	if req.Command == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Command cannot be empty",
		})
	}

	if len(req.Command) > 1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Command exceeds maximum length of 1024 characters",
		})
	}

	unsafePattern := regexp.MustCompile(`[;&|<>$\(\)\{\}\[\]!\\]`)
	if unsafePattern.MatchString(req.Command) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Command contains unsafe characters",
		})
	}

	parts := strings.Fields(req.Command)
	if len(parts) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid command format",
		})
	}

	if !AllowedCommands[parts[0]] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Command not allowed",
		})
	}

	if len(parts) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Too many arguments (max 10)",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd.exe", "/C", req.Command)
	} else {
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", req.Command)
	}
	cmd.Env = []string{}

	if runtime.GOOS == "windows" && parts[0] == "ping" {
		cmd = exec.CommandContext(ctx, "C:\\Windows\\System32\\ping.exe", parts[1:]...)
	}

	output, err := cmd.CombinedOutput()

	if len(output) > 1024*1024 {
		output = output[:1024*1024]
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"output":  string(output),
			"error":   "Output truncated: exceeds 1MB limit",
		})
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
				"success": false,
				"output":  string(output),
				"error":   "Command execution timeout after 5 seconds",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"output":  string(output),
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"output":  string(output),
	})
}
