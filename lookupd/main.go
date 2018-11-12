package main

import (
	"os"
	"fmt"
	"net"
	"flag"
	"syscall"
	"runtime"
	"path/filepath"
	"github.com/judwhite/go-svc/svc"
	"time"
	"bufio"
)

type program struct {
	tcpListener *net.TCPListener
	clients []*net.TCPConn
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		fmt.Println("program exit! err:%s", err.Error())
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	p.clients = make([]*net.TCPConn, 0)
	return nil
}

func (p *program) Start() error {
	var err error
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", ":16688")
	if err != nil {
		fmt.Println("resolve tcp addr failed! %s", err.Error())
		return err
	}

	p.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("listen tcp addr failed! %s", err.Error())
		return err
	}

	// err = p.Accept()
	// if err != nil {
	// 	fmt.Println("run net accept failed!", err.Error())
	// 	return err
	// }

	go p.Accept()
	go p.HeartBeat()

	fmt.Println("program start...")
	return nil
}

func (p *program) Stop() error {
	p.tcpListener.Close()
	for _, client := range p.clients {
		client.Close()
	}
	return nil
}

func (p * program) Accept() error {
	for {
		tcpConn, err := p.tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println("New client connected : " + tcpConn.RemoteAddr().String())
		p.clients = append(p.clients, tcpConn)
		go p.IOLoop(tcpConn)
	}
	return nil
}

func (p * program) IOLoop(conn *net.TCPConn) error {
	// ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read from client failed! %s", err.Error())
			return err
		}

		fmt.Println(string(message))

		msg := time.Now().String() + ", server: ni hao!\n"
		b := []byte(msg)
		conn.Write(b)
	}
	return nil
}

func (p * program) HeartBeat() error {
	fmt.Println("start server heart beat.")
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case t := <- ticker.C:
			{
				fmt.Println(t, " tick server heart beat.")
				for _, client := range p.clients {
					p.PushMsg(client)
				}
			}
		}
	}
	return nil
}

func (p * program) PushMsg(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	fmt.Println("remote client ip : " + ipStr)
	msg := time.Now().String() + ", push hart bit!\n"
	b := []byte(msg)
	conn.Write(b)
}