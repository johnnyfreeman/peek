package domain

type RequestGroup struct {
    Name      string
    Env       Environment
    Requests  []Request
}

