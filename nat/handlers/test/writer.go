package test

import (
	"testing"

	"github.com/gopherx/net/nat"
)

// Writer implements the ResponseWriter interface for testing purposes.
type Writer struct {
	T    *testing.T
	Msgs []nat.Message
}

func (w *Writer) Write(m nat.Message, opts *nat.PrintOptions) error {
	w.Msgs = append(w.Msgs, m)
	return nil
}

// Pop returns the first message from the list of written messages.
func (w *Writer) Pop() nat.Message {
	if len(w.Msgs) == 0 {
		w.T.Fatal("no message available to pop")
	}

	msg := w.Msgs[0]
	w.Msgs = w.Msgs[1:]
	return msg
}
