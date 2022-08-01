package agent

import (
	"github.com/jkstack/anet"
)

// App app interface
type App interface {
	AgentName() string
	OnConnect()
	Disconnect()
	OnMessage(*anet.Msg)
}

// Run run agent
func Run(app App, cfg Config) {
}
