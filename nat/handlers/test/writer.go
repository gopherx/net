package test

import (
	"github.com/gopherx/net/nat"
)

// Writer implements the ResponseWriter interface for testing purposes.
type Writer struct {
	Written []nat.Message
}

func (w *Writer) Write(m nat.Message) error {
	w.Written = append(w.Written, m)
	return nil
}
