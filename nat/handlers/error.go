package handlers

import (
	"github.com/gopherx/net/nat"
)

// NewServerErrorResponse returns a new Message with a ERROR-CODE attribute.
func NewServerErrorResponse(
	method uint16,
	reqID nat.TransactionID,
	err error,
	attrs ...nat.Attribute) nat.Message {

	return nat.NewErrorResponse(method, reqID, 5, 0, err.Error(), attrs...)
}
