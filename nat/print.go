package nat

import (
	"fmt"

	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"

	"github.com/golang/glog"
)

var (
	DefaultInitialBufferSize = 512
	DefaultPrinter           = &MessagePrinter{DefaultRegistry, DefaultInitialBufferSize}
	DefaultPrintOptions      = &PrintOptions{nil, false, DefaultInitialBufferSize}
)

// AttributePrinterFunc prints the attribute into the byte buffer.
type AttributePrinterFunc func(w *write.BigEndian, a Attribute) error

type MessagePrinter struct {
	Registry          AttributeRegistry
	InitialBufferSize int
}

type PrintOptions struct {
	MessageIntegrityKey []byte
	Fingerprint         bool
	InitialBufferSize   int
}

func Print(m Message, opts *PrintOptions) ([]byte, error) {
	return DefaultPrinter.Print(m, opts)
}

func (p *MessagePrinter) Print(m Message, opts *PrintOptions) ([]byte, error) {
	if opts == nil {
		opts = DefaultPrintOptions
	}

	size := p.InitialBufferSize
	if opts.InitialBufferSize > size {
		size = opts.InitialBufferSize
	}
	if size == 0 {
		size = DefaultInitialBufferSize
	}

	used := uint16(0)
	for i := 0; i < 5; i++ {
		bytes := make([]byte, size)
		n, err := p.writeMsg(bytes, m, opts)
		used += n
		if err != nil {
			if _, tooSmall := err.(BufferTooSmallError); !tooSmall {
				return nil, err
			}

			newSize := size * 2
			if v := glog.V(11); v {
				v.Infof("buffer too small; %d -> %d", size, newSize)
			}
			size = newSize
			continue
		}

		if v := glog.V(11); v {
			v.Infof("buffer use: %d/%d", used, size)
		}
		return bytes[0:used], nil
	}
	return nil, nil
}

type BufferTooSmallError struct {
	w *write.BigEndian
}

func (b BufferTooSmallError) Error() string {
	return fmt.Sprintf("buffer too small; len:%d, Err:%+v", len(b.w.Dest), b.w.Err)
}

func (p *MessagePrinter) writeMsg(b []byte, m Message, opts *PrintOptions) (uint16, error) {
	w := &write.BigEndian{b, 0, nil}
	w.Uint16(MessageTypeMask & uint16(m.Type))
	w.Uint16(0xABCD)
	w.Uint32(MagicCookie)
	w.Uint32(m.TID.p0)
	w.Uint32(m.TID.p1)
	w.Uint32(m.TID.p2)

	for _, a := range m.Attrs {
		_, ok := p.Registry[a.Type]
		if !ok {
			//...unknown attribute...
			panic("handle unknown attributes")
			continue
		}

		a, _ := m.Attr(a.Type)
		if _, isMIA := a.(MessageIntegrityAttribute); isMIA {
			//...ignore any message integrity attribute!
			glog.Error("MessageIntegrityAttribute should not be added to Message; use PrintOptions instead. Ignoring")
			continue
		}

		if _, isFPA := a.(FingerprintAttribute); isFPA {
			//...ignore any fingerprint attribute!
			glog.Error("FingerprintAttribtue should not be added to Message; use PrintOptions instead. Ignoring")
			continue
		}

		err := a.Print(w)
		if err != nil {
			return 0, errors.Internal(err, "Failed to print attribute", a, w)
		}

		if w.Err != nil {
			return 0, BufferTooSmallError{w}
		}
	}

	if opts.MessageIntegrityKey != nil {
		PrintMessageIntegrityAttribute(w, opts.MessageIntegrityKey)
	}

	if opts.Fingerprint {
		//...fingerprint is calculated with the length already written; including the fingerprint.
		w.Uint16At(2, uint16(w.Offset)-HeaderSize+TLVHeaderSize+FingerprintAttributeSize)
		PrintFingerprintAttribute(w)
	} else {
		//...no fingerprint; just write the current size
		w.Uint16At(2, uint16(w.Offset)-HeaderSize+TLVHeaderSize)
	}

	if w.Err != nil {
		return 0, BufferTooSmallError{w}
	}

	return uint16(w.Offset), nil
}

func WriteTLVHeader(w *write.BigEndian, at AttributeType, l uint16) {
	w.Uint16(uint16(at))
	w.Uint16(l)
}

var (
	padding = [][]byte{
		{},
		{0x20},
		{0x20, 0x20},
		{0x20, 0x20, 0x20},
	}
)

func WriteTLVPadding(w *write.BigEndian, written uint16) uint16 {
	pad := uint16(0)
	if written%4 > 0 {
		pad = (4 - written%4)
	}

	pb := padding[pad]
	if len(pb) == 0 {
		return 0
	}

	w.Bytes(pb)
	return uint16(len(pb))
}
