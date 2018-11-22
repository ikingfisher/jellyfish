package codec

import (
	"bytes"
	"encoding/gob"
)

type Header struct {
	T byte
	Seq int64
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

func Encode(seq int64, obj interface{}) ([]byte, error) {
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
	header.Size = int32(len(obj_buff.Bytes()))
	header_enc := gob.NewEncoder(&header_buff)
	err = header_enc.Encode(header)
	if err != nil {
		return nil, err
	}

	buff.Write(header_buff.Bytes())
	buff.Write(obj_buff.Bytes())
	// buff.Write([]byte("E"))

	return buff.Bytes(), nil
}

func EncodeHeartBeat(seq int64, obj interface{}) ([]byte, error) {
	var buff bytes.Buffer

	var obj_buff bytes.Buffer
	enc := gob.NewEncoder(&obj_buff)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	var header_buff bytes.Buffer
	var header Header
	header.T = 'H'
	header.Seq = seq
	header.Size = int32(len(obj_buff.Bytes()))
	header_enc := gob.NewEncoder(&header_buff)
	err = header_enc.Encode(header)
	if err != nil {
		return nil, err
	}

	buff.Write(header_buff.Bytes())
	buff.Write(obj_buff.Bytes())
	// buff.Write([]byte("E"))

	return buff.Bytes(), nil
}

func HeaderSize() int {
	var buff bytes.Buffer
	var header Header
	header.T = 'H'
	header.Seq = int64(1542852055751046236)
	header.Size = 12348
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(header)
    if err != nil {
		return -1
	}
	return len(buff.Bytes())
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