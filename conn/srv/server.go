package srv

type Server interface {
	Start() error
	Stop() error
	Shutdown() error
}
