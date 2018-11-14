package main

import (
	"net"
	"bufio"
	"core/lg"
)

type MsgHandler struct {
	logger *lg.Logger
}

func (this MsgHandler) HandleMessage(conn net.Conn, reader *bufio.Reader) error {
	message, err := reader.ReadString('\n')
	if err != nil {
		this.logger.Error("read from client failed! %s", err.Error())
		return err
	}

	this.logger.Debug("MsgHandler receive: %s", string(message))

	msg := "server: rsp to client.\n"
	b := []byte(msg)
	conn.Write(b)
	return nil
}