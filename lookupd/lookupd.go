package lookupd

import (
	"github.com/wondywang/rpclookup/core/client"
	"github.com/wondywang/rpclookup/core/lg"
	"github.com/wondywang/rpclookup/core/util"
	"os"
	"net"
	"io"
	"time"
	"sync/atomic"
)

type Lookupd struct {
	logger *lg.Logger
	tcpListener *net.TCPListener
	clients map[int64]*client.Client
	clientIDSequence int64
	heartbeatInterval int64
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
		heartbeatInterval: 10,
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
			this.logger.Error("run net accept failed! %s", err.Error())
			this.Exit()
			os.Exit(1)
		}
	})

	this.waitGroup.Wrap(func() { this.Ticker() })
	return
}

func (this * Lookupd) lookupLoop() error {
	for {
		this.tcpListener.SetDeadline(time.Now().Add(1e9))
		tcpConn, err := this.tcpListener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			return err
		}

		clientID := atomic.AddInt64(&this.clientIDSequence, 1)
		this.logger.Debug("New client[%d] connected : %s", clientID, tcpConn.RemoteAddr().String())
		client, cerr := client.NewClient(this.logger, clientID, tcpConn)
		if cerr != nil {
			this.logger.Error("new client failed! %s", err.Error())
			continue
		}
		this.clients[clientID] = client
		this.waitGroup.Wrap(func() { this.IOLoop(client) })

		for _, client := range this.clients {
			this.logger.Debug("client pool id[%d] ip: %s", client.ID, client.Conn.RemoteAddr().String())
		}
	}
	return nil
}

func (this * Lookupd) IOLoop(client *client.Client) error {
	for {
		select {
		case <- client.ExitChan:
			this.logger.Warning("ExitChan effect. close client[%d]", client.ID)
			this.CloseClient(client)
			return nil
		default:
			this.logger.Trace("default. begin receive msg.")
			//do nothing
		}
	
		buf := make([]byte, 17)
		n, err := io.ReadFull(client.Conn, buf)
		if err != nil || n != len(buf) {
			this.logger.Error("conn receive header failed! %s", err.Error())
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			this.logger.Warning("ExitChan set. client[%d]", client.ID)
			client.ExitChan <- 1	//occur error! disconnect client.
			continue
		}

		protocolMagic := string(buf[:1])
		this.logger.Trace("protocolMagic:%s", protocolMagic)
		bodySize := util.BytesToInt64(buf[1:9])
		seq := util.BytesToInt64(buf[9:])
		this.logger.Debug("seq:%d, buf:%v", seq, buf)

		body := make([]byte, bodySize)
		n, err = io.ReadFull(client.Conn, body)
		if err != nil {
			this.logger.Error("seq:%d, conn receive body failed! %s", seq, err.Error())
			continue
		}

		switch protocolMagic {
		case "H":
			//todo heart beat logic
			err = this.HeartBeat(client, seq, body)
		case "D":
			//todo reveive data logic
			err := this.handler.HandleMessage(client.Conn, body)
			if err != nil {
				this.logger.Error("client[%d] handler error. %s", client.ID, err.Error())
				continue
			}
		default:
			this.logger.Error("unkown magic code:%s", protocolMagic)
			continue
		}
	}
	return nil
}

func (this *Lookupd) HeartBeat(client *client.Client, seq int64, body []byte) error {
	ipStr := client.Conn.RemoteAddr().String()
	this.logger.Debug("%d, %s say: %s", seq, ipStr, string(body))

	err := client.PushHeartBeat()
	if err != nil {
		this.logger.Error("push heart beat failed! ", err.Error())
		return err
	}
	return nil
}


func (this *Lookupd) Ticker() {
	ticker := time.NewTicker(time.Duration(this.heartbeatInterval) * time.Second)
	for {
		select {
		case <- ticker.C:
			{
				timestamp := time.Now().Unix()
				for _, client := range this.clients {
					if timestamp - client.HeartbeatTime > 3 * this.heartbeatInterval {
						this.logger.Warning("ticker ExitChan set. close client[%d]", client.ID)
						client.ExitChan <- 1	//occur error! disconnect client.
					}
				}
			}
		}
	}
	return
}

func (this *Lookupd) SetLogger(l *lg.Logger) {
	this.logger = l
}

func (this *Lookupd) CloseClient(client *client.Client) error {
	client.Exit()
	delete(this.clients, client.ID)
	this.logger.Warning("client[%d] close. clients num:%d", client.ID, len(this.clients))
	return nil
}

func (this *Lookupd) Exit() error {
	this.tcpListener.Close()
	for _, client := range this.clients {
		client.Close()
	}
	return nil
}