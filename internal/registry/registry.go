package registry

import (
	"context"
	"time"
)

type Metadata struct {
	Exists    bool
	CreatedAt time.Time
}

type Registry interface {
	GetMetadata(ctx context.Context, name string) (*Metadata, error)
}
