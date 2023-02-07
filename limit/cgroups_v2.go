//go:build !windows && !aix && !darwin
// +build !windows,!aix,!darwin

package limit

import (
	"os"

	"github.com/containerd/cgroups/v3/cgroup2"
	"github.com/dustin/go-humanize"
	"github.com/jkstack/jkframe/logging"
)

func (cfg *Configure) doV2(agentName string) {
	if !wantCGroup(cfg) {
		return
	}
	logging.Info("use cgroups_v2")
	group, err := cgroup2.NewSystemd("/", agentName+".slice", -1, &cgroup2.Resources{})
	if err != nil {
		logging.Warning("can not create cgroup: %v", err)
		return
	}
	limitCPUV2(group, cfg.CPUQuota)
	limitMemoryV2(group, int64(cfg.Memory))
	limitDiskV2(group, cfg.Disks)
	err = group.AddProc(uint64(os.Getpid()))
	if err != nil {
		logging.Warning("can not add current pid: %v", err)
		return
	}
}

func limitCPUV2(group *cgroup2.Manager, limit int64) {
	quota := limit * 1000
	period := uint64(1000000)
	err := group.Update(&cgroup2.Resources{
		CPU: &cgroup2.CPU{
			Max: cgroup2.NewCPUMax(&quota, &period),
		},
	})
	if err != nil {
		logging.Warning("can not set cpu_quota to %d%%: %v", limit, err)
		return
	}
	logging.Info("set cpu_quota to %d%%", limit)
}

func limitMemoryV2(group *cgroup2.Manager, limit int64) {
	err := group.Update(&cgroup2.Resources{
		Memory: &cgroup2.Memory{
			Max: &limit,
		},
	})
	if err != nil {
		logging.Warning("can not set memory_limit to %s: %v",
			humanize.IBytes(uint64(limit)), err)
		return
	}
	logging.Info("set memory_limit to %s", humanize.IBytes(uint64(limit)))
}

func limitDiskV2(group *cgroup2.Manager, limits diskLimits) {
	var io cgroup2.IO
	for _, disk := range limits {
		write := func(t cgroup2.IOType, value uint64) {
			major, minor := parseDev(disk.Dev)
			io.Max = append(io.Max, cgroup2.Entry{
				Type:  t,
				Major: major,
				Minor: minor,
				Rate:  value,
			})
		}
		if disk.ReadBytes > 0 {
			write(cgroup2.ReadBPS, disk.ReadBytes.Bytes())
			logging.Info("  - set read_bytes limit by dev [%s]: %s",
				disk.Dev, disk.ReadBytes.String())
		}
		if disk.WriteBytes > 0 {
			write(cgroup2.WriteBPS, disk.WriteBytes.Bytes())
			logging.Info("  - set write_bytes limit by dev [%s]: %s",
				disk.Dev, disk.WriteBytes.Bytes())
		}
		if disk.ReadIOPS > 0 {
			write(cgroup2.ReadIOPS, disk.ReadIOPS)
			logging.Info("  - set read_iops limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.ReadIOPS))
		}
		if disk.WriteIOPS > 0 {
			write(cgroup2.WriteIOPS, disk.WriteIOPS)
			logging.Info("  - set write_iops limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.WriteIOPS))
		}
	}
	err := group.Update(&cgroup2.Resources{
		IO: &io,
	})
	if err != nil {
		logging.Warning("can not set disk_limit: %v", err)
		return
	}
	logging.Info("set disk_limit successed")
}
