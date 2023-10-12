package uptime

import (
	"time"

	"github.com/shirou/gopsutil/host"
)

func GetUptime() (time.Duration, error) {
	uptime, err := host.Uptime()
	if err != nil {
		return 0, err
	}
	return time.Duration(uptime) * time.Second, nil
}
