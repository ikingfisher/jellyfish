package main

import (
	"net"
	"github.com/ikingfisher/jellyfish/core/lg"
	"github.com/ikingfisher/jellyfish/core/codec"
)

type MsgHandler struct {
	logger *lg.Logger
}

func (this MsgHandler) HandleMessage(conn net.Conn, reqBuf []byte) error {
	ipStr := conn.RemoteAddr().String()
	this.logger.Debug("MsgHandler receive from %s", ipStr)

	req, err := codec.RspDecode(reqBuf)
	if err != nil {
		this.logger.Error("decode failed! %s", err.Error())
		return err
	}

	this.logger.Debug("cmd:%s, body:%s", req.Cmd, string(req.Body))

	msg := "server: rsp to client.\n"
	b := []byte(msg)
	conn.Write(b)
	return nil
}