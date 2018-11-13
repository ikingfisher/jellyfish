package lookupd

import (
	"net"
	"time"
	"bufio"
	"sync/atomic"
	"core/client"
	"core/lg"
)

type Lookupd struct {
	logger *lg.Logger
	tcpListener *net.TCPListener
	clients map[int64]*client.Client
	clientIDSequence int64
	listenPort string
}

func NewLookupd(port string, logger *lg.Logger) (*Lookupd, error) {
	l := &Lookupd{
		logger:      logger,
		clients: 	 make(map[int64]*client.Client),
		listenPort:  port,
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

func (this * Lookupd) Accept() error {
	for {
		tcpConn, err := this.tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		clientID := atomic.AddInt64(&this.clientIDSequence, 1)
		this.logger.Debug("New client[%d] connected : %s", clientID, tcpConn.RemoteAddr().String())
		client := &client.Client{
			ID: clientID,
			Conn: tcpConn,
		}
		this.clients[clientID] = client
		go this.IOLoop(client)
	}
	return nil
}

func (this * Lookupd) IOLoop(client *client.Client) error {
	reader := bufio.NewReader(client.Conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			this.logger.Error("read from client failed! %s", err.Error())
			return err
		}

		this.logger.Debug(string(message))

		msg := time.Now().String() + ", server: ni hao!\n"
		b := []byte(msg)
		client.Conn.Write(b)
	}
	return nil
}

func (this * Lookupd) HeartBeat() error {
	this.logger.Debug("start server heart beat.")
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <- ticker.C:
			{
				this.logger.Debug("tick server heart beat.")
				for _, client := range this.clients {
					err := this.PushMsg(client)
					if err != nil {
						timestamp := time.Now().Unix()
						if timestamp - client.HeartbeatTime > 10 {
							client.Conn.Close()
							delete(this.clients, client.ID)
							this.logger.Debug("client[%d] close. clients count: %d", client.ID, len(this.clients))
						}
					}
				}
			}
		}
	}
	return nil
}

func (this * Lookupd) PushMsg(client *client.Client) error {
	ipStr := client.Conn.RemoteAddr().String()
	this.logger.Debug("client[%d] push heart beat, remote client ip: %s, last time: %d", client.ID, ipStr, client.HeartbeatTime)
	msg := time.Now().String() + ", push heart beat.\n"
	b := []byte(msg)
	_, err := client.Conn.Write(b)
	if err != nil {
		this.logger.Error("push heart beat failed! ", err.Error())
		return err
	}
	client.HeartbeatTime = time.Now().Unix()
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