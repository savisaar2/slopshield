package registry

import "time"

type Metadata struct {
	Exists    bool
	CreatedAt time.Time
}

type Registry interface {
	GetMetadata(name string) (*Metadata, error)
}
