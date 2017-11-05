package handlers

import (
	"github.com/gopherx/net/nat"
)

type Namer interface {
	Name() string
}

func newSoftwareAttribute(n Namer) nat.SoftwareAttribute {
	return nat.SoftwareAttribute{n.Name()}
}
