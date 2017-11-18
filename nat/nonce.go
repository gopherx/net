package nat

import (
	"crypto/rand"

	"github.com/golang/glog"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	NonceAttributeType    AttributeType = 0x0015
	NonceAttributeRfcName string        = "NONCE"
)

var (
	alphabet = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	}

	nonceDataChan = make(chan string, 1024)
)

func init() {
	go produceNonce()
	RegisterNonceAttribute(DefaultRegistry)
}

func RegisterNonceAttribute(r AttributeRegistry) {
	r.Register(
		NonceAttributeType,
		NonceAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseNonceAttribute(r, l)
		},
	)
}

func ParseNonceAttribute(r *read.BigEndian, l uint16) (NonceAttribute, error) {
	txt, err := Read127CharString(r, l)
	return NonceAttribute{txt}, err
}

func (a NonceAttribute) Print(w *write.BigEndian) error {
	bytes, err := Check127CharString(a.Nonce)
	if err != nil {
		return err
	}
	WriteTLVHeader(w, NonceAttributeType, uint16(len(bytes)))
	w.Bytes(bytes)
	WriteTLVPadding(w, uint16(len(bytes)))
	return nil
}

func produceNonce() {
	bytes := make([]byte, 127)
	tmp := bytes[63:]
	for {
		const maxRetries = 5
		i := 0
		for i = 0; i < maxRetries; i++ {
			_, err := rand.Read(tmp)
			if err != nil {
				glog.Error("failed to read crypto.Rand; err:", err)
				continue
			}

			at := 0
			for _, c := range tmp {
				bytes[at] = alphabet[c>>4]
				at++
				if at < len(bytes) {
					bytes[at] = alphabet[c&0x0F]
					at++
				}
			}

			nonceDataChan <- string(bytes)
			break
		}

		if i == maxRetries {
			glog.Fatal("failed to create nonce; terminating!!!")
		}
	}
}

// NewNonceAttribute creates a new NonceAttribute
func NewNonceAttribute() (NonceAttribute, error) {
	a := NonceAttribute{}

	a.Nonce = <-nonceDataChan

	return a, nil
}

type NonceAttribute struct {
	Nonce string
}

func (f NonceAttribute) Type() AttributeType {
	return NonceAttributeType
}
