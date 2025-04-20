package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Name         string
	Method       string
	URL          string
	Headers      map[string]string
	Body         io.Reader
	Dependencies []Dependency
}

func (r *Request) ToHTTPRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return nil, err
	}

	for key, val := range r.Headers {
		req.Header.Set(key, val)
	}

	return req, nil
}

func (r *Request) Resolve(ctx context.Context, resolverCtx *ResolverContext) error {
	resolved := make(map[string]string)

	for _, dep := range r.Dependencies {
		val, err := dep.Resolve(ctx, resolverCtx)
		if err != nil {
			return fmt.Errorf("resolving %q: %w", dep.Placeholder(), err)
		}
		resolved[dep.Placeholder()] = val
	}

	// Replace placeholders in URL
	r.URL = replacePlaceholders(r.URL, resolved)

	// Replace placeholders in headers
	for k, v := range r.Headers {
		r.Headers[k] = replacePlaceholders(v, resolved)
	}

	// Body templating is not supported (yet)
	// Could be added here if needed.

	return nil
}

func replacePlaceholders(s string, values map[string]string) string {
	for key, val := range values {
		s = strings.ReplaceAll(s, "{"+key+"}", val)
	}
	return s
}
