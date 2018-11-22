package codec

import (
	"testing"
)

func Test_EncodeHeartBeat(t *testing.T) {
	msg := "hello"
	seq := int64(123)
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

	t.Logf("seq : %d", header.Seq)
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