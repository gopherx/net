package nat

import (
	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
)

const (
	RealmAttributeType    AttributeType = 0x0014
	RealmAttributeRfcName string        = "REALM"
)

func init() {
	RegisterRealmAttribute(DefaultParser)
}

func RegisterRealmAttribute(p *MessageParser) {
	p.Register(
		RealmAttributeType,
		RealmAttributeRfcName,
		func(r *read.BigEndian, l uint16) (Attribute, error) {
			return ParseRealmAttribute(r, l)
		},
		func(w *write.BigEndian, a Attribute) error {
			return PrintRealmAttribute(w, a.(RealmAttribute))
		},
	)
}

func ParseRealmAttribute(r *read.BigEndian, l uint16) (RealmAttribute, error) {
	realm, err := Read127CharString(r, l)
	return RealmAttribute{realm}, err
}

func PrintRealmAttribute(w *write.BigEndian, a RealmAttribute) error {
	bytes, err := Check127CharString(a.Realm)
	if err != nil {
		return err
	}

	WriteTLVHeader(w, RealmAttributeType, uint16(len(bytes)))
	w.Bytes(bytes)
	WriteTLVPadding(w, uint16(len(bytes)))
	return nil
}

type RealmAttribute struct {
	Realm string
}

func (f RealmAttribute) Type() AttributeType {
	return RealmAttributeType
}
