package codec

import (
	"bytes"
	"encoding/gob"
	"unsafe"
)

type Header struct {
	T byte
	Seq uint64
	Size int32
}

type Request struct {
	Cmd  string
	Body []byte
}

type Response struct {
	Cmd  string
	Body []byte
}

func Encode(seq uint64, obj interface{}) ([]byte, error) {
	var buff bytes.Buffer

	var obj_buff bytes.Buffer
	enc := gob.NewEncoder(&obj_buff)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	var header_buff bytes.Buffer
	var header Header
	header.T = 'D'
	header.Seq = seq
	header.Size = len(obj_buff.Bytes())
	header_enc := gob.NewEncoder(&header_buff)
	err = header_enc.Encode(header)
	if err != nil {
		return nil, err
	}

	buff.Write(header_buff.Bytes())
	buff.Write(obj_buff.Bytes())
	buff.Write([]byte("E"))

	return buff.Bytes(), nil
}

func HeaderSize() int {
	return unsafe.Sizeof(Header{})
}

func Decode(body []byte, obj interface{}) error {
	buff := bytes.NewBuffer(body)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(obj)
	if err != nil {
		return err
	}
	return nil
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