package lookupd

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

func (this * Lookupd) ReqEncode(req Request) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(req)
    if err != nil {
		this.logger.Error("encode error: %s", err.Error())
		return nil, err
	}
	return buff.Bytes(), nil
}

func (this * Lookupd) RspDecode() (*Response, error) {
	var buff bytes.Buffer
	dec := gob.NewDecoder(&buff)
	var rsp Response
    err := dec.Decode(&rsp)
    if err != nil {
		this.logger.Error("decode error: %s", err.Error())
		return nil, err
	}
	return &rsp, nil
}