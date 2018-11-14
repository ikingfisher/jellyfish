package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"core/util"
)

func main() {
	conn, err := net.Dial("tcp", "10.100.71.218:16688")
	if err != nil {
		fmt.Println("connecting failed!", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("connecting sucess.", conn.RemoteAddr().String())
	HeartBeat(conn)
}

func HeartBeat(conn net.Conn) {
	done := make(chan string)
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
			case t := <- ticker.C:
			{
				fmt.Println(t, " client heart beat.")
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
		fmt.Println("Error to send message. ", err.Error())
		return
	}
	done <- "Send"
}

func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message. ", err.Error())
		return
	}
	fmt.Println(string(buf[:len-1]))
	done <- "Receive"
}