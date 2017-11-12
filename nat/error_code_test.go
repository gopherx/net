package nat

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

func TestErrorCode(t *testing.T) {
	rawErrorCode := []byte{
		0x00, 0x09, 0x00, 0x0B,
		0x00, 0x00, 0x03, 0x0A,
		'f', 'a', 'i', 'l',
		'e', 'd', '!', 0x20,
	}

	rawNoHeader := rawErrorCode[TLVHeaderSize:]

	// the error code reason is padded and that byte is not included in length.
	l := uint16(len(rawNoHeader) - 1)

	r := read.NewBigEndian(bytes.NewBuffer(rawNoHeader))

	a, err := ParseErrorCodeAttribute(r, l)
	if err != nil {
		t.Fatal(err)
	}

	want := ErrorCodeAttribute{3, 10, "failed!"}
	if !reflect.DeepEqual(a, want) {
		t.Fatalf("bad result; got:%+v want:%+v", a, want)
	}

	writer := &write.BigEndian{make([]byte, 1500), 0, nil}
	PrintErrorCodeAttribute(writer, want)
	rawReversed := writer.Dest[0:writer.Offset]
	if !reflect.DeepEqual(rawReversed, rawErrorCode) {
		t.Fatalf("bad result; got:%+v want:%+v", rawReversed, rawErrorCode)
	}
}
