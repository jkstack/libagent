//go:build !windows && !aix && !darwin
// +build !windows,!aix,!darwin

package limit

import "github.com/containerd/cgroups/v3"

// Do set cgroups limit
func (cfg *Configure) Do(agentName string) {
	if cgroups.Mode() == cgroups.Unified {
		cfg.doV2(agentName)
		return
	}
	cfg.doV1(agentName)
}
