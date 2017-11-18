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

var (
	SomeMethod = 0x123

	MessageIntegrityInWrongPlaceBytes = []byte{
		0x00, 0x01, 0x00, 0x58, //   Request type and message length
		0x21, 0x12, 0xa4, 0x42, //   Magic cookie
		0xb7, 0xe7, 0xa7, 0x01, //}
		0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
		0xfa, 0x87, 0xdf, 0xae, //}
		0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
		0x53, 0x54, 0x55, 0x4e, //}
		0x20, 0x74, 0x65, 0x73, //}  User-agent...
		0x74, 0x20, 0x63, 0x6c, //}  ...name
		0x69, 0x65, 0x6e, 0x74, //}
		0x00, 0x06, 0x00, 0x09, //   USERNAME attribute header
		0x65, 0x76, 0x74, 0x6a, //}
		0x3a, 0x68, 0x36, 0x76, //}  Username (9 bytes) and padding (3 bytes)
		0x59, 0x20, 0x20, 0x20, //}
		0x00, 0x08, 0x00, 0x14, //   MESSAGE-INTEGRITY attribute header
		0x9a, 0xea, 0xa7, 0x0c, //}
		0xbf, 0xd8, 0xcb, 0x56, //}
		0x78, 0x1e, 0xf2, 0xb5, //}  HMAC-SHA1 fingerprint
		0xb2, 0xd3, 0xf2, 0x49, //}
		0xc1, 0xb5, 0x71, 0xa2, //}
		0x00, 0x24, 0x00, 0x04, //   PRIORITY attribute header
		0x6e, 0x00, 0x01, 0xff, //   ICE priority value
		0x80, 0x29, 0x00, 0x08, //   ICE-CONTROLLED attribute header
		0x93, 0x2f, 0xf9, 0xb1, //}  Pseudo-random tie breaker...
		0x51, 0x26, 0x3b, 0x36, //}   ...for ICE control
		0x80, 0x28, 0x00, 0x04, //   FINGERPRINT attribute header
		0xe5, 0x7a, 0x3b, 0xcf, //   CRC32 fingerprint
	}

	MessageIntegrityInWrongPlaceMessage = MakeMessage(
		MessageType(0x0001),
		TransactionID{
			0xb7e7a701,
			0xbc34d686,
			0xfa87dfae,
		},
		[]Attribute{
			SoftwareAttribute{
				"STUN test client",
			},
			UsernameAttribute{
				"evtj:h6vY",
			},
			MessageIntegrityAttribute{
				HMAC: []byte{
					0x9a, 0xea, 0xa7, 0x0c,
					0xbf, 0xd8, 0xcb, 0x56,
					0x78, 0x1e, 0xf2, 0xb5,
					0xb2, 0xd3, 0xf2, 0x49,
					0xc1, 0xb5, 0x71, 0xa2,
				},
				raw: []byte{
					0x00, 0x01, 0x00, 0x58, //   Request type and message length
					0x21, 0x12, 0xa4, 0x42, //   Magic cookie
					0xb7, 0xe7, 0xa7, 0x01, //}
					0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
					0xfa, 0x87, 0xdf, 0xae, //}
					0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
					0x53, 0x54, 0x55, 0x4e, //}
					0x20, 0x74, 0x65, 0x73, //}  User-agent...
					0x74, 0x20, 0x63, 0x6c, //}  ...name
					0x69, 0x65, 0x6e, 0x74, //}
					0x00, 0x06, 0x00, 0x09, //   USERNAME attribute header
					0x65, 0x76, 0x74, 0x6a, //}
					0x3a, 0x68, 0x36, 0x76, //}  Username (9 bytes) and padding (3 bytes)
					0x59, 0x20, 0x20, 0x20, //}
				},
			},
		},
	)

	MessageIntegrityInWrongPlaceWithFingerprintBytes = []byte{
		0x00, 0x01, 0x00, 0x58, //   Request type and message length
		0x21, 0x12, 0xa4, 0x42, //   Magic cookie
		0xb7, 0xe7, 0xa7, 0x01, //}
		0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
		0xfa, 0x87, 0xdf, 0xae, //}
		0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
		0x53, 0x54, 0x55, 0x4e, //}
		0x20, 0x74, 0x65, 0x73, //}  User-agent...
		0x74, 0x20, 0x63, 0x6c, //}  ...name
		0x69, 0x65, 0x6e, 0x74, //}
		0x00, 0x08, 0x00, 0x14, //   MESSAGE-INTEGRITY attribute header
		0x9a, 0xea, 0xa7, 0x0c, //}
		0xbf, 0xd8, 0xcb, 0x56, //}
		0x78, 0x1e, 0xf2, 0xb5, //}  HMAC-SHA1 fingerprint
		0xb2, 0xd3, 0xf2, 0x49, //}
		0xc1, 0xb5, 0x71, 0xa2, //}
		0x80, 0x28, 0x00, 0x04, //   FINGERPRINT attribute header
		0xe5, 0x7a, 0x3b, 0xcf, //   CRC32 fingerprint
		0x00, 0x24, 0x00, 0x04, //   PRIORITY attribute header
		0x6e, 0x00, 0x01, 0xff, //   ICE priority value
		0x80, 0x29, 0x00, 0x08, //   ICE-CONTROLLED attribute header
		0x93, 0x2f, 0xf9, 0xb1, //}  Pseudo-random tie breaker...
		0x51, 0x26, 0x3b, 0x36, //}   ...for ICE control
		0x00, 0x06, 0x00, 0x09, //   USERNAME attribute header
		0x65, 0x76, 0x74, 0x6a, //}
		0x3a, 0x68, 0x36, 0x76, //}  Username (9 bytes) and padding (3 bytes)
		0x59, 0x20, 0x20, 0x20, //}
	}

	MessageIntegrityInWrongPlaceWithFingerprintMessage = MakeMessage(
		MessageType(0x0001),
		TransactionID{
			0xb7e7a701,
			0xbc34d686,
			0xfa87dfae,
		},
		[]Attribute{
			SoftwareAttribute{
				"STUN test client",
			},
			MessageIntegrityAttribute{
				HMAC: []byte{
					0x9a, 0xea, 0xa7, 0x0c,
					0xbf, 0xd8, 0xcb, 0x56,
					0x78, 0x1e, 0xf2, 0xb5,
					0xb2, 0xd3, 0xf2, 0x49,
					0xc1, 0xb5, 0x71, 0xa2,
				},
				raw: []byte{
					0x00, 0x01, 0x00, 0x58, //   Request type and message length
					0x21, 0x12, 0xa4, 0x42, //   Magic cookie
					0xb7, 0xe7, 0xa7, 0x01, //}
					0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
					0xfa, 0x87, 0xdf, 0xae, //}
					0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
					0x53, 0x54, 0x55, 0x4e, //}
					0x20, 0x74, 0x65, 0x73, //}  User-agent...
					0x74, 0x20, 0x63, 0x6c, //}  ...name
					0x69, 0x65, 0x6e, 0x74, //}
				},
			},
			FingerprintAttribute{
				0xe57a3bcf,
			},
		},
	)

	FingerprintNotLastBytes = []byte{
		0x00, 0x01, 0x00, 0x58, //   Request type and message length
		0x21, 0x12, 0xa4, 0x42, //   Magic cookie
		0xb7, 0xe7, 0xa7, 0x01, //}
		0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
		0xfa, 0x87, 0xdf, 0xae, //}
		0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
		0x53, 0x54, 0x55, 0x4e, //}
		0x20, 0x74, 0x65, 0x73, //}  User-agent...
		0x74, 0x20, 0x63, 0x6c, //}  ...name
		0x69, 0x65, 0x6e, 0x74, //}
		0x80, 0x28, 0x00, 0x04, //   FINGERPRINT attribute header
		0xe5, 0x7a, 0x3b, 0xcf, //   CRC32 fingerprint
		0x00, 0x24, 0x00, 0x04, //   PRIORITY attribute header
		0x6e, 0x00, 0x01, 0xff, //   ICE priority value
		0x80, 0x29, 0x00, 0x08, //   ICE-CONTROLLED attribute header
		0x93, 0x2f, 0xf9, 0xb1, //}  Pseudo-random tie breaker...
		0x51, 0x26, 0x3b, 0x36, //}   ...for ICE control
		0x00, 0x06, 0x00, 0x09, //   USERNAME attribute header
		0x65, 0x76, 0x74, 0x6a, //}
		0x3a, 0x68, 0x36, 0x76, //}  Username (9 bytes) and padding (3 bytes)
		0x59, 0x20, 0x20, 0x20, //}
	}

	FingerprintNotLastMessage = MakeMessage(
		MessageType(0x0001),
		TransactionID{
			0xb7e7a701,
			0xbc34d686,
			0xfa87dfae,
		},
		[]Attribute{
			SoftwareAttribute{
				"STUN test client",
			},
			FingerprintAttribute{
				0xe57a3bcf,
			},
		},
	)
)

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
		{
			"modified rfc5769 sample request - message integrity in wrong place",
			MessageIntegrityInWrongPlaceBytes,
			codes.OK,
			MessageIntegrityInWrongPlaceMessage,
			1,
			MessageClassRequest,
		},
		{
			"modified rfc5769 sample request - message integrity (with fingerprint) in wrong place",
			MessageIntegrityInWrongPlaceWithFingerprintBytes,
			codes.OK,
			MessageIntegrityInWrongPlaceWithFingerprintMessage,
			1,
			MessageClassRequest,
		},
		{
			"modified rfc5769 sample request - fingerprint not last",
			FingerprintNotLastBytes,
			codes.InvalidArgument,
			FingerprintNotLastMessage,
			1,
			MessageClassRequest,
		},
	}

	for _, tc := range tests {
		mp := &MessageParser{DefaultRegistry}

		m, err := mp.Parse(tc.b)
		code := errors.Code(err)
		if code != tc.code {
			t.Log(err)
			t.Fatalf("[%s] Unexpected code; got:%+v want:%+v", tc.desc, code, tc.code)
		}

		if code != codes.OK {
			continue
		}

		if !reflect.DeepEqual(m, tc.m) {
			t.Log(m)
			t.Log(tc.m)
			t.Fatalf("[%s] wrong message", tc.desc)
		}

		if m.Type.Method() != tc.method {
			t.Fatalf("[%s] wrong method; got:%d want:%d", tc.desc, m.Type.Method(), tc.method)
		}

		if m.Type.Class() != tc.class {
			t.Fatalf("[%s] wrong class; got:%v want:%v", tc.desc, m.Type.Class(), tc.class)
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
