//go:build windows || aix || darwin
// +build windows aix darwin

package limit

import (
	"math"
	"runtime"
	"runtime/debug"
)

// Do set cgroups limit
func (cfg *Configure) Do(_ string) {
	cpu := math.Ceil(float64(cfg.CPUQuota / 100))
	runtime.GOMAXPROCS(int(cpu))
	debug.SetMemoryLimit(int64(cfg.Memory))
}
