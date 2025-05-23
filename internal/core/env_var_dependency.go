package core

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
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
	log.Debug("looking up env variable",
		"name", d.name,
	)

	if val, ok := os.LookupEnv(d.name); ok {
		log.Debug("found env variable",
			"name", d.name,
			"value", val,
		)
		return val, nil
	}

	if rctx.Prompt != nil {
		return rctx.Prompt(d.name, d.prompt)
	}

	return "", fmt.Errorf("env var %q not set and no prompt available", d.name)
}
