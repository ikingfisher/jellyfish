package main

import (
	"net"
	"time"
	"github.com/ikingfisher/jellyfish/core/codec"
)

func HeartBeatWrite(conn net.Conn, done chan int) {
	body := "hello"
	// buf := []byte("H")
	// bodySize := util.Int64ToBytes(int64(len(body)))
	seq := time.Now().UnixNano()
	// seq := util.Int64ToBytes(nanoTime)
	logger.Debug("seq: %d", seq)

	// buf = append(buf, bodySize...)
	// buf = append(buf, seq...)
	// buf = append(buf, body...)
	buf, err := codec.EncodeHeartBeat(seq, body)
	if err != nil {
		logger.Error("heart beat encode failed! %s", err.Error())
		return
	}

	logger.Info("%v", buf)

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

func HeartBeatRead(body []byte) error {
	var msg string
	err := codec.Decode(body, &msg)
	if err != nil {
		logger.Error("header decode failed! %s", err.Error())
		return err
	}
	logger.Debug("heart beat, body:%s", string(msg))
	return nil
}