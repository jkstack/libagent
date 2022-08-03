package conf

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jkstack/agent/utils"
	"github.com/jkstack/jkframe/logging"
	runtime "github.com/jkstack/jkframe/utils"
)

// RewriteServer rewrite server configure
//   - use $... to set value by envionment variables
func (cfg *Configure) RewriteServer() {
	cfg.Server = strings.TrimSpace(cfg.Server)
	if len(cfg.Server) == 0 {
		panic("no server address set")
	}
	if cfg.Server[0] == '$' {
		logging.Info("rewrite server address by env %s", cfg.Server)
		str, ok := os.LookupEnv(cfg.Server[1:])
		if !ok {
			panic(fmt.Sprintf("can not get server address by env %s", cfg.Server))
		}
		cfg.Server = str
		logging.Info("now server address is %s", cfg.Server)
	}
}

// RewriteID rewrite agent id configure
//   - use $IP to set value by ip address in interface to connect server
//   - use $HOSTNAME to set value by hostname
//   - use $... to set value by envionment variables
func (cfg *Configure) RewriteID() {
	cfg.RewriteServer()
	cfg.ID = strings.TrimSpace(cfg.ID)
	if len(cfg.ID) == 0 {
		return
	}
	if cfg.ID[0] != '$' {
		return
	}
	logging.Info("rewrite agent id by env %s", cfg.ID)
	switch cfg.ID {
	case "$IP":
		conn, err := net.DialTimeout("tcp", cfg.Server, 10*time.Second)
		runtime.Assert(err)
		defer conn.Close()
		cfg.ID = utils.GetIP(conn).String()
	case "$HOSTNAME":
		name, err := os.Hostname()
		runtime.Assert(err)
		cfg.ID = name
	default:
		str, ok := os.LookupEnv(cfg.ID[1:])
		if !ok {
			panic(fmt.Sprintf("can not get agent id by env %s", cfg.ID))
		}
		cfg.ID = str
	}
	logging.Info("now agent id is %s", cfg.ID)
}

// SetAgentID reset agent id by value
func (cfg *Configure) SetAgentID(id string) {
	cfg.ID = id
}
