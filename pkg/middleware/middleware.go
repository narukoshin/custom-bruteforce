package middleware

import (
	"net/http"
)

type Middleware struct {
	Client *http.Client
	Request *http.Request
}

func (m *Middleware) Do() error {
	// your custom code here...
	return nil
}
