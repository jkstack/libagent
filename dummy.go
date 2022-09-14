package agent

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/utils"
	"github.com/jkstack/libagent/conf"
)

var errUnsupported = errors.New("unsupported")

type dummy struct {
	name    string
	confDir string
}

// NewDummyApp 创建一个dummy的APP用于系统服务注册等
func NewDummyApp(name, confDir string) *dummy {
	dir, err := filepath.Abs(confDir)
	utils.Assert(err)
	return &dummy{
		name:    name,
		confDir: dir,
	}
}

func (dm *dummy) AgentName() string {
	return dm.name
}

func (dm *dummy) Version() string {
	return "unset"
}

func (dm *dummy) ConfDir() string {
	return dm.confDir
}

func (dm *dummy) Configure() *conf.Configure {
	return nil
}

func (dm *dummy) OnRewriteConfigure() error {
	return errUnsupported
}

func (dm *dummy) OnConnect() {
}

func (dm *dummy) OnDisconnect() {
}

func (dm *dummy) OnReportMonitor() {
}

func (dm *dummy) OnMessage(*anet.Msg) error {
	return errUnsupported
}

func (dm *dummy) LoopWrite(context.Context, chan *anet.Msg) error {
	return errUnsupported
}
