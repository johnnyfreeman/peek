package core

type RequestGroup struct {
	Name     string
	Env      Environment
	Requests []Request
}
