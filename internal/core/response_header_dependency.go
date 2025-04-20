package core

import (
	"context"
	"fmt"
)

type ResponseHeaderDependency struct {
	placeholder string
	request     string
	key         string
}

func (d ResponseHeaderDependency) Placeholder() string {
	return d.placeholder
}

func (d ResponseHeaderDependency) Resolve(ctx context.Context, rctx *ResolverContext) (string, error) {
	result, ok := rctx.Results[d.request]
	if !ok {
		return "", fmt.Errorf("result for request %q not found", d.request)
	}
	values := result.Headers[d.key]
	if len(values) == 0 {
		return "", fmt.Errorf("header %q not found in result for %q", d.key, d.request)
	}
	return values[0], nil
}
