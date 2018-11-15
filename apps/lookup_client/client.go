package main

import (
	"io"
	"net"
	"os"
	"log"
	"time"
	"core/util"
	"core/lg"
)

var logger *lg.Logger

func main() {
	logger = &lg.Logger{}
	logger.SetFlags(log.Ltime | log.Lshortfile)
	logger.SetOutput(os.Stdout)

	conn, err := net.Dial("tcp", "127.0.0.1:16688")
	if err != nil {
		logger.Fatal("connecting failed! %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	logger.Debug("connecting sucess. %s", conn.RemoteAddr().String())
	HeartBeat(conn)
	logger.Warning("client disconnect. %s", conn.RemoteAddr().String())
}

func HeartBeat(conn net.Conn) {
	done := make(chan int)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
		{
			logger.Debug("client heart beat.")
			go handleWrite(conn, done)
			go handleRead(conn, done)
		}
		case <- done:
			logger.Warning("conn close. exit!")
			return
		}
	}
}

func handleWrite(conn net.Conn, done chan int) {
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

func handleRead(conn net.Conn, done chan int) {
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