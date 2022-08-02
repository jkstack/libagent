package agent

type builtinService interface {
	Install() error
	Uninstall() error
	Run() error
	Stop() error
}
