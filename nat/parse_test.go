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
		desc   string
		b      []byte
		code   codes.Code
		m      Message
		method uint16
		class  MessageClass
	}{
		{
			"empty",
			[]byte{},
			codes.InvalidArgument,
			EmptyMessage,
			0,
			MessageClassRequest,
		},
		{
			"too short",
			[]byte{1, 2, 3, 4, 5, 6},
			codes.InvalidArgument,
			EmptyMessage,
			0,
			MessageClassRequest,
		},
		{
			"initial bits not zero",
			pad2header([]byte{0xC0, 0xFF}),
			codes.InvalidArgument,
			EmptyMessage,
			0,
			MessageClassRequest,
		},
		{
			"rfc5769 sample request",
			RFC5769SampleRequest,
			codes.OK,
			RFC5769SampleRequestMessage,
			1,
			MessageClassRequest,
		},
		{
			"rfc5769 sample request - long term creds",
			RFC5769SampleRequestWithLongTermCreds,
			codes.OK,
			RFC5769SampleRequestWithLongTermCredsMessage,
			1,
			MessageClassRequest,
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
			t.Errorf("[%s] wrong message", tc.desc)
			t.Log(m)
			t.Log(tc.m)
		}

		if m.Type.Method() != tc.method {
			t.Errorf("[%s] wrong method; got:%d want:%d", tc.desc, m.Type.Method(), tc.method)
		}

		if m.Type.Class() != tc.class {
			t.Errorf("[%s] wrong class; got:%v want:%v", tc.desc, m.Type.Class(), tc.class)
		}
	}
}

func TestMessageType(t *testing.T) {
	tests := []struct {
		desc   string
		method uint16
		class  MessageClass
		mt     MessageType
	}{
		{
			"bind-request",
			0x0001,
			MessageClassRequest,
			MessageType(0x0001),
		},
		{
			"bind-response-success",
			0x0001,
			MessageClassResponseSuccess,
			MessageType(0x0101),
		},
		{
			"allocate-request",
			0x0003,
			MessageClassRequest,
			MessageType(0x0003),
		},
		{
			"allocate-response-error",
			0x0003,
			MessageClassResponseError,
			MessageType(0x0113),
		},
		{
			"data-indication",
			0x0007,
			MessageClassIndication,
			MessageType(0x0017),
		},
	}

	for _, tc := range tests {
		mt := NewMessageType(tc.method, tc.class)
		if mt != tc.mt {
			t.Fatalf("[%s] wrong type; got:%v want:%v", tc.desc, mt, tc.mt)
		}

		method := mt.Method()
		if method != tc.method {
			t.Fatalf("[%s] wrong method; got:%v want:%v", tc.desc, method, tc.method)
		}

		class := mt.Class()
		if class != tc.class {
			t.Fatalf("[%s] wrong class; got:%v (0x%x) want:%v", tc.desc, class, uint16(class), tc.class)
		}
	}
}
