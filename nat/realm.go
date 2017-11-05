package nat

import (
	"fmt"

	"github.com/gopherx/base/binary/read"
	"github.com/gopherx/base/binary/write"
	"github.com/gopherx/base/errors"
)

const (
	RealmAttributeType     AttributeType = 0x0014
	RealmAttributeRfcName  string        = "REALM"
	RealmAttributeMaxChars uint16        = 128
	RealmAttributeMaxBytes uint16        = 763
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
	realm := RealmAttribute{}
	if l > RealmAttributeMaxBytes {
		return realm, errors.InvalidArgument(nil, fmt.Sprintf("too many bytes in realm; max=%d current=%d", RealmAttributeMaxBytes, l))
	}

	txt := string(r.Bytes(int(l)))
	if uint16(len(txt)) > RealmAttributeMaxChars {
		return realm, errors.InvalidArgument(nil, fmt.Sprintf("too many chars in realm; max=%d current=%d", RealmAttributeMaxChars, len(txt)))
	}

	realm.Realm = txt
	return realm, nil
}

func PrintRealmAttribute(w *write.BigEndian, a RealmAttribute) error {
	bytes := []byte(a.Realm)
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
