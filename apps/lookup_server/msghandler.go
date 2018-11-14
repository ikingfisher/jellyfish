package main

import (
	"net"
	"core/lg"
)

type MsgHandler struct {
	logger *lg.Logger
}

func (this MsgHandler) HandleMessage(conn net.Conn, reqBuf []byte) error {
	ipStr := conn.RemoteAddr().String()
	this.logger.Debug("MsgHandler %s receive: %s", ipStr, string(reqBuf))
	msg := "server: rsp to client.\n"
	b := []byte(msg)
	conn.Write(b)
	return nil
}