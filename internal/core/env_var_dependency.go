package core

import (
	"context"
	"fmt"
	"os"
)

type EnvVarDependency struct {
	placeholder string
	name        string
	prompt      string
}

func NewEnvVarDependency(placeholder, name, prompt string) EnvVarDependency {
	return EnvVarDependency{
		placeholder: placeholder,
		name:        name,
		prompt:      prompt,
	}
}

func (d EnvVarDependency) Placeholder() string {
	return d.placeholder
}

func (d EnvVarDependency) Resolve(ctx context.Context, rctx *ResolverContext) (string, error) {
	if val, ok := os.LookupEnv(d.name); ok {
		return val, nil
	}
	if rctx.Prompt != nil {
		return rctx.Prompt(d.name, d.prompt)
	}
	return "", fmt.Errorf("env var %q not set and no prompt available", d.name)
}
