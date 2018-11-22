package main

import (
	"net"
	"time"
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

	var rsp codec.Response
	rsp.Cmd = req.Cmd
	rsp.Body = []byte(string("server: rsp from server."))
	// body, err := codec.RspEncode(rsp)
	// if err != nil {
	// 	logger.Error("request encode failed! %s", err.Error())
	// 	return err
	// }
	// bodySize := util.Int64ToBytes(int64(len(body)))

	seq := time.Now().UnixNano()
	// seq := util.Int64ToBytes(nanoTime)
	logger.Debug("seq: %d", seq)

	// buf := []byte("D")
	// buf = append(buf, bodySize...)
	// buf = append(buf, seq...)
	// buf = append(buf, body...)
	buf, err := codec.Encode(seq, rsp)
	if err != nil {
		logger.Error("respone encode failed! %s", err.Error())
		return err
	}

	_, err = conn.Write(buf)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			logger.Debug("write message timeout. %s", err.Error())
			return err
		}
		logger.Error("Error to send message. %s", err.Error())
		return err
	}

	return nil
}