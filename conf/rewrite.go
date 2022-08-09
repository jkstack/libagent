package conf

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jkstack/jkframe/logging"
	runtime "github.com/jkstack/jkframe/utils"
	"github.com/jkstack/libagent/internal/utils"
)

// RewriteServer 重写服务器地址配置，每次连接时自动调用
//   - 当给定值以$开头则表示使用系统环境变量
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

// RewriteID 重写agent id配置，每次连接时自动调用
//   - $HOSTNAME: 使用当前主机名作为agent id
//   - $IP: 使用连接到服务器端的网卡IP作为agent id
//   - ${env}: 使用环境变量作为agent id
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

// SetAgentID 重设当前agent id，当服务器端自动分配了新的agent id时被调用
func (cfg *Configure) SetAgentID(id string) {
	cfg.ID = id
}
