package codec

import (
	"bytes"
	"encoding/gob"
)

type Request struct {
	Cmd  string
	Body []byte
}

type Response struct {
	Cmd  string
	Body []byte
}

func ReqEncode(req Request) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(req)
    if err != nil {
		// this.logger.Error("encode error: %s", err.Error())
		return nil, err
	}
	return buff.Bytes(), nil
}

func ReqDecode(body []byte) (*Request, error) {
	buff := bytes.NewBuffer(body)
	dec := gob.NewDecoder(buff)
	var req Request
    err := dec.Decode(&req)
    if err != nil {
		// this.logger.Error("decode error: %s", err.Error())
		return nil, err
	}
	return &req, nil
}

func RspEncode(rsp Response) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(rsp)
    if err != nil {
		// this.logger.Error("encode error: %s", err.Error())
		return nil, err
	}
	return buff.Bytes(), nil
}

func RspDecode(body []byte) (*Response, error) {
	buff := bytes.NewBuffer(body)
	dec := gob.NewDecoder(buff)
	var rsp Response
    err := dec.Decode(&rsp)
    if err != nil {
		// this.logger.Error("decode error: %s", err.Error())
		return nil, err
	}
	return &rsp, nil
}