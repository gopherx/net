package nat

import (
	"fmt"

	"reflect"
	"testing"

	"github.com/gopherx/base/binary/format"
	"github.com/gopherx/base/binary/write"
)

func TestMessagePrinter(t *testing.T) {
	mp := &MessagePrinter{DefaultRegistry, 8}

	opts := &PrintOptions{
		MessageIntegrityKey: RFC5769SampleRequestPwd,
		Fingerprint:         true,
		InitialBufferSize:   32,
	}
	mb, err := mp.Print(RFC5769SampleRequestMessage, opts)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(mb, RFC5769SampleRequest) {
		fmt.Println("got:-------------------------------------")
		fmt.Println(len(mb))
		format.OutHex4C(mb)

		fmt.Println("want:-------------------------------------")
		fmt.Println(len(RFC5769SampleRequest))
		format.OutHex4C(RFC5769SampleRequest)
		t.Error("wrong; result")
	}
}

func TestWriteTLVHeader(t *testing.T) {
	b := []byte{0, 0, 0, 0, 0, 0}
	w := &write.BigEndian{b, 1, nil}
	WriteTLVHeader(w, SoftwareAttributeType, 0xF00D)
	want := []byte{0, 0x80, 0x22, 0xF0, 0x0D, 0}

	if !reflect.DeepEqual(b, want) {
		t.Fatalf("got:%+v, want:%+v", b, want)
	}
}

func TestWriteTLVPadding(t *testing.T) {
	tests := []struct {
		written uint16
		want    []byte
	}{
		{1, []byte{0xFF, 0xFF, 0x20, 0x20, 0x20, 0xFF, 0xFF, 0xFF}},
		{2, []byte{0xFF, 0xFF, 0x20, 0x20, 0xFF, 0xFF, 0xFF, 0xFF}},
		{3, []byte{0xFF, 0xFF, 0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{4, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{7, []byte{0xFF, 0xFF, 0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for i, tc := range tests {
		b := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
		w := &write.BigEndian{b, 2, nil}
		WriteTLVPadding(w, tc.written)

		if !reflect.DeepEqual(b, tc.want) {
			t.Fatalf("%d got:%+v want:%+v", i, b, tc.want)
		}
	}
}

func TestCheckWriteLen(t *testing.T) {
	msg := MakeMessage(
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
			PriorityAttribute{
				0x6e0001ff,
			},
			IceControlledAttribute{
				0x932ff9b151263b36,
			},
			UsernameAttribute{
				"evtj:h6vY",
			},
		},
	)

	mp := &MessagePrinter{DefaultRegistry, 8}
	bytes, err := mp.Print(msg, nil)
	if err != nil {
		t.Fatal("Print failed; err:", err)
	}

	for i := 0; i < 8; i++ {
		t.Logf("%x", bytes[i])
	}
}
