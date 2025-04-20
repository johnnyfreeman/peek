package core

type RedactionRule struct {
	JSONPath string
	Regex    string
	Replace  string
}
