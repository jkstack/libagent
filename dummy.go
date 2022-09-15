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

// DummyApp 用于系统服务注册的空App
type DummyApp struct {
	name    string
	confDir string
}

// NewDummyApp 创建一个dummy的APP用于系统服务注册等
func NewDummyApp(name, confDir string) *DummyApp {
	dir, err := filepath.Abs(confDir)
	utils.Assert(err)
	return &DummyApp{
		name:    name,
		confDir: dir,
	}
}

// AgentName agent名称
func (dm *DummyApp) AgentName() string {
	return dm.name
}

// Version 版本号（unset）
func (dm *DummyApp) Version() string {
	return "unset"
}

// ConfDir 配置文件路径
func (dm *DummyApp) ConfDir() string {
	return dm.confDir
}

// Configure 获取Configure对象
func (dm *DummyApp) Configure() *conf.Configure {
	return nil
}

// OnRewriteConfigure 重写配置文件回调接口
func (dm *DummyApp) OnRewriteConfigure() error {
	return errUnsupported
}

// OnConnect 连接成功时的回调接口
func (dm *DummyApp) OnConnect() {
}

// OnDisconnect 连接断开时的回调接口
func (dm *DummyApp) OnDisconnect() {
}

// OnReportMonitor 上报监控数据时的回调接口
func (dm *DummyApp) OnReportMonitor() {
}

// OnMessage 收到数据时的回调接口
func (dm *DummyApp) OnMessage(*anet.Msg) error {
	return errUnsupported
}

// LoopWrite 发送数据的回调接口
func (dm *DummyApp) LoopWrite(context.Context, chan *anet.Msg) error {
	return errUnsupported
}
