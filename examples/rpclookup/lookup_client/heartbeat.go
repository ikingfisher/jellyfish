package main

import (
	"net"
	"time"
	// "bufio"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HeartBeatWrite(conn net.Conn, done chan int) {
	body := "hello"
	seq := time.Now().UnixNano()
	logger.Debug("seq: %d", seq)

	err := codec.EncodeHeartBeat(conn, seq, body)
	if err != nil {
		logger.Error("heart beat encode failed! %s", err.Error())
		return
	}
}

func HeartBeatRead(conn net.Conn) error {
	var msg string
	err := codec.DecodeBody(conn, &msg)
	if err != nil {
		logger.Error("body decode failed! %s", err.Error())
		return err
	}
	logger.Debug("heart beat, body:%s", string(msg))
	return nil
}