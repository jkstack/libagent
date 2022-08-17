package limit

import (
	"encoding/json"

	"github.com/jkstack/jkframe/utils"
)

// DiskLimit 磁盘限制配置
type DiskLimit struct {
	// 磁盘设备编号，可使用lsblk进行查询，如: 8:0
	Dev string `json:"dev" yaml:"dev" kv:"dev"`
	// 每秒读取字节数，为0表示不限制
	ReadBytes utils.Bytes `json:"read_bytes" yaml:"read_bytes" kv:"read_bytes"`
	// 每秒写入字节数，为0表示不限制
	WriteBytes utils.Bytes `json:"write_bytes" yaml:"write_bytes" kv:"write_bytes"`
	// 每秒并发读取次数，为0表示不限制
	ReadIOPS uint64 `json:"read_iops" yaml:"read_iops" kv:"read_iops"`
	// 每秒并发写入次数，为0表示不限制
	WriteIOPS uint64 `json:"write_iops" yaml:"write_iops" kv:"write_iops"`
}

// Configure 资源限制配置
type Configure struct {
	// CPU使用率限制，100表示1个核心
	CPUQuota int64 `json:"cpu_quota" yaml:"cpu_quota" kv:"cpu_quota"`
	// 内存使用限制
	Memory utils.Bytes `json:"memory_limit" yaml:"memory_limit" kv:"memory_limit"`
	// 磁盘限制
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
