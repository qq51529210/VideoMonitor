package sip

import (
	"errors"
)

var (
	errTXTimeout = errors.New("transaction timeout")
	errTXFinish  = errors.New("transaction finished")
)

type transaction interface {
	writeMessage(Conn, *Message) error
	Key() string
}
