//go:build (windows || aix || darwin) && !go1.17 && !go1.18
// +build windows aix darwin
// +build !go1.17
// +build !go1.18

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
