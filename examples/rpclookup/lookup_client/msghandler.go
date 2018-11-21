package main

import (
	"io"
	"net"
	"time"
	"github.com/ikingfisher/jellyfish/core/util"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HandleMessage(conn net.Conn, done chan int) error {
	buf := make([]byte, codec.HeaderSize())
	n, err := io.ReadFull(conn, buf)
	if err != nil || n != len(buf) {
		logger.Error("conn receive header failed! %s", err.Error())
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			done <- 1
			return err
		}
		done <- 1
		return err
	}
	var header codec.Header
	err = codec.Decode(buf, &header)
	if err != nil {
		logger.Error("header decode failed! %s", err.Error())
		done <- 1
		return err
	}
	protocolMagic := string(header.T)
	logger.Trace("protocolMagic:%s", protocolMagic)
	// bodySize := util.BytesToInt64(buf[1:9])
	seq := header.Seq
	logger.Debug("seq:%d, buf:%v", seq, buf)

	body := make([]byte, header.Size)
	n, err = io.ReadFull(conn, body)
	if err != nil {
		logger.Error("seq:%d, conn receive body failed! %s", seq, err.Error())
		done <- 1
		return err
	}

	switch protocolMagic {
	case "H":
		err = HeartBeatRead(body)
		if err != nil {
			logger.Error("heart beat error. %s", err.Error())
			return err
		}
	case "D":
		err := HandleMsgRead(body)
		if err != nil {
			logger.Error("mgs handler error. %s", err.Error())
			return err
		}
	default:
		logger.Error("unkown magic code:%s", protocolMagic)
		return nil
	}
	return nil
}

func HandleMsgWrite(conn net.Conn, done chan int) {
	buf := []byte("D")
	var req codec.Request
	req.Cmd = "HandleMsg"
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

func HandleMsgRead(body []byte) error {
	rsp, err := codec.RspDecode(body)
	if err != nil {
		logger.Error("decode failed! %s", err.Error())
		return err
	}

	logger.Debug("cmd:%s, body:%s", rsp.Cmd, string(rsp.Body))
	return nil
}