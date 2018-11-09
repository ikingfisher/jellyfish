package main

import (
	"fmt"
	"bufio"
	"time"
	"net"
)

func tcpPipe(conn *net.TCPConn) {
	// ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read from client failed! %s", err.Error())
			return
		}

		fmt.Println(string(message))
		msg := time.Now().String() + ", server: ni hao!\n"
		b := []byte(msg)
		conn.Write(b)
	}
}

func pushMsg(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	fmt.Println("remote client ip : " + ipStr)
	msg := time.Now().String() + ", push hart bit!\n"
	b := []byte(msg)
	conn.Write(b)
}