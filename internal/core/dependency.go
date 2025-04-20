package core

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
)

// Core interface
type Dependency interface {
	Placeholder() string
	Resolve(ctx context.Context, rctx *ResolverContext) (string, error)
}

// Shared context for all resolvers (e.g. previous results, prompt function)
type ResolverContext struct {
	Requests map[string]Request
	Results  map[string]Result
	Prompt   func(name, prompt string) (string, error)
	Runner   Runner
}

func (rctx *ResolverContext) GetResult(ctx context.Context, request string) (Result, error) {
	if result, ok := rctx.Results[request]; ok {
		return result, nil
	}

	if request, ok := rctx.Requests[request]; ok {
		if err := request.Resolve(ctx, rctx); err != nil {
			panic(err)
		}
		log.Debug("request resolved", "url", request.URL)

		result, err := rctx.Runner.Run(ctx, request)
		if err != nil {
			return Result{}, err
		}

		rctx.Results[request.Name] = result

		return result, nil
	}

	return Result{}, fmt.Errorf("result for request %q not found", request)
}
