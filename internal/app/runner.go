// internal/app/runner.go
package app

import (
	"context"

	"github.com/johnnyfreeman/peek/internal/core/domain"
)

type Runner struct{}

func NewDefaultRunner() Runner {
	return Runner{}
}

func (r Runner) Run(ctx context.Context, group domain.RequestGroup) ([]domain.Result, error) {
	return []domain.Result{
		{RequestName: "stub", StatusCode: 200},
	}, nil
}
