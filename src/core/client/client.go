package client

import (
	"net"
	"time"
	"core/lg"
	// "context"
)

type Client struct {
	// ctx *context
	logger *lg.Logger
	ID int64
	net.Conn

	HeartbeatInterval int64
	HeartbeatTime int64
}

func NewClient(logger *lg.Logger, id int64, conn net.Conn, heartbeatInterval int64) (*Client, error){
	c := &Client{
		logger: logger,
		ID: id,
		Conn: conn,
		HeartbeatInterval: heartbeatInterval,
	}
	return c, nil
}

func (this *Client) Ticker() error {
	ticker := time.NewTicker(time.Duration(this.HeartbeatInterval) * time.Second)
	for {
		select {
		case <- ticker.C:
			{
				timestamp := time.Now().Unix()
				if timestamp - this.HeartbeatTime > 3 * this.HeartbeatInterval {
					this.Conn.Close()
				}
			}
		}
	}
}

func (this * Client) PushMsg(message  string) error {
	ipStr := this.Conn.RemoteAddr().String()
	this.logger.Debug("client[%d] push heart beat, remote client ip: %s, last time: %d", this.ID, ipStr, this.HeartbeatTime)
	b := []byte(message)
	_, err := this.Conn.Write(b)
	if err != nil {
		this.logger.Error("push heart beat failed! ", err.Error())
		return err
	}
	this.HeartbeatTime = time.Now().Unix()
	return nil
}