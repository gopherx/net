package nat

import (
	"reflect"
	"testing"

	"github.com/gopherx/base/errors"
	"github.com/gopherx/base/errors/codes"
)

func pad2header(b []byte) []byte {
	tmp := make([]byte, HeaderSize)
	copy(tmp, b)
	return tmp
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		desc string
		b    []byte
		code codes.Code
		m    Message
	}{
		{
			"empty",
			[]byte{},
			codes.InvalidArgument,
			EmptyMessage,
		},
		{
			"too short",
			[]byte{1, 2, 3, 4, 5, 6},
			codes.InvalidArgument,
			EmptyMessage,
		},
		{
			"initial bits not zero",
			pad2header([]byte{0xC0, 0xFF}),
			codes.InvalidArgument,
			EmptyMessage,
		},
		{
			"rfc5769 sample request",
			RFC5769SampleRequest,
			codes.OK,
			RFC5769SampleRequestMessage,
		},
	}

	for _, tc := range tests {
		m, err := ParseMessage(tc.b)
		code := errors.Code(err)
		if code != tc.code {
			t.Log(err)
			t.Errorf("[%s] Unexpected code; got:%+v want:%+v", tc.desc, code, tc.code)
		}

		if !reflect.DeepEqual(m, tc.m) {
			t.Errorf("[%s] wrong message")
			t.Log(m)
			t.Log(tc.m)
		}
	}
}
