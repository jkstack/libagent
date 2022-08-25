//go:build (windows || aix || darwin) && (go1.17 || go1.18)
// +build windows aix darwin
// +build go1.17 go1.18

package limit

// Do set cgroups limit
func (cfg *Configure) Do(_ string) {
}
