package core

import (
	"context"
	"io"
	"net/http"
)

type HttpRunner struct {
	Client *http.Client
}

func NewDefaultRunner(client *http.Client) Runner {
	return HttpRunner{
		Client: client,
	}
}

func (r HttpRunner) Run(ctx context.Context, request Request) (Result, error) {
	httpReq, err := request.ToHTTPRequest()
	if err != nil {
		return Result{}, err
	}

	httpReq = httpReq.WithContext(ctx)

	resp, err := r.Client.Do(httpReq)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, err
	}

	return Result{
		RequestName: request.Name,
		StatusCode:  resp.StatusCode,
		Body:        body,
		Headers:     resp.Header,
	}, nil
}
