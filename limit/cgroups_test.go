package limit

import (
	"testing"

	"github.com/jkstack/jkframe/utils"
)

func TestCGroups(t *testing.T) {
	var cfg Configure
	cfg.CPUQuota = 100                          // 1core
	cfg.Memory = utils.Bytes(100 * 1024 * 1024) // 100MB
	cfg.Disks = append(cfg.Disks, DiskLimit{
		Dev:       "8:0",
		ReadBytes: 1024, // 1KB
	})
	// want root permission
	cfg.Do("testing")
}
