package app

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/johnnyfreeman/peek/internal/core/domain"
	"github.com/tidwall/pretty"
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

	// Styles
	nameStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("35")).Foreground(lipgloss.Color("#fff")).Padding(0, 1)
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("239"))
	value := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// Request
	buf.WriteString(nameStyle.Render(result.RequestName) + " " + value.Render(fmt.Sprintf("%d", result.StatusCode)) + "\n")

	// Headers
	if len(result.Headers) > 0 {
		buf.WriteString("\n")
	}

	for key, vals := range result.Headers {
		for _, val := range vals {
			buf.WriteString(header.Render(key+": ") + value.Render(val) + "\n")
		}
	}

	// Body
	buf.WriteString("\n")
	json := pretty.Pretty(result.Body)
	buf.Write(pretty.Color(json, nil))
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
