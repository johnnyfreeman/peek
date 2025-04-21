package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/johnnyfreeman/peek/internal/core"
)

type YamlLoader struct{}

func NewYAMLLoader() YamlLoader {
	return YamlLoader{}
}

func (l YamlLoader) Load(ctx context.Context, filename string) (core.RequestGroup, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return core.RequestGroup{}, err
	}

	path := filepath.Join(home, ".config", "peek", "requests", filename+".yaml")
	content, err := os.ReadFile(path)
	if err != nil {
		return core.RequestGroup{}, fmt.Errorf("failed to read file: %w", err)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return core.RequestGroup{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	group := core.RequestGroup{
		Name: raw["name"].(string),
	}

	rawRequests, ok := raw["requests"].([]any)
	if !ok {
		return core.RequestGroup{}, fmt.Errorf("missing or invalid 'requests'")
	}

	for _, r := range rawRequests {
		requestMap := r.(map[string]any)

		req := core.Request{
			Name:   requestMap["name"].(string),
			Method: requestMap["method"].(string),
			URL:    requestMap["url"].(string),
		}

		// headers
		if headersRaw, ok := requestMap["headers"].(map[string]any); ok {
			headers := make(map[string]string)
			for k, v := range headersRaw {
				headers[k] = fmt.Sprint(v)
			}
			req.Headers = headers
		}

		// dependencies
		if depsRaw, ok := requestMap["dependencies"].([]any); ok {
			for _, d := range depsRaw {
				depMap := d.(map[string]any)
				typ := depMap["type"].(string)
				placeholder := depMap["placeholder"].(string)

				switch typ {
				case "envvar":
					req.Dependencies = append(req.Dependencies, core.NewEnvVarDependency(
						placeholder,
						depMap["name"].(string),
						depMap["prompt"].(string),
					))
				case "response_body":
					req.Dependencies = append(req.Dependencies, core.NewResponseBodyDependency(
						placeholder,
						depMap["request"].(string),
						depMap["path"].(string),
					))
				default:
					return core.RequestGroup{}, fmt.Errorf("unknown dependency type: %s", typ)
				}
			}
		}

		group.Requests = append(group.Requests, req)
	}

	return group, nil
}
