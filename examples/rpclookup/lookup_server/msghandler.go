package main

import (
	"net"
	"time"
	"bufio"
	"github.com/ikingfisher/jellyfish/core/lg"
	"github.com/ikingfisher/jellyfish/core/codec"
	// "github.com/ikingfisher/jellyfish/core/util"
)

type MsgHandler struct {
	logger *lg.Logger
}

func (this MsgHandler) HandleMessage(conn net.Conn, reqBuf []byte) error {
	ipStr := conn.RemoteAddr().String()
	this.logger.Debug("MsgHandler receive from %s", ipStr)

	var req codec.Request
	err := codec.Decode(reqBuf, &req)
	if err != nil {
		this.logger.Error("decode failed! %s", err.Error())
		return err
	}
	this.logger.Debug("cmd:%s, body:%s", req.Cmd, string(req.Body))

	seq := time.Now().UnixNano()
	logger.Debug("seq: %d", seq)

	var header codec.Header
	header.T = 'D'
	header.Seq = time.Now().UnixNano()

	codec.Encode(conn, header)
	if err != nil {
		logger.Error("respone encode failed! %s", err.Error())
		return err
	}

	var rsp codec.Response
	rsp.Cmd = req.Cmd
	rsp.Body = []byte(string("server: rsp from server."))
	err = codec.Encode(conn, rsp)
	if err != nil {
		logger.Error("respone encode failed! %s", err.Error())
		return err
	}

	buf := bufio.NewWriter(conn)
	buf.Flush()
	return nil
}