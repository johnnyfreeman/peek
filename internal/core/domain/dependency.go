package domain

import (
	"context"
)

// Core interface
type Dependency interface {
	Placeholder() string
	Resolve(ctx context.Context, rctx *ResolverContext) (string, error)
}

// Shared context for all resolvers (e.g. previous results, prompt function)
type ResolverContext struct {
	Results map[string]Result
	Prompt  func(name, prompt string) (string, error)
}
