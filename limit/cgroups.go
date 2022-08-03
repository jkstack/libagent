//go:build !windows && !aix
// +build !windows,!aix

package limit

import (
	"os"
	"strconv"
	"strings"

	"github.com/containerd/cgroups"
	cgroupsv2 "github.com/containerd/cgroups/v2"
	"github.com/dustin/go-humanize"
	"github.com/jkstack/jkframe/logging"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// Do set cgroups limit
func (cfg *Configure) Do(agentName string) {
	if cgroups.Mode() == cgroups.Unified {
		cfg.doV2(agentName)
		return
	}
	cfg.doV1(agentName)
}

func (cfg *Configure) doV2(agentName string) {
	// TODO
	logging.Info("use cgroups_v2")
	cgroupsv2.NewSystemd("/", "my-group.slice", -1, &cgroupsv2.Resources{})
}

func (cfg *Configure) doV1(agentName string) {
	if !wantCGroup(cfg) {
		return
	}
	logging.Info("use cgroups_v1")
	dir := "/jkstack/agent/" + agentName
	group, err := cgroups.New(cgroups.V1, cgroups.StaticPath(dir), &specs.LinuxResources{})
	if err != nil {
		logging.Warning("can not create cgroup %s: %v", dir, err)
		return
	}
	limitCPU(group, cfg.CpuQuota)
	limitMemory(group, int64(cfg.Memory))
	limitDisk(group, cfg.Disks)
	err = group.Add(cgroups.Process{
		Pid: os.Getpid(),
	})
	if err != nil {
		logging.Warning("can not add current pid: %v", err)
		return
	}
}

func wantCGroup(cfg *Configure) bool {
	if cfg.CpuQuota > 0 {
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

func limitCPU(group cgroups.Cgroup, limit int64) {
	cpu := limit * 1000
	err := group.Update(&specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Quota: &cpu,
		},
	})
	if err != nil {
		logging.Warning("can not set cpu_quota to %d%%: %v", limit, err)
		return
	}
	logging.Info("set cpu_quota to %d%%", limit)
}

func limitMemory(group cgroups.Cgroup, limit int64) {
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

func limitDisk(group cgroups.Cgroup, limits diskLimits) {
	var block specs.LinuxBlockIO
	for _, disk := range limits {
		write := func(value uint64, target []specs.LinuxThrottleDevice) []specs.LinuxThrottleDevice {
			var device specs.LinuxThrottleDevice
			device.Major, device.Minor = parseDev(disk.Dev)
			device.Rate = value
			return append(target, device)
		}
		if disk.ReadBytes > 0 {
			block.ThrottleReadBpsDevice = write(disk.ReadBytes, block.ThrottleReadBpsDevice)
			logging.Info("  - set read_bytes limit by mount_point %s: %s",
				disk.Dev, humanize.IBytes(disk.ReadBytes))
		}
		if disk.WriteBytes > 0 {
			block.ThrottleWriteBpsDevice = write(disk.WriteBytes, block.ThrottleWriteBpsDevice)
			logging.Info("  - set write_bytes limit by mount_point %s: %s",
				disk.Dev, humanize.IBytes(disk.WriteBytes))
		}
		if disk.ReadIOPS > 0 {
			block.ThrottleReadIOPSDevice = write(disk.ReadIOPS, block.ThrottleReadIOPSDevice)
			logging.Info("  - set read_iops limit by mount_point %s: %s",
				disk.Dev, humanize.IBytes(disk.ReadIOPS))
		}
		if disk.WriteIOPS > 0 {
			block.ThrottleWriteIOPSDevice = write(disk.WriteIOPS, block.ThrottleWriteIOPSDevice)
			logging.Info("  - set write_iops limit by mount_point %s: %s",
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
