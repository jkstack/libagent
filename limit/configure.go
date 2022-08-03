package limit

import (
	"encoding/json"

	"github.com/containerd/cgroups"
	"github.com/jkstack/agent/utils"
)

type DiskLimit struct {
	// device version of lsblk command eg: 8:0
	Dev string `json:"dev" yaml:"dev" kv:"dev"`
	// read bytes
	ReadBytes uint64 `json:"read_bytes" yaml:"read_bytes" kv:"read_bytes"`
	// write bytes
	WriteBytes uint64 `json:"write_bytes" yaml:"write_bytes" kv:"write_bytes"`
	// read iops
	ReadIOPS uint64 `json:"read_iops" yaml:"read_iops" kv:"read_iops"`
	// write iops
	WriteIOPS uint64 `json:"write_iops" yaml:"write_iops" kv:"write_iops"`
}

// Configure limit configure
type Configure struct {
	// cpu usage 100 means 1 core
	CpuQuota int64 `json:"cpu_quota" yaml:"cpu_quota" kv:"cpu_quota"`
	// memory size limit in bytes
	Memory utils.Bytes `json:"memory_limit" yaml:"memory_limit" kv:"memory_limit"`
	// limit of disk
	Disks diskLimits `json:"disk_limit" yaml:"disk_limit" kv:"disk_limit"`
}

type diskLimits []DiskLimit

func (limits diskLimits) MarshalKV() (string, error) {
	data, err := json.Marshal(limits)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (limits *diskLimits) UnmarshalKV(data string) error {
	return json.Unmarshal([]byte(data), limits)
}

// Do set cgroups limit
func (cfg *Configure) Do(agentName string) {
	if cgroups.Mode() == cgroups.Unified {
		cfg.doV2(agentName)
		return
	}
	cfg.doV1(agentName)
}
