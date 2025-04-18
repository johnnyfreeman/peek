package domain

type Dependency struct {
	Kind    string // "envvar", "response", "vault", etc.
	Details map[string]string
}
