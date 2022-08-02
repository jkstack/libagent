package limit

import (
	"testing"

	"github.com/jkstack/agent/utils"
)

func TestCGroups(t *testing.T) {
	var cfg Configure
	cfg.CpuQuota = 100                          // 1core
	cfg.Memory = utils.Bytes(100 * 1024 * 1024) // 100MB
	cfg.Disks = append(cfg.Disks, DiskLimit{
		MountPoint: "/dev/sdb",
		ReadBytes:  1024, // 1KB
	})
	// want root permission
	cfg.Do("testing")
}
