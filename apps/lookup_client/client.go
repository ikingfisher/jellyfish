package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "10.100.71.218:16688")
	if err != nil {
		fmt.Println("connecting failed!", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("connecting sucess.", conn.RemoteAddr().String())

	done := make(chan string)
	go handleWrite(conn, done)
	go handleRead(conn, done)

	fmt.Println(<-done)
	fmt.Println(<-done)
}

func handleWrite(conn net.Conn, done chan string) {
	_, err := conn.Write([]byte("client: hello.\n"))
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