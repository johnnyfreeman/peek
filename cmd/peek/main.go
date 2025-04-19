package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/johnnyfreeman/peek/internal/app"
	"github.com/johnnyfreeman/peek/internal/core/domain"
	"github.com/johnnyfreeman/peek/internal/infra/file"
)

type Loader interface {
	Load(ctx context.Context, filename string) (domain.RequestGroup, error)
}

type Runner interface {
	Run(context.Context, domain.Request) (domain.Result, error)
}

type Formatter interface {
	Format(results domain.Result) ([]byte, error)
}

func main() {
	code, out := Run(os.Args[1:], file.NewYAMLLoader(), app.NewDefaultRunner(http.DefaultClient), app.NewPrettyFormatter())
	fmt.Println(out)
	os.Exit(code)
}

func Run(args []string, loader Loader, runner Runner, formatter Formatter) (int, string) {
	ctx := context.Background()

	// if len(args) < 2 || args[0] != "run" {
	// 	return 1, "Usage: peek run <file>"
	// }

	// filename := args[1]

	// group, err := loader.Load(ctx, filename)
	// if err != nil {
	// 	return 1, fmt.Sprintf("load error: %v", err)
	// }

	// TODO: prompt user to pick request from group and pass to runner
	request := domain.Request{
		Name:   "Get Current Weather",
		Method: http.MethodGet,
		URL:    "https://api.openweathermap.org/data/2.5/weather?q={city}&appid={api_key}",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}

	result, err := runner.Run(ctx, request)
	if err != nil {
		return 1, fmt.Sprintf("execution error: %v", err)
	}

	out, err := formatter.Format(result)
	if err != nil {
		return 1, fmt.Sprintf("format error: %v", err)
	}

	return 0, string(out)
}
