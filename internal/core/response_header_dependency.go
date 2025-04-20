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

func NewResponseHeaderDependency(placeholder, request, key string) ResponseHeaderDependency {
	return ResponseHeaderDependency{
		placeholder: placeholder,
		request:     request,
		key:         key,
	}
}

func (d ResponseHeaderDependency) Placeholder() string {
	return d.placeholder
}

func (d ResponseHeaderDependency) Resolve(ctx context.Context, rctx *ResolverContext) (string, error) {
	result, err := rctx.GetResult(ctx, d.request)
	if err != nil {
		return "", err
	}

	values := result.Headers[d.key]
	if len(values) == 0 {
		return "", fmt.Errorf("header %q not found in result for %q", d.key, d.request)
	}

	return values[0], nil
}
