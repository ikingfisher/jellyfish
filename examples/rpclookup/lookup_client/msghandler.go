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

func HandleMsgRead(conn net.Conn, done chan int) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			logger.Debug("read null message. %s", err.Error())
			return
		}
		done <- 1
		logger.Error("Error to read message. %s", err.Error())
		return
	}
	logger.Debug(string(buf[:len-1]))
}