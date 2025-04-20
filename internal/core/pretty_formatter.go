package core

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/tidwall/pretty"
)

type PrettyFormatter struct{}

func NewPrettyFormatter() Formatter {
	return PrettyFormatter{}
}

func (f PrettyFormatter) Format(result Result) ([]byte, error) {
	var buf bytes.Buffer

	// Styles
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("239"))
	value := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	statusStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fff")).Padding(0, 1)

	switch {
	case result.StatusCode >= 200 && result.StatusCode <= 299:
		statusStyle = statusStyle.Background(lipgloss.Color("35")) // green
	case result.StatusCode >= 300 && result.StatusCode <= 399:
		statusStyle = statusStyle.Background(lipgloss.Color("27")) // blue
	case result.StatusCode >= 400 && result.StatusCode <= 499:
		statusStyle = statusStyle.Background(lipgloss.Color("208")) // orange
	case result.StatusCode >= 500 && result.StatusCode <= 599:
		statusStyle = statusStyle.Background(lipgloss.Color("160")) // red
	default:
		statusStyle = statusStyle.Background(lipgloss.Color("238")) // gray fallback
	}

	// Request
	buf.WriteString(statusStyle.Render(fmt.Sprintf("%d", result.StatusCode)) + " " + value.Render(result.RequestName) + "\n")

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
