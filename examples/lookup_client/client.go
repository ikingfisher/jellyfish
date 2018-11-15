package main

import (
	"github.com/wondywang/rpclookup/core/lg"
	"github.com/wondywang/rpclookup/core/util"
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

	conn, err := net.Dial("tcp", "10.100.71.218:16688")
	if err != nil {
		logger.Fatal("connecting failed! %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	logger.Debug("connecting sucess. %s", conn.RemoteAddr().String())

	var waitGroup util.WaitGroupWrapper
	waitGroup.Wrap(func() { HeartBeat(conn) })
	// waitGroup.Wrap(func() { ProcessMsg(conn) })
	ProcessMsg(conn)

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
			go HeartBeatWrite(conn, done)
			go HeartBeatRead(conn, done)
		}
		case <- done:
			logger.Warning("conn close. exit!")
			return
		}
	}
}

func ProcessMsg(conn net.Conn) {
	done := make(chan int)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
		{
			logger.Debug("client process msg.")
			go HandleMsgWrite(conn, done)
			go HandleMsgRead(conn, done)
		}
		case <- done:
			logger.Warning("conn close. exit!")
			return
		}
	}
}