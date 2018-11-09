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

	for {
		tcpConn, err := p.tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println("New client connected : " + tcpConn.RemoteAddr().String())
		p.clients = append(p.clients, tcpConn)
		go tcpPipe(tcpConn)

		go func() {
			fmt.Println("set server hartbit.")
			select {
			case <- time.After(5 * time.Second):
				{
					for _, client := range p.clients {
						go pushMsg(client)
					}
				}
			}
		}()
	}

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