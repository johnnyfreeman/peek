package domain

type Result struct {
	RequestName string
	StatusCode  int
	Headers     map[string][]string
	Body        []byte
	Exports     map[string]string
}
