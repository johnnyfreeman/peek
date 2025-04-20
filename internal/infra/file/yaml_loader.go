// internal/infra/file/yaml_loader.go
package file

import (
	"context"

	"github.com/johnnyfreeman/peek/internal/core"
)

type Loader struct{}

func NewYAMLLoader() Loader {
	return Loader{}
}

func (l Loader) Load(ctx context.Context, filename string) (core.RequestGroup, error) {
	return core.RequestGroup{Name: "stub"}, nil
}
