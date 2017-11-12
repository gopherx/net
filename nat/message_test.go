package nat

import (
	"bytes"
	"testing"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/errors"
	"github.com/gopherx/base/errors/codes"
)

func longString(c byte, n int) string {
	buff := make([]byte, n)
	for i := range buff {
		buff[i] = c
	}
	return string(buff)
}

func TestRead127CharString(t *testing.T) {
	tests := []struct {
		desc string
		txt  string
		l    uint16
		want string
		code codes.Code
	}{
		{
			"smoke test",
			"hello world",
			11,
			"hello world",
			codes.OK,
		},
		{
			"empty string",
			"",
			0,
			"",
			codes.OK,
		},
		{
			"too many chars",
			"0123456789012345678901234567890123456789012345678901234567890123456789" +
				"0123456789012345678901234567890123456789012345678901234567890123456789",
			140,
			"",
			codes.InvalidArgument,
		},
		{
			"too many byte",
			longString('a', 784),
			784,
			"",
			codes.InvalidArgument,
		},
	}

	for _, tc := range tests {
		r := read.NewBigEndian(bytes.NewBufferString(tc.txt))
		txt, err := Read127CharString(r, tc.l)
		if errors.Code(err) != tc.code {
			t.Errorf("[%s] wrong code; got:%+v want:%+v", tc.desc, errors.Code(err), tc.code)
		}

		if txt != tc.want {
			t.Errorf("[%s] wrong text; got:%q want:%q", tc.desc, txt, tc.want)
		}
	}
}
