package domain

type RedactionRule struct {
	JSONPath string
	Regex    string
	Replace  string
}
