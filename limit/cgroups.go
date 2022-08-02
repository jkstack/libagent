//go:build !windows && !aix
// +build !windows,!aix

package limit

import (
	"bufio"
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
	cf, want := buildLinuxResource(cfg)
	if !want {
		return
	}
	logging.Info("use cgroups_v1")
	dir := "/jkstack/agent/" + agentName
	group, err := cgroups.New(cgroups.V1, cgroups.StaticPath(dir), &cf)
	if err != nil {
		logging.Warning("can not create cgroup %s: %v", dir, err)
		return
	}
	err = group.Add(cgroups.Process{
		Pid: os.Getpid(),
	})
	if err != nil {
		logging.Warning("can not add current pid: %v", err)
		return
	}
}

func buildLinuxResource(cfg *Configure) (specs.LinuxResources, bool) {
	var ret specs.LinuxResources
	want := false
	if cfg.CpuQuota > 0 {
		cpu := int64(cfg.CpuQuota) * 1000
		ret.CPU = &specs.LinuxCPU{
			Quota: &cpu,
		}
		logging.Info("set cpu_quota to %d%%", cfg.CpuQuota)
		want = true
	}
	if cfg.Memory > 0 {
		mem := int64(cfg.Memory.Bytes())
		ret.Memory = &specs.LinuxMemory{
			Limit: &mem,
		}
		logging.Info("set memory_limit to %s", cfg.Memory.String())
		want = true
	}
	if len(cfg.Disks) > 0 {
		mnt := readMnt()
		var block specs.LinuxBlockIO
		for _, disk := range cfg.Disks {
			dev := mnt[disk.MountPoint]
			if len(dev) == 0 {
				logging.Warning("can not get dev of mount_point %s", disk.MountPoint)
				continue
			}
			write := func(value uint64, target []specs.LinuxThrottleDevice) []specs.LinuxThrottleDevice {
				var device specs.LinuxThrottleDevice
				device.Major, device.Minor = parseDev(dev)
				device.Rate = value
				return append(target, device)
			}
			if disk.ReadBytes > 0 {
				block.ThrottleReadBpsDevice = write(disk.ReadBytes, block.ThrottleReadBpsDevice)
				logging.Info("set read_bytes limit by mount_point %s(%s): %s",
					disk.MountPoint, dev, humanize.IBytes(disk.ReadBytes))
			}
			if disk.WriteBytes > 0 {
				block.ThrottleWriteBpsDevice = write(disk.WriteBytes, block.ThrottleWriteBpsDevice)
				logging.Info("set write_bytes limit by mount_point %s(%s): %s",
					disk.MountPoint, dev, humanize.IBytes(disk.WriteBytes))
			}
			if disk.ReadIOPS > 0 {
				block.ThrottleReadIOPSDevice = write(disk.ReadIOPS, block.ThrottleReadIOPSDevice)
				logging.Info("set read_iops limit by mount_point %s(%s): %s",
					disk.MountPoint, dev, humanize.IBytes(disk.ReadIOPS))
			}
			if disk.WriteIOPS > 0 {
				block.ThrottleWriteIOPSDevice = write(disk.WriteIOPS, block.ThrottleWriteIOPSDevice)
				logging.Info("set write_iops limit by mount_point %s(%s): %s",
					disk.MountPoint, dev, humanize.IBytes(disk.WriteIOPS))
			}
		}
		ret.BlockIO = &block
		want = true
	}
	return ret, want
}

// mount_point => dev
func readMnt() map[string]string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		logging.Warning("can not open /proc/self/mountinfo: %v", err)
		return nil
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	ret := make(map[string]string)
	for s.Scan() {
		tmp := strings.Split(s.Text(), " ")
		//          id: 1
		//      parent: 2
		//     maj:min: 3
		//        root: 4
		//      target: 5
		// vfs options: 6
		//    optional: 7, next "-"
		//     fs type: 8
		//      source: 9
		//  fs options: 10
		// ret[tmp[2]] = tmp[8]
		if tmp[6] == "-" {
			ret[tmp[8]] = tmp[2]
		} else {
			ret[tmp[9]] = tmp[2]
		}
	}
	return ret
}

func parseDev(str string) (int64, int64) {
	tmp := strings.SplitN(str, ":", 2)
	major, _ := strconv.ParseInt(tmp[0], 10, 64)
	minor, _ := strconv.ParseInt(tmp[1], 10, 64)
	return major, minor
}
