package agent

import (
	"context"

	"github.com/jkstack/agent/conf"
	"github.com/jkstack/anet"
)

// App app interface
type App interface {
	// get agent name
	AgentName() string
	// get agent version
	Version() string
	// get configure file path
	ConfDir() string
	// get agent configure
	Configure() *conf.Configure
	// rewrite configure file
	RewriteConfigure() error

	// connect callback
	OnConnect()
	// disconnect callback
	OnDisconnect()
	// report callback
	OnReportMonitor()
	// message received
	OnMessage(*anet.Msg) error

	// loop write
	LoopWrite(context.Context, chan *anet.Msg) error
}

// RegisterService register system service
func RegisterService(app App) error {
	svc := newService(app)
	return svc.Install()
}

// UnregisterService unregister system service
func UnregisterService(app App) error {
	svc := newService(app)
	svc.Stop()
	return svc.Uninstall()
}

// Run run agent
func Run(app App) error {
	svc := newService(app)
	return svc.Run()
}
