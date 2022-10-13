module github.com/jkstack/libagent

go 1.17

replace github.com/kardianos/service => github.com/lwch/service v1.2.1-1

require (
	github.com/btcsuite/winsvc v1.0.0
	github.com/containerd/cgroups v1.0.4
	github.com/dustin/go-humanize v1.0.0
	github.com/gorilla/websocket v1.5.0
	github.com/jkstack/anet v0.0.0-20221010100306-9a88844af68f
	github.com/jkstack/jkframe v1.0.9
	github.com/kardianos/service v1.2.1
	github.com/opencontainers/runtime-spec v1.0.2
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/shirou/gopsutil/v3 v3.22.9
)

require (
	github.com/cilium/ebpf v0.9.3 // indirect
	github.com/coreos/go-systemd/v22 v22.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.5.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.0.0-20221010170243-090e33056c14 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
