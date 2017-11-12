package nat

import (
	"testing"
	"time"
)

func TestNewNonceAttribute(t *testing.T) {
	a, err := NewNonceAttribute()
	if err != nil {
		t.Fatal(err)
	}

	if len(a.Nonce) != 127 {
		t.Fatalf("wrong nonce size; got:%d want:%d", len(a.Nonce), 127)
	}

	t.Log(a)
}

func BenchmarkNewNonceAttribute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewNonceAttribute()
	}
}

func TestNewNonceAttributeT0(t *testing.T) {
	t0 := time.Now()
	NewNonceAttribute()
	t.Log(time.Now().Sub(t0))
}
