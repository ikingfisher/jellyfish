package main

import (
	// "io"
	"net"
	"time"
	"bufio"
	"encoding/gob"
	// "github.com/ikingfisher/jellyfish/core/util"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HandleMessage(conn net.Conn, done chan int) error {
	buf := make([]byte, codec.HeaderSize())
	var header codec.Header
	dec := gob.NewDecoder(conn)
	err := codec.DecodeHeader(dec, &header)
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
		err = HeartBeatRead(dec)
		if err != nil {
			logger.Error("heart beat error. %s", err.Error())
			return err
		}
	case "D":
		err := HandleMsgRead(dec)
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
	buf := bufio.NewWriter(conn)
	enc := gob.NewEncoder(buf)
	codec.Encode(enc, header)

	var req codec.Request
	req.Cmd = "HandleMsg"
	req.Body = append(req.Body, string("request from client")...)
	codec.Encode(enc, req)
	buf.Flush()
}

func HandleMsgRead(dec *gob.Decoder) error {
	var req codec.Request
	err := codec.DecodeBody(dec, &req)
	if err != nil {
		logger.Error("decode body failed! %s", err.Error())
		return err
	}

	logger.Debug("read body. cmd:%s, body:%s", req.Cmd, string(req.Body))
	return nil
}