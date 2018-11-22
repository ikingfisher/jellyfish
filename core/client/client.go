package client

import (
	"github.com/ikingfisher/jellyfish/core/codec"
	"github.com/ikingfisher/jellyfish/core/lg"
	"net"
	"time"
	"encoding/gob"
	"bufio"
)

type Client struct {
	// ctx *context
	logger *lg.Logger
	ID int64
	net.Conn
	Dec *gob.Decoder
	Enc *gob.Encoder
	EncBuf *bufio.Writer

	HeartbeatTime int64
	ExitChan chan int
}

func NewClient(logger *lg.Logger, id int64, conn net.Conn) (*Client, error){
	buf := bufio.NewWriter(conn)
	c := &Client{
		logger: logger,
		ID: id,
		Conn: conn,
		Dec: gob.NewDecoder(conn),
		Enc: gob.NewEncoder(buf),
		EncBuf: buf,
		ExitChan: make(chan int, 1),
	}
	return c, nil
}

func (this * Client) PushHeartBeat() error {
	ipStr := this.Conn.RemoteAddr().String()
	this.logger.Debug("client[%d] push heartbeat, remote ip: %s, last time: %d",
		this.ID, ipStr, this.HeartbeatTime)

	body := "hello"
	seq := time.Now().UnixNano()
	this.logger.Debug("seq: %d", seq)
	err := codec.EncodeHeartBeat(this.Enc, seq, body)
	if err != nil {
		this.logger.Error("heart beat encode failed! %s", err.Error())
		return err
	}
	this.EncBuf.Flush()
	this.HeartbeatTime = time.Now().Unix()
	return nil
}

func (this *Client) Exit() error {
	this.logger.Warning("client id:%d HeartbeatTime:%d exit!", this.ID, this.HeartbeatTime)
	this.Conn.Close()
	close(this.ExitChan)
	return nil
}