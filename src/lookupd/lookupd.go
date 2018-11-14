package lookupd

import (
	"os"
	"net"
	"time"
	"bufio"
	"sync/atomic"
	"core/client"
	"core/lg"
	"core/util"
)

type Lookupd struct {
	logger *lg.Logger
	tcpListener *net.TCPListener
	clients map[int64]*client.Client
	clientIDSequence int64
	listenPort string
	handler Handler
	waitGroup util.WaitGroupWrapper
}

func NewLookupd(port string, logger *lg.Logger, handler Handler) (*Lookupd, error) {
	l := &Lookupd{
		logger:      logger,
		clients: 	 make(map[int64]*client.Client),
		listenPort:  port,
		handler:	 handler,
	}
	return l, nil
}

func (this * Lookupd) Init() error {
	var err error
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", ":" + this.listenPort)
	if err != nil {
		this.logger.Error("resolve tcp addr failed! %s", err.Error())
		return err
	}

	this.tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		this.logger.Error("listen tcp addr failed! %s", err.Error())
		return err
	}
	return nil
}

func (this *Lookupd) Main() {
	this.waitGroup.Wrap(func() {
		err := this.lookupLoop()
		if err != nil {
			this.logger.Error("run net accept failed!", err.Error())
			this.Exit()
			os.Exit(1)
		}
	})
	return
}

func (this * Lookupd) lookupLoop() error {
	for {
		tcpConn, err := this.tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		clientID := atomic.AddInt64(&this.clientIDSequence, 1)
		this.logger.Debug("New client[%d] connected : %s", clientID, tcpConn.RemoteAddr().String())
		client, cerr := client.NewClient(this.logger, clientID, tcpConn, 10)
		if cerr != nil {
			this.logger.Error("new client failed! %s", err.Error())
			continue
		}
		this.clients[clientID] = client
		go this.IOLoop(client)
	}
	return nil
}

func (this * Lookupd) IOLoop(client *client.Client) error {
	reader := bufio.NewReader(client.Conn)

	for {
		err := this.handler.HandleMessage(client.Conn, reader)
		if err != nil {
			this.logger.Error("client[%d] handler error. %s", client.ID, err.Error())
			continue
		}
		client.HeartbeatTime = time.Now().Unix()
	}
	return nil
}

func (this *Lookupd) SetLogger(l *lg.Logger) {
	this.logger = l
}

func (this *Lookupd) Exit() error {
	this.tcpListener.Close()
	for _, client := range this.clients {
		client.Close()
	}
	return nil
}