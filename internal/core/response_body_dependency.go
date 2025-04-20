package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/oliveagle/jsonpath"
)

type ResponseBodyDependency struct {
	placeholder string
	request     string
	pointer     string
	runner      Runner
}

func NewResponseBodyDependency(placeholder, request, pointer string, runner Runner) ResponseBodyDependency {
	return ResponseBodyDependency{
		placeholder: placeholder,
		request:     request,
		pointer:     pointer,
		runner:      runner,
	}
}

func (d ResponseBodyDependency) Placeholder() string {
	return d.placeholder
}

func (d ResponseBodyDependency) Resolve(ctx context.Context, rctx *ResolverContext) (string, error) {

	if result, ok := rctx.Results[d.request]; ok {
		return extractJSONPointer(result.Body, d.pointer)
	}

	if request, ok := rctx.Requests[d.request]; ok {
		if err := request.Resolve(ctx, rctx); err != nil {
			panic(err)
		}
		log.Debug("request resolved", "url", request.URL)

		result, err := d.runner.Run(ctx, request)
		if err != nil {
			return "", err
		}

		rctx.Results[request.Name] = result

		return extractJSONPointer(result.Body, d.pointer)
	}

	return "", fmt.Errorf("result for request %q not found", d.request)
}

func extractJSONPointer(body []byte, pointer string) (string, error) {
	var doc any
	if err := json.Unmarshal(body, &doc); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	val, err := jsonpath.JsonPathLookup(doc, pointer)
	if err != nil {
		return "", fmt.Errorf("invalid JSON pointer %q: %w", pointer, err)
	}

	// Flatten result to string
	switch v := val.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%v", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case nil:
		return "", nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to re-marshal pointer value: %w", err)
		}
		return string(b), nil
	}
}
