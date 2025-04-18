package app

import (
	"github.com/johnnyfreeman/peek/internal/core/domain"
)

type Formatter interface {
	Format(results []domain.Result) ([]byte, error)
}

type PrettyFormatter struct{}

func NewPrettyFormatter() Formatter {
	return PrettyFormatter{}
}

func (f PrettyFormatter) Format(results []domain.Result) ([]byte, error) {
	return []byte("formatted results go here"), nil
}

