package quic

import (
	"sync"

	"github.com/lucas-clemente/quic-go/protocol"
)

var bufferPool sync.Pool

func getPacketBuffer() []byte {
	return bufferPool.Get().([]byte)
}

func putPacketBuffer(buf []byte) {
	if cap(buf) != int(protocol.MaxPacketSize) {
		println(buf)
		panic("putPacketBuffer called with packet of wrong size!")
	}
	bufferPool.Put(buf[:0])
}

func init() {
	bufferPool.New = func() interface{} {
		return make([]byte, 0, protocol.MaxPacketSize)
	}
}
