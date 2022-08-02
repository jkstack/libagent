//go:build !(windows && 386)
// +build !windows !386

package agent

import (
	rt "runtime"

	"github.com/jkstack/jkframe/utils"
	"github.com/kardianos/service"
)

type svr struct {
	app *app
}

func newService(app App) builtinService {
	var user string
	var depends []string
	if rt.GOOS != "windows" {
		user = "root"
		depends = append(depends, "After=network.target")
	}
	appCfg := &service.Config{
		Name:         app.AgentName(),
		DisplayName:  app.AgentName(),
		Description:  app.AgentName(),
		UserName:     user,
		Arguments:    []string{"-conf", app.ConfDir()},
		Dependencies: depends,
	}
	svr, err := service.New(&svr{app: newApp(app)}, appCfg)
	utils.Assert(err)
	return svr
}

func (svr *svr) Start(s service.Service) error {
	go svr.app.start()
	return nil
}

func (svr *svr) Stop(s service.Service) error {
	svr.app.stop()
	return nil
}
