package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/johnnyfreeman/peek/internal/core"
	"github.com/johnnyfreeman/peek/internal/infra/file"
	"github.com/samber/lo"
)

func main() {
	code, out := Run(os.Args[1:], file.NewYAMLLoader(), core.NewDefaultRunner(http.DefaultClient), core.NewPrettyFormatter())
	log.Error(out)
	os.Exit(code)
}

func Run(args []string, loader core.Loader, runner core.Runner, formatter core.Formatter) (int, string) {
	ctx := context.Background()

	filename := args[0]

	requestGroup, err := loader.Load(ctx, filename)
	if err != nil {
		return 1, fmt.Sprintf("load error: %v", err)
	}

	// TODO: prompt user to pick request from group
	request := requestGroup.Requests[0]

	resolverCtx := &core.ResolverContext{
		Requests: lo.KeyBy(requestGroup.Requests, func(request core.Request) string {
			return request.Name
		}),
		Results: map[string]core.Result{},
		Prompt: func(name, prompt string) (string, error) {
			fmt.Printf("%s: ", prompt)
			var input string
			_, err := fmt.Scanln(&input)
			return input, err
		},
		Runner: runner,
	}

	if err := request.Resolve(ctx, resolverCtx); err != nil {
		panic(err)
	}
	log.Debug("request resolved", "url", request.URL)

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
