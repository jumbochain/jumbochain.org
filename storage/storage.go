package storage

import (
	"github.com/shirou/gopsutil/disk"
)

type HardDisk struct {
	Name      []string
	SizeBytes float64
	FreeBytes float64
}

func GetHardDisks() []HardDisk {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil
	}
	var disks []HardDisk
	var totalUsage, totalFree float64
	var deviceName []string
	for _, partition := range partitions {
		if partition.Fstype == "ext4" || partition.Fstype == "ext3" || partition.Fstype == "xfs" || partition.Fstype == "NTFS" {
			usage, err := disk.Usage(partition.Mountpoint)
			if err != nil {
				return nil
			}
			totalUsage += float64(usage.Total)
			totalFree += float64(usage.Free)
			deviceName = append(deviceName, partition.Device)
		}
	}
	disks = append(disks, HardDisk{
		Name:      deviceName,
		SizeBytes: totalUsage,
		FreeBytes: totalFree,
	})

	return disks
}
