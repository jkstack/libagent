package limit

import "github.com/jkstack/agent/utils"

type DiskLimit struct {
	// mountpoint eg: /
	MountPoint string `json:"mount_point" yaml:"mount_point"`
	// read bytes
	ReadBytes uint64 `json:"read_bytes" yaml:"read_bytes"`
	// write bytes
	WriteBytes uint64 `json:"write_bytes" yaml:"write_bytes"`
	// read iops
	ReadIOPS uint64 `json:"read_iops" yaml:"read_iops"`
	// write iops
	WriteIOPS uint64 `json:"write_iops" yaml:"write_iops"`
}

// Config limit configure
type Config struct {
	// cpu usage 100 means 1 core
	CpuQuota int64 `json:"cpu_quota" yaml:"cpu_quota"`
	// memory size limit in bytes
	Memory utils.Bytes `json:"memory_limit" yaml:"memory_limit"`
	// limit of disk
	Disks []DiskLimit `json:"disk_limit" yaml:"disk_limit"`
}
