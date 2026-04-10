package registry

type Registry interface {
	Exists(name string) (bool, error)
}
