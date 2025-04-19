package app

import (
	"bytes"
	"fmt"
	"strings"

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
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	label := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("245"))
	value := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(strings.Repeat("─", 60))

	// Request
	buf.WriteString(header.Render("Request: ") + value.Render(result.RequestName) + "\n")
	buf.WriteString(header.Render("Status:  ") + value.Render(fmt.Sprintf("%d", result.StatusCode)) + "\n")
	buf.WriteString(divider + "\n")

	// Headers
	if len(result.Headers) > 0 {
		buf.WriteString(header.Render("Headers:") + "\n")
		for key, vals := range result.Headers {
			for _, val := range vals {
				buf.WriteString("  " + label.Render(key) + ": " + value.Render(val) + "\n")
			}
		}
		buf.WriteString(divider + "\n")
	}

	// Body
	buf.WriteString(header.Render("Body:") + "\n")
	json := pretty.Pretty(result.Body)
	buf.Write(pretty.Color(json, nil))
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
