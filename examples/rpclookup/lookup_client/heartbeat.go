package main

import (
	"net"
	"time"
	"bufio"
	"encoding/gob"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HeartBeatWrite(conn net.Conn, done chan int) {
	body := "hello"
	seq := time.Now().UnixNano()
	logger.Debug("seq: %d", seq)

	buf := bufio.NewWriter(conn)
	enc := gob.NewEncoder(buf)
	err := codec.EncodeHeartBeat(enc, seq, body)
	if err != nil {
		logger.Error("heart beat encode failed! %s", err.Error())
		done <- 1
		return
	}
	buf.Flush()
}

func HeartBeatRead(dec *gob.Decoder) error {
	var msg string
	err := codec.DecodeBody(dec, &msg)
	if err != nil {
		logger.Error("body decode failed! %s", err.Error())
		return err
	}
	logger.Debug("heart beat, body:%s", string(msg))
	return nil
}