package rollinghash

import (
	"hash"
	"math"
)

const defaultBase = 71 // https://qr.ae/psrgIj

type RollingHash struct {
	hash     uint32
	winStart int64
	winSize  int64
	window   []byte
}

func NewRollingHash(chunkSize int64) hash.Hash {
	return &RollingHash{
		hash:     0,
		winStart: 0,
		winSize:  chunkSize,
	}
}

func (rh *RollingHash) Write(data []byte) (int, error) {
	for _, b := range data {
		rh.hash += uint32(b) * uint32(math.Pow(defaultBase, float64(len(rh.window))))
		rh.window = append(rh.window, b)
	}
	return len(data), nil
}

func (rh *RollingHash) Sum(in []byte) []byte {
	v := rh.hash
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (rh *RollingHash) Size() int {
	return 4
}

func (rh *RollingHash) BlockSize() int {
	return 1
}

func (rh *RollingHash) Reset() {
	rh.hash = 0
	rh.winStart = 0
	rh.window = []byte{}
}

func (rh *RollingHash) RollHash(b byte) byte {
	if len(rh.window) < int(rh.winSize) {
		rh.winStart++
		rh.Write([]byte{b})
		return 0
	}

	byteOut := rh.window[0]
	rh.hash -= uint32(rh.window[0]) * uint32(math.Pow(defaultBase, 0))
	rh.hash /= defaultBase

	if len(rh.window) > 0 {
		rh.winStart++
		rh.window = rh.window[1:]
	} else {
		return 0
	}

	rh.Write([]byte{b})
	return byteOut
}

func (rh *RollingHash) Signature() uint32 {
	return rh.hash
}

func (rh *RollingHash) Offset() int64 {
	return rh.winStart
}
