package waiter

type Waiter struct {
	Host   string
	Port   int
	Path   string
	Status int
}

func NewWaiter(host, path string, port int, status int) *Waiter {
	if (path == "") || (path[0] != '/') {
		path = "/" + path
	}
	return &Waiter{
		Host:   host,
		Port:   port,
		Path:   path,
		Status: status,
	}
}

type WaiterExecutor interface {
	ShouldExecute() bool
	IsReady() bool
}
