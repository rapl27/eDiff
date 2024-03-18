package rollinghash

import (
	"hash"
)

const defaultBase = 7 // https://qr.ae/psrgIj

type RollingHash struct {
	hash     uint32
	winStart int64
	winEnd   int64
	winSize  int64
	window   []byte
}

func NewRollingHash(chunkSize int64) (hash.Hash, error) {
	return &RollingHash{
		hash:     0,
		winStart: 0,
		winEnd:   0,
		winSize:  chunkSize,
	}, nil
}

func (rh *RollingHash) Write(data []byte) (int, error) {
	for _, b := range data {
		rh.hash = rh.hash*uint32(defaultBase) + uint32(b)
	}
	rh.window = append(rh.window, data...)
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
	rh.winEnd = 0
	rh.window = []byte{}
}

func (rh *RollingHash) RollHash(b byte) {
	rh.hash -= uint32(rh.window[0])
	rh.hash /= defaultBase
	if len(rh.window) > 0 {
		rh.window = rh.window[1:]
	}

	rh.Write([]byte{b})
}

func (rh *RollingHash) Signature() uint32 {
	return rh.hash
}
