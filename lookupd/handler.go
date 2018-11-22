package lookupd

import (
	// "net"
	"github.com/ikingfisher/jellyfish/core/client"
)

type Handler interface {
	HandleMessage(client *client.Client, reqBuf []byte) error
}