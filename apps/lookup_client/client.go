package main

import (
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
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.SetOutput(os.Stdout)

	conn, err := net.Dial("tcp", "10.100.71.218:16688")
	if err != nil {
		logger.Fatal("connecting failed! %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	logger.Debug("connecting sucess. %s", conn.RemoteAddr().String())
	HeartBeat(conn)
}

func HeartBeat(conn net.Conn) {
	done := make(chan string)
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
		}
	}
}

func handleWrite(conn net.Conn, done chan string) {
	body := "hello"
	buf := []byte("H")
	bodySize := util.Int64ToBytes(int64(len(body)))
	buf = append(buf, bodySize...)
	buf = append(buf, body...)
	_, err := conn.Write(buf)
	if err != nil {
		logger.Error("Error to send message. %s", err.Error())
		return
	}
	done <- "Send"
}

func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		logger.Error("Error to read message. %s", err.Error())
		return
	}
	logger.Debug(string(buf[:len-1]))
	done <- "Receive"
}