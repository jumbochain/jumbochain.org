package memory

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetTotalRAM() float64 {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("Error getting total RAM: %v\n", err)
		return 0
	}
	return float64(vmStat.Total) / 1024 / 1024 / 1024
}

func GetCPUCoreCount() uint64 {
	cpuStat, err := cpu.Info()
	if err != nil {
		return 0
	}
	var cpu_core uint64
	for i := 0; i < len(cpuStat); i++ {
		cpu_core += uint64(cpuStat[i].Cores)
	}
	return cpu_core
}
