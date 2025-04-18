package main

import (
	"context"
	"fmt"
	"os"

	"github.com/johnnyfreeman/peek/internal/app"
	"github.com/johnnyfreeman/peek/internal/core/domain"
	"github.com/johnnyfreeman/peek/internal/infra/file"
)

type Loader interface {
	Load(ctx context.Context, filename string) (domain.RequestGroup, error)
}

type Runner interface {
	Run(ctx context.Context, group domain.RequestGroup) ([]domain.Result, error)
}

type Formatter interface {
	Format(results []domain.Result) ([]byte, error)
}

func main() {
	code, out := Run(os.Args[1:], file.NewYAMLLoader(), app.NewDefaultRunner(), app.NewPrettyFormatter())
	fmt.Println(out)
	os.Exit(code)
}

func Run(args []string, loader Loader, runner Runner, formatter Formatter) (int, string) {
	ctx := context.Background()

	if len(args) < 2 || args[0] != "run" {
		return 1, "Usage: peek run <file>"
	}

	filename := args[1]

	group, err := loader.Load(ctx, filename)
	if err != nil {
		return 1, fmt.Sprintf("load error: %v", err)
	}

	results, err := runner.Run(ctx, group)
	if err != nil {
		return 1, fmt.Sprintf("execution error: %v", err)
	}

	out, err := formatter.Format(results)
	if err != nil {
		return 1, fmt.Sprintf("format error: %v", err)
	}

	return 0, string(out)
}
