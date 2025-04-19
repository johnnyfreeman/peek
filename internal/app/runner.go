package app

import (
	"context"
	"io"
	"net/http"

	"github.com/johnnyfreeman/peek/internal/core/domain"
)

type Runner struct {
	Client *http.Client
}

func NewDefaultRunner(client *http.Client) Runner {
	return Runner{
		Client: client,
	}
}

func (r Runner) Run(ctx context.Context, request domain.Request) (domain.Result, error) {
	httpReq, err := request.ToHTTPRequest()
	if err != nil {
		return domain.Result{}, err
	}

	httpReq = httpReq.WithContext(ctx)

	resp, err := r.Client.Do(httpReq)
	if err != nil {
		return domain.Result{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Result{}, err
	}

	return domain.Result{
		RequestName: request.Name,
		StatusCode:  resp.StatusCode,
		Body:        body,
		Headers:     resp.Header,
	}, nil
}
