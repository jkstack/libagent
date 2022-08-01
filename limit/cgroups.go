//go:build !windows && !aix
// +build !windows,!aix

package limit

import (
	"github.com/containerd/cgroups"
	cgroupsv2 "github.com/containerd/cgroups/v2"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// Do set cgroups limit
func Do(agentName string, cfg *Config) {
	if cgroups.Mode() == cgroups.Unified {
		doV2(agentName, cfg)
		return
	}
	doV1(agentName, cfg)
}

func doV2(agentName string, cfg *Config) {
	// TODO
	cgroupsv2.NewSystemd("/", "my-group.slice", -1, &cgroupsv2.Resources{})
}

func doV1(agentName string, cfg *Config) {
	// TODO
	cgroups.New(cgroups.V1, cgroups.StaticPath("/"), &specs.LinuxResources{})
}
