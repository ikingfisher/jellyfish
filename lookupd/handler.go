package lookupd

import (
	"net"
)

type Handler interface {
	HandleMessage(conn net.Conn, reqBuf []byte) error
}