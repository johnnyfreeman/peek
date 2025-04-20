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
	fmt.Println(out)
	os.Exit(code)
}

func Run(args []string, loader core.Loader, runner core.Runner, formatter core.Formatter) (int, string) {
	ctx := context.Background()

	// if len(args) < 2 || args[0] != "run" {
	// 	return 1, "Usage: peek run <file>"
	// }

	// filename := args[1]

	// group, err := loader.Load(ctx, filename)
	// if err != nil {
	// 	return 1, fmt.Sprintf("load error: %v", err)
	// }

	requestGroup := core.RequestGroup{
		Name: "Openweather API",
		Env:  core.Environment{},
		Requests: []core.Request{
			{
				Name:   "Get Current Weather",
				Method: http.MethodGet,
				URL:    "https://api.openweathermap.org/data/2.5/weather?q={city}&appid={api_key}",
				Headers: map[string]string{
					"Accept": "application/json",
				},
				Dependencies: []core.Dependency{
					core.NewEnvVarDependency("city", "CITY", "Enter the city name"),
					core.NewEnvVarDependency("api_key", "OPENWEATHER_API_KEY", "Enter your OpenWeather API key"),
				},
			},
			{
				Name:   "Get 5-Day Forecast",
				Method: http.MethodGet,
				URL:    "https://api.openweathermap.org/data/2.5/forecast?lat={lat}&lon={lon}&appid={api_key}",
				Headers: map[string]string{
					"Accept": "application/json",
				},
				Dependencies: []core.Dependency{
					core.NewResponseBodyDependency("lat", "Get Current Weather", "$.coord.lat"),
					core.NewResponseBodyDependency("lon", "Get Current Weather", "$.coord.lon"),
					core.NewEnvVarDependency("api_key", "OPENWEATHER_API_KEY", "Enter your OpenWeather API key"),
					// core.NewOnePasswordDependency("Private", "OpenWeather", "api-key"),
				},
			},
		},
	}

	// TODO: prompt user to pick request from group
	request := requestGroup.Requests[1]

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
