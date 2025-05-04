package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	home, err := os.UserHomeDir()
	if err != nil {
		return 1, fmt.Sprintf("user home directory error: %v", err)
	}

	path := filepath.Join(home, ".config", "peek", "requests")
	files, err := os.ReadDir(path)

	file, err := chooseRequestGroup(files)
	if err != nil {
		return 1, fmt.Sprintf("choose request error: %v", err)
	}

	requestGroup, err := loader.Load(ctx, file.Name())
	if err != nil {
		return 1, fmt.Sprintf("load error: %v", err)
	}

	request, err := chooseRequest(requestGroup)
	if err != nil {
		return 1, fmt.Sprintf("choose request error: %v", err)
	}

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

func chooseRequest(group core.RequestGroup) (core.Request, error) {
	fmt.Println("Choose a request:")

	for i, req := range group.Requests {
		fmt.Printf("[%d] %s\n", i+1, req.Name)
	}

	fmt.Print("Enter number: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return core.Request{}, err
	}

	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(group.Requests) {
		return core.Request{}, fmt.Errorf("invalid selection")
	}

	return group.Requests[index-1], nil
}

func chooseRequestGroup(groups []os.DirEntry) (os.DirEntry, error) {
	fmt.Println("Choose a request group:")

	for i, file := range groups {
		fmt.Printf("[%d] %s\n", i+1, file.Name())
	}

	fmt.Print("Enter number: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	input = strings.TrimSpace(input)
	index, err := strconv.Atoi(input)
	if err != nil || index < 1 || index > len(groups) {
		return nil, fmt.Errorf("invalid selection")
	}

	return groups[index-1], nil
}
