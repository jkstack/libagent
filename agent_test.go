package agent

import (
	"context"
	"testing"

	"github.com/jkstack/anet"
	"github.com/jkstack/libagent/conf"
)

type agent struct{}

func (a *agent) AgentName() string {
	return "example-agent"
}

func (a *agent) Version() string {
	return "0.0.0"
}

func (a *agent) ConfDir() string {
	return "./conf"
}

func (a *agent) Configure() *conf.Configure {
	return &conf.Configure{}
}

func (a *agent) OnRewriteConfigure() error {
	return nil
}

func (a *agent) OnConnect() {
}

func (a *agent) OnDisconnect() {
}

func (a *agent) OnReportMonitor() {
}

func (a *agent) OnMessage(*anet.Msg) error {
	return nil
}

func (a *agent) LoopWrite(context.Context, chan *anet.Msg) error {
	select {}
}

func TestRestart(t *testing.T) {
	Restart(&agent{})
}
