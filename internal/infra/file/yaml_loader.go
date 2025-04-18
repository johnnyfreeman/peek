// internal/infra/file/yaml_loader.go
package file

import (
	"context"

	"github.com/johnnyfreeman/peek/internal/core/domain"
)

type Loader struct{}

func NewYAMLLoader() Loader {
	return Loader{}
}

func (l Loader) Load(ctx context.Context, filename string) (domain.RequestGroup, error) {
	return domain.RequestGroup{Name: "stub"}, nil
}
