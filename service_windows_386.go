package agent

import (
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/winsvc/eventlog"
	"github.com/btcsuite/winsvc/mgr"
	"github.com/btcsuite/winsvc/svc"
	"github.com/jkstack/jkframe/logging"
	"github.com/kardianos/service"
)

type svr struct {
	agentName string
	app       *app
	exepath   string
}

func newService(app App) (builtinService, error) {
	exepath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	return &svr{
		agentName: app.AgentName(),
		app:       newApp(app),
		exepath:   exepath,
	}, nil
}

func (svr *svr) Install() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svr.agentName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", svr.agentName)
	}
	s, err = m.CreateService(svr.agentName, svr.exepath, mgr.Config{
		DisplayName: svr.agentName,
	})
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(svr.agentName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}

func (svr *svr) Uninstall() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svr.agentName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", svr.agentName)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(svr.agentName)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}
	return nil
}

func (svr *svr) Run() error {
	return svc.Run(svr.agentName, svr)
}

func (svr *svr) Start() error {
	m, err := mgr.Connect()
	if err != nil {
		logging.Error("mgr connect: %v", err)
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svr.agentName)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start([]string{"p1", "p2", "p3"})
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func (svr *svr) Stop() error {
	svr.app.stop()
	m, err := mgr.Connect()
	if err != nil {
		logging.Error("mgr connect: %v", err)
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svr.agentName)
	if err != nil {
		logging.Error("open service: %v", err)
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(svc.Stop)
	if err != nil {
		logging.Error("control stop: %v", err)
		return fmt.Errorf("could not send control=%d: %v", svc.Stop, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			logging.Error("timeout")
			return fmt.Errorf("timeout waiting for service to go to state=%d", svc.Stopped)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			logging.Error("query status: %v", err)
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func (svr *svr) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	go svr.app.start()
	for c := range r {
		logging.Info("request: %d => %d", c.CurrentStatus.State, c.Cmd)
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			break
		default:
			logging.Error("unexpected control request #%d", c)
		}
	}
	svr.app.stop()
	changes <- svc.Status{State: svc.StopPending}
	return false, 0
}

func (svr *svr) Platform() string {
	return "windows-service"
}

func (svr *svr) Status() (service.Status, error) {
	m, err := mgr.Connect()
	if err != nil {
		logging.Error("mgr connect: %v", err)
		return service.StatusUnknown, err
	}
	defer m.Disconnect()
	s, err := m.OpenService(svr.agentName)
	if err != nil {
		return service.StatusUnknown, fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Query()
	if err != nil {
		logging.Error("query status: %v", err)
		return service.StatusUnknown, fmt.Errorf("could not retrieve service status: %v", err)
	}
	switch status.State {
	case svc.Stopped:
		return service.StatusStopped, nil
	case svc.Running:
		return service.StatusRunning, nil
	default:
		return service.StatusUnknown, nil
	}
}

func (svr *svr) Restart() error {
	err := svr.Stop()
	if err != nil {
		return err
	}
	return svr.Start()
}
