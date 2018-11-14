package lookupd

import (
	"net"
	"bufio"
)

type Handler interface {
	HandleMessage(conn net.Conn, reader *bufio.Reader) error
}