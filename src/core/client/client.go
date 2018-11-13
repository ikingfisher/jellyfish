package client

import (
	"net"
)

type Client struct {
	ID int64
	net.Conn

	HeartbeatTime int64
}
