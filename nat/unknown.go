package nat

import (
	"github.com/gopherx/base/binary/write"
)

const (
	UnknownAttributeRfcName string = "UNKNOWN"
)

func (u UnknownAttribute) Print(w *write.BigEndian) error {
	WriteTLVHeader(w, u.RawType, uint16(len(u.Bytes)))
	w.Bytes(u.Bytes)
	WriteTLVPadding(w, uint16(len(u.Bytes)))

	return nil
}

type UnknownAttribute struct {
	RawType AttributeType
	Bytes   []byte
}

func (a UnknownAttribute) Type() AttributeType {
	return a.RawType
}
