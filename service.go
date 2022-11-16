package agent

import "github.com/kardianos/service"

type builtinService interface {
	Install() error
	Uninstall() error
	Run() error
	Start() error
	Stop() error
	Restart() error
	Status() (service.Status, error)
	Platform() string
}
