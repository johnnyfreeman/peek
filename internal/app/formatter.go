package app

import (
	"bytes"
	"fmt"

	"github.com/johnnyfreeman/peek/internal/core/domain"
)

type Formatter interface {
	Format(result domain.Result) ([]byte, error)
}

type PrettyFormatter struct{}

func NewPrettyFormatter() Formatter {
	return PrettyFormatter{}
}

func (f PrettyFormatter) Format(result domain.Result) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("→ Request: %s\n", result.RequestName))
	buf.WriteString(fmt.Sprintf("→ Status: %d\n", result.StatusCode))

	if len(result.Headers) > 0 {
		buf.WriteString("→ Headers:\n")
		for key, values := range result.Headers {
			for _, val := range values {
				buf.WriteString(fmt.Sprintf("   %s: %s\n", key, val))
			}
		}
	}

	buf.WriteString("→ Body:\n")
	buf.Write(result.Body)

	return buf.Bytes(), nil
}
