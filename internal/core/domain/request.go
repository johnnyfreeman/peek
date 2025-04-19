package domain

import (
	"io"
	"net/http"
)

type Request struct {
	Name         string
	Method       string
	URL          string
	Headers      map[string]string
	Body         io.Reader
	Placeholders []Placeholder
}

func (r *Request) ToHTTPRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range r.Headers {
		req.Header.Set(key, val)
	}

	return req, nil
}
