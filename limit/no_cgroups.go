//go:build windows || aix || darwin
// +build windows aix darwin

package limit

// Do set cgroups limit
func (cfg *Configure) Do(agentName string) {
}
