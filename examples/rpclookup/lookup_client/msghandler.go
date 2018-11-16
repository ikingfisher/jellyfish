package main

import (
	"io"
	"net"
	"time"
	"github.com/ikingfisher/jellyfish/core/util"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HandleMsgWrite(conn net.Conn, done chan int) {
	// body := "hello"
	buf := []byte("D")
	var req codec.Request
	req.Cmd = "QueryList"
	req.Body = append(req.Body, string("request from client")...)
	body, err := codec.ReqEncode(req)
	if err != nil {
		logger.Error("request encode failed! %s", err.Error())
		return
	}

	bodySize := util.Int64ToBytes(int64(len(body)))
	nanoTime := time.Now().UnixNano()
	seq := util.Int64ToBytes(nanoTime)
	logger.Debug("seq: %d", nanoTime)

	buf = append(buf, bodySize...)
	buf = append(buf, seq...)
	buf = append(buf, body...)

	_, err = conn.Write(buf)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			logger.Debug("write message timeout. %s", err.Error())
			return
		}
		done <- 1
		logger.Error("Error to send message. %s", err.Error())
		return
	}
}

func HandleMsgRead(conn net.Conn, done chan int) error {
	buf := make([]byte, 17)
	n, err := io.ReadFull(conn, buf)
	if err != nil || n != len(buf) {
		logger.Error("conn receive header failed! %s", err.Error())
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			return err
		}
		return err
	}
	protocolMagic := string(buf[:1])
	logger.Trace("protocolMagic:%s", protocolMagic)
	bodySize := util.BytesToInt64(buf[1:9])
	seq := util.BytesToInt64(buf[9:])
	logger.Debug("seq:%d, buf:%v", seq, buf)
	body := make([]byte, bodySize)
	n, err = io.ReadFull(conn, body)
	if err != nil {
		logger.Error("seq:%d, conn receive body failed! %s", seq, err.Error())
		return err
	}

	rsp, err := codec.RspDecode(body)
	if err != nil {
		logger.Error("decode failed! %s", err.Error())
		return err
	}

	logger.Debug("cmd:%s, body:%s", rsp.Cmd, string(rsp.Body))
	return nil
}