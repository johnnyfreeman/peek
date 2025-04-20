package core

import (
	"context"
)

type Runner interface {
	Run(context.Context, Request) (Result, error)
}
