//go:build !windows && !aix && !darwin
// +build !windows,!aix,!darwin

package limit

import (
	"os"
	"strconv"
	"strings"

	"github.com/containerd/cgroups/v3/cgroup1"
	"github.com/dustin/go-humanize"
	"github.com/jkstack/jkframe/logging"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func (cfg *Configure) doV1(agentName string) {
	if !wantCGroup(cfg) {
		return
	}
	logging.Info("use cgroups_v1")
	dir := "/jkstack/agent/" + agentName
	group, err := cgroup1.New(cgroup1.StaticPath(dir), &specs.LinuxResources{})
	if err != nil {
		logging.Warning("can not create cgroup %s: %v", dir, err)
		return
	}
	limitCPUV1(group, cfg.CPUQuota)
	limitMemoryV1(group, int64(cfg.Memory))
	limitDiskV1(group, cfg.Disks)
	err = group.Add(cgroup1.Process{
		Pid: os.Getpid(),
	})
	if err != nil {
		logging.Warning("can not add current pid: %v", err)
		return
	}
}

func wantCGroup(cfg *Configure) bool {
	if cfg.CPUQuota > 0 {
		return true
	}
	if cfg.Memory > 0 {
		return true
	}
	if len(cfg.Disks) > 0 {
		return true
	}
	return false
}

func limitCPUV1(group cgroup1.Cgroup, limit int64) {
	quota := limit * 1000
	err := group.Update(&specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Quota: &quota,
		},
	})
	if err != nil {
		logging.Warning("can not set cpu_quota to %d%%: %v", limit, err)
		return
	}
	logging.Info("set cpu_quota to %d%%", limit)
}

func limitMemoryV1(group cgroup1.Cgroup, limit int64) {
	err := group.Update(&specs.LinuxResources{
		Memory: &specs.LinuxMemory{
			Limit: &limit,
		},
	})
	if err != nil {
		logging.Warning("can not set memory_limit to %s: %v",
			humanize.IBytes(uint64(limit)), err)
		return
	}
	logging.Info("set memory_limit to %s", humanize.IBytes(uint64(limit)))
}

func limitDiskV1(group cgroup1.Cgroup, limits diskLimits) {
	var block specs.LinuxBlockIO
	for _, disk := range limits {
		write := func(value uint64, target []specs.LinuxThrottleDevice) []specs.LinuxThrottleDevice {
			var device specs.LinuxThrottleDevice
			device.Major, device.Minor = parseDev(disk.Dev)
			device.Rate = value
			return append(target, device)
		}
		if disk.ReadBytes > 0 {
			block.ThrottleReadBpsDevice = write(disk.ReadBytes.Bytes(), block.ThrottleReadBpsDevice)
			logging.Info("  - set read_bytes limit by dev [%s]: %s",
				disk.Dev, disk.ReadBytes.String())
		}
		if disk.WriteBytes > 0 {
			block.ThrottleWriteBpsDevice = write(disk.WriteBytes.Bytes(), block.ThrottleWriteBpsDevice)
			logging.Info("  - set write_bytes limit by dev [%s]: %s",
				disk.Dev, disk.WriteBytes.String())
		}
		if disk.ReadIOPS > 0 {
			block.ThrottleReadIOPSDevice = write(disk.ReadIOPS, block.ThrottleReadIOPSDevice)
			logging.Info("  - set read_iops limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.ReadIOPS))
		}
		if disk.WriteIOPS > 0 {
			block.ThrottleWriteIOPSDevice = write(disk.WriteIOPS, block.ThrottleWriteIOPSDevice)
			logging.Info("  - set write_iops limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.WriteIOPS))
		}
	}
	err := group.Update(&specs.LinuxResources{
		BlockIO: &block,
	})
	if err != nil {
		logging.Warning("can not set disk_limit: %v", err)
		return
	}
	logging.Info("set disk_limit successed")
}

func parseDev(str string) (int64, int64) {
	tmp := strings.SplitN(str, ":", 2)
	major, _ := strconv.ParseInt(tmp[0], 10, 64)
	minor, _ := strconv.ParseInt(tmp[1], 10, 64)
	return major, minor
}
