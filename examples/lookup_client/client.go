package main

import (
	"errors"
	"github.com/wondywang/rpclookup/core/lg"
	"github.com/wondywang/rpclookup/core/util"
	"io"
	"net"
	"os"
	"log"
	"time"
)

var logger *lg.Logger

func main() {
	logger = &lg.Logger{}
	logger.SetFlags(log.Ltime | log.Lshortfile)
	logger.SetOutput(os.Stdout)

	err := HeartBeat()
	if err != nil {
		logger.Warning("client disconnect, err: %s", err.Error())
	}
}

func HeartBeat() error{
	errorOccur := make(chan int)
	//retry 3 times at most, retry time gap start from 1 second
	return util.Retry(3, 1 * time.Second, func() error {
		conn, err := net.Dial("tcp", "127.0.0.1:16688")
		if err != nil {
			logger.Fatal("connecting failed! %s, retrying...", err.Error())
			return err
		}
		logger.Debug("connecting sucess. %s", conn.RemoteAddr().String())
		defer conn.Close()
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <- ticker.C:
				{
					logger.Debug("client heart beat.")
					go handleWrite(conn, errorOccur)
					go handleRead(conn, errorOccur)
				}
			case <- errorOccur:
				logger.Warning("conn close. exit!")
				return errors.New("conn close.")
			}
		}
	})
}

func handleWrite(conn net.Conn, errOccur chan int) {
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
		errOccur <- 1
		logger.Error("Error to send message. %s", err.Error())
		return
	}
}

func handleRead(conn net.Conn, errOccur chan int) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			logger.Debug("read null message. %s", err.Error())
			return
		}
		errOccur <- 1
		logger.Error("Error to read message. %s", err.Error())
		return
	}
	logger.Debug(string(buf[:len-1]))
}