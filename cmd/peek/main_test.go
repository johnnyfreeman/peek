package main

import (
	"context"
	"strings"
	"testing"

	"github.com/johnnyfreeman/peek/internal/core"
)

type fakeLoader struct {
	Group core.RequestGroup
	Err   error
}

func (f fakeLoader) Load(ctx context.Context, filename string) (core.RequestGroup, error) {
	return f.Group, f.Err
}

type fakeRunner struct {
	Result core.Result
	Err    error
}

func (f fakeRunner) Run(ctx context.Context, request core.Request) (core.Result, error) {
	return f.Result, f.Err
}

type fakeFormatter struct {
	Out string
	Err error
}

func (f fakeFormatter) Format(result core.Result) ([]byte, error) {
	if f.Err != nil {
		return nil, f.Err
	}
	return []byte(f.Out), nil
}

func TestRun_Success(t *testing.T) {
	args := []string{"run", "fake.yml"}

	loader := fakeLoader{Group: core.RequestGroup{Name: "test"}}
	runner := fakeRunner{Result: core.Result{RequestName: "one", StatusCode: 200}}
	formatter := fakeFormatter{Out: "Request one: 200 OK"}

	code, out := Run(args, loader, runner, formatter)

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}

	if !strings.Contains(out, "Request one") {
		t.Errorf("unexpected output: %s", out)
	}
}
