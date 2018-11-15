package main

import (
	"errors"
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
					go HeartBeatWrite(conn, errorOccur)
					go HeartBeatRead(conn, errorOccur)
				}
			case <- errorOccur:
				logger.Warning("conn close. exit!")
				return errors.New("conn close.")
			}
		}
	})
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