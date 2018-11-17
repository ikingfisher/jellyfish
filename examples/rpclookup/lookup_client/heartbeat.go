package main

import (
	"net"
	"time"
	"github.com/ikingfisher/jellyfish/core/util"
)

func HeartBeatWrite(conn net.Conn, done chan int) {
	body := "hello"
	buf := []byte("H")
	bodySize := util.Int64ToBytes(int64(len(body)))
	nanoTime := time.Now().UnixNano()
	seq := util.Int64ToBytes(nanoTime)
	logger.Debug("seq: %d", nanoTime)

	buf = append(buf, bodySize...)
	buf = append(buf, seq...)
	buf = append(buf, body...)
	_, err := conn.Write(buf)
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
	logger.Debug("heart beat, body:%s", string(body))
	return nil
}