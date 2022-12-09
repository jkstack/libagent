package agent

import (
	"context"
	"fmt"

	"github.com/jkstack/anet"
	"github.com/jkstack/libagent/conf"
	"github.com/jkstack/libagent/internal/utils"
	"github.com/kardianos/service"
)

// App app 接口，每一个agent必须实现以下接口
type App interface {
	// 获取当前agent名称
	AgentName() string
	// 获取当前agent版本号
	Version() string
	// 获取配置文件路径
	ConfDir() string
	// 获取libagent所需配置
	// 该对象必须是一个相对全局作用域的变量，在后续执行过程中该变量将会被更新
	Configure() *conf.Configure

	// 重置配置文件时的回调函数，在以下情况下会回调
	//   - 连接成功后服务端分配了新的agent id
	OnRewriteConfigure() error
	// 连接成功后的回调函数
	OnConnect()
	// 断开连接时的回调函数
	OnDisconnect()
	// 触发上报监控信息时的回调函数，该回调一般被用来上报一些自定义监控数据
	OnReportMonitor()
	// 收到数据包时的回调函数
	OnMessage(*anet.Msg) error

	// 返回数据包时的回调函数，该函数必须是一个循环，
	// 且在有数据需要返回时将其放入第二个参数中的队列内
	LoopWrite(context.Context, chan *anet.Msg) error
}

func deferCallback(name string, fn func()) {
	defer utils.Recover(name)
	fn()
}

// RegisterService 注册系统服务
func RegisterService(app App) error {
	svc, err := newService(app)
	if err != nil {
		return err
	}
	fmt.Printf("service name: %s\n", app.AgentName())
	fmt.Printf("platform: %s\n", svc.Platform())
	return svc.Install()
}

// UnregisterService 卸载系统服务
func UnregisterService(app App) error {
	svc, err := newService(app)
	if err != nil {
		return err
	}
	svc.Stop()
	return svc.Uninstall()
}

// Run 运行agent
func Run(app App) error {
	svc, err := newService(app)
	if err != nil {
		return err
	}
	return svc.Run()
}

// Start 启动agent
func Start(app App) error {
	svc, err := newService(app)
	if err != nil {
		return err
	}
	return svc.Start()
}

// Stop 停止agent
func Stop(app App) error {
	svc, err := newService(app)
	if err != nil {
		return err
	}
	return svc.Stop()
}

// Status 获取服务状态
func Status(app App) (service.Status, error) {
	svc, err := newService(app)
	if err != nil {
		return service.StatusUnknown, err
	}
	return svc.Status()
}

// Restart 重启agent
func Restart(app App) error {
	svc, err := newService(app)
	if err != nil {
		return err
	}
	status, err := svc.Status()
	if err != nil {
		return err
	}
	if status == service.StatusStopped {
		return svc.Start()
	}
	return svc.Restart()
}
