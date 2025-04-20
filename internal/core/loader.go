package core

import (
	"context"
)

type Loader interface {
	Load(ctx context.Context, filename string) (RequestGroup, error)
}
