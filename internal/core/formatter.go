package core

type Formatter interface {
	Format(result Result) ([]byte, error)
}
