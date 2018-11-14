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
	HeartbeatTime int64
}

func NewClient(logger *lg.Logger, id int64, conn net.Conn) (*Client, error){
	c := &Client{
		logger: logger,
		ID: id,
		Conn: conn,
	}
	return c, nil
}

func (this * Client) PushHeartBeat() error {
	ipStr := this.Conn.RemoteAddr().String()
	this.logger.Debug("client[%d] push heartbeat, remote ip: %s, last time: %d",
		this.ID, ipStr, this.HeartbeatTime)
	
	msg := time.Now().String() + ", push heartbeat."
	b := []byte(msg)
	_, err := this.Conn.Write(b)
	if err != nil {
		this.logger.Error("push heart beat failed! ", err.Error())
		return err
	}
	this.HeartbeatTime = time.Now().Unix()
	return nil
}

func (this *Client) Exit() error {
	this.Conn.Close()
	return nil
}