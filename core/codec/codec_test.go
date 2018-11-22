package codec

import (
	"bytes"
	"testing"
	"encoding/gob"
)

func Test_EncodeHeartBeat(t *testing.T) {
	msg := "hello"
	seq := int64(1542862774224076350)
	buf, err := EncodeHeartBeat(seq, msg)
	if err != nil {
		t.Errorf("encode heart beat failed! %s", err.Error())
	}
	t.Logf("%v", buf)

	header_buff := buf[:HeaderSize()]
	t.Logf("%v", header_buff)

	var header Header
	err = Decode(header_buff, &header)
	if err != nil {
		t.Errorf("decode header failed! %s", err.Error())
	}

	t.Logf("header seq : %d, size : %d", header.Seq, header.Size)

	var obj_buff bytes.Buffer
	enc := gob.NewEncoder(&obj_buff)
	err = enc.Encode(msg)
	if err != nil {
		t.Errorf("encode body failed! %s", err.Error())
	}
	t.Logf("obj_buff: %v", obj_buff.Bytes())
}

func Test_Encode(t *testing.T) {
}

func Test_HeaderSize(t *testing.T) {
	size := HeaderSize()
	t.Logf("header size = %d", size)
}

// func Benchmark_Encode(b *testing.B) {
// 	t.Log("benchmark Encode...")
// }