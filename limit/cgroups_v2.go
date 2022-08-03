package limit

import (
	"os"

	v2 "github.com/containerd/cgroups/v2"
	"github.com/dustin/go-humanize"
	"github.com/jkstack/jkframe/logging"
)

func (cfg *Configure) doV2(agentName string) {
	if !wantCGroup(cfg) {
		return
	}
	logging.Info("use cgroups_v2")
	group, err := v2.NewSystemd("/", agentName+".slice", -1, &v2.Resources{})
	if err != nil {
		logging.Warning("can not create cgroup: %v", err)
		return
	}
	limitCPUV2(group, cfg.CpuQuota)
	limitMemoryV2(group, int64(cfg.Memory))
	limitDiskV2(group, cfg.Disks)
	err = group.AddProc(uint64(os.Getpid()))
	if err != nil {
		logging.Warning("can not add current pid: %v", err)
		return
	}
}

func limitCPUV2(group *v2.Manager, limit int64) {
	quota := limit * 1000
	period := uint64(1000000)
	err := group.Update(&v2.Resources{
		CPU: &v2.CPU{
			Max: v2.NewCPUMax(&quota, &period),
		},
	})
	if err != nil {
		logging.Warning("can not set cpu_quota to %d%%: %v", limit, err)
		return
	}
	logging.Info("set cpu_quota to %d%%", limit)
}

func limitMemoryV2(group *v2.Manager, limit int64) {
	err := group.Update(&v2.Resources{
		Memory: &v2.Memory{
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

func limitDiskV2(group *v2.Manager, limits diskLimits) {
	var io v2.IO
	for _, disk := range limits {
		write := func(t v2.IOType, value uint64) {
			major, minor := parseDev(disk.Dev)
			io.Max = append(io.Max, v2.Entry{
				Type:  t,
				Major: major,
				Minor: minor,
				Rate:  value,
			})
		}
		if disk.ReadBytes > 0 {
			write(v2.ReadBPS, disk.ReadBytes)
			logging.Info("  - set read_bytes limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.ReadBytes))
		}
		if disk.WriteBytes > 0 {
			write(v2.WriteBPS, disk.WriteBytes)
			logging.Info("  - set write_bytes limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.WriteBytes))
		}
		if disk.ReadIOPS > 0 {
			write(v2.ReadIOPS, disk.ReadIOPS)
			logging.Info("  - set read_iops limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.ReadIOPS))
		}
		if disk.WriteIOPS > 0 {
			write(v2.WriteIOPS, disk.WriteIOPS)
			logging.Info("  - set write_iops limit by dev [%s]: %s",
				disk.Dev, humanize.IBytes(disk.WriteIOPS))
		}
	}
	err := group.Update(&v2.Resources{
		IO: &io,
	})
	if err != nil {
		logging.Warning("can not set disk_limit: %v", err)
		return
	}
	logging.Info("set disk_limit successed")
}
