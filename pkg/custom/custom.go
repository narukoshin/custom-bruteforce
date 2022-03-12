package custom

import (
	"net/http"
)

type Middleware struct {
	Client *http.Client
	Request *http.Request
}

func (m *Middleware) Do() error {
	return nil
}