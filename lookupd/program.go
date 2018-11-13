package main

import (
	"os"
	"net"
	"flag"
	"runtime"
	"path/filepath"
	"time"
	"bufio"
	"log"
	"sync/atomic"
	"github.com/judwhite/go-svc/svc"
)


type Program struct {
	tcpListener *net.TCPListener

	clients map[int64]*Client
	clientIDSequence int64
}

func (p *Program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	p.clients = make(map[int64]*Client)

	logger.SetFlags(log.Ldate | log.Lshortfile)
	logger.SetOutput(os.Stdout)
	return nil
}

func (p *Program) Start() error {
	var err error
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", ":16688")
	if err != nil {
		logger.Error("resolve tcp addr failed! %s", err.Error())
		return err
	}

	p.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Error("listen tcp addr failed! %s", err.Error())
		return err
	}

	go func() {
		err = p.Accept()
		if err != nil {
			logger.Error("run net accept failed!", err.Error())
			os.Exit(1)
		}
	}()

	go p.HeartBeat()

	logger.Debug("Program start...")
	return nil
}

func (p *Program) Stop() error {
	p.tcpListener.Close()
	for _, client := range p.clients {
		client.Close()
	}
	return nil
}

func (p * Program) Accept() error {
	for {
		tcpConn, err := p.tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		clientID := atomic.AddInt64(&p.clientIDSequence, 1)
		logger.Debug("New client[%d] connected : %s", clientID, tcpConn.RemoteAddr().String())
		client := &Client{
			Conn: tcpConn,
		}
		p.clients[clientID] = client
		go p.IOLoop(client)
	}
	return nil
}

func (p * Program) IOLoop(client *Client) error {
	// ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(client.Conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("read from client failed! %s", err.Error())
			return err
		}

		logger.Debug(string(message))

		msg := time.Now().String() + ", server: ni hao!\n"
		b := []byte(msg)
		client.Conn.Write(b)
	}
	return nil
}

func (p * Program) HeartBeat() error {
	logger.Debug("start server heart beat.")
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <- ticker.C:
			{
				logger.Debug("tick server heart beat.")
				for _, client := range p.clients {
					err := p.PushMsg(client)
					if err != nil {
						timestamp := time.Now().Unix()
						if client.HeartbeatTime - timestamp > 10 {
							client.Conn.Close()
							delete(p.clients, client.ID)
							logger.Debug("client[%d] close.", client.ID)
						}
					}
				}
			}
		}
	}
	return nil
}

func (p * Program) PushMsg(client *Client) error {
	ipStr := client.Conn.RemoteAddr().String()
	logger.Debug("remote client ip : " + ipStr)
	msg := time.Now().String() + ", push heart beat.\n"
	b := []byte(msg)
	_, err := client.Conn.Write(b)
	if err != nil {
		logger.Error("push heart beat failed! ", err.Error())
		return err
	}
	client.HeartbeatTime = time.Now().Unix()
	return nil
}