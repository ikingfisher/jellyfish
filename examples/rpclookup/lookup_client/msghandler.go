package main

import (
	// "io"
	"net"
	"time"
	"bufio"
	// "github.com/ikingfisher/jellyfish/core/util"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HandleMessage(conn net.Conn, done chan int) error {
	buf := make([]byte, codec.HeaderSize())
	var header codec.Header
	err := codec.DecodeHeader(conn, &header)
	if err != nil {
		logger.Error("header decode failed! %s", err.Error())
		done <- 1
		return err
	}
	protocolMagic := string(header.T)
	logger.Trace("protocolMagic:%s", protocolMagic)
	seq := header.Seq
	logger.Debug("seq:%d, buf:%v", seq, buf)

	switch protocolMagic {
	case "H":
		err = HeartBeatRead(conn)
		if err != nil {
			logger.Error("heart beat error. %s", err.Error())
			return err
		}
	case "D":
		err := HandleMsgRead(conn)
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
	var header codec.Header
	header.T = 'H'
	header.Seq = time.Now().UnixNano()

	codec.Encode(conn, header)

	var req codec.Request
	req.Cmd = "HandleMsg"
	req.Body = append(req.Body, string("request from client")...)
	codec.Encode(conn, req)
	buf := bufio.NewWriter(conn)
	buf.Flush()
}

func HandleMsgRead(conn net.Conn) error {
	var req codec.Request
	err := codec.DecodeBody(conn, &req)
	if err != nil {
		logger.Error("decode body failed! %s", err.Error())
		return err
	}

	logger.Debug("read body. cmd:%s, body:%s", req.Cmd, string(req.Body))
	return nil
}