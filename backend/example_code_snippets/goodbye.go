package example_code_snippets

import (
	"log"
	"net/http"
)

// Goodbye is a simple handler
type Goodbye struct {
	l *log.Logger
}

// NewGoodbye creates a new Goodbye handler with the given logger
func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (g *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// write the client
	_, err := w.Write([]byte("Goodbye"))
	if err != nil {
		return
	}
}
