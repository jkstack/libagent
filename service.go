package agent

type builtinService interface {
	Install() error
	Uninstall() error
	Run() error
	Start() error
	Stop() error
	Platform() string
}
