package delta

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type rollingDiffer struct {
	hash      RollingHash
	chunkSize int64
}

type Chunk struct {
	Index     int64
	Signature uint32

	Operation string // todo: insert / modify / remove / unmodiffied
	Buf       []byte
}

type Delta []Chunk

func NewRollingDiffer(chunkSize int64) (Differ, error) {
	rh, err := NewRollingHash(chunkSize)
	if err != nil {
		fmt.Println("Error: Failed to create the rolling hash: %v", err)
		return nil, nil
	}

	return &rollingDiffer{
		hash:      rh,
		chunkSize: chunkSize,
	}, nil
}

func (rd *rollingDiffer) Delta(oldSigs []uint32, newFilename string) (Delta, error) {
	// Open file
	newFile, err := os.Open(newFilename)
	if err != nil {
		fmt.Println("Error: Failed to open file [%s]: %v", newFile, err)
		return nil, err
	}
	defer newFile.Close()

	// initialize delta with chunk indexes and signatures
	delta := initDelta(oldSigs)

	// Initialize rolling hash with the first chunck in input file
	reader := bufio.NewReader(newFile)
	buf := make([]byte, rd.chunkSize)
	n, err := io.ReadAtLeast(reader, buf, int(rd.chunkSize))
	if n < int(rd.chunkSize) {
		//todo: do something
	}
	rd.hash.Write(buf)

	// Search for a signature match function
	sigMatch := func(sig uint32, delta Delta) int {
		for i, d := range delta {
			if d.Signature == sig {
				return i
			}
		}

		return -1
	}

	//Todo:  carefull when reaching the last byte not to roll the window outside
	for {
		index := sigMatch(rd.rh.Signature(), delta)
		if index != -1 {
			// populate delta
			/*
					type Chunk struct {
					Index     int64
					Signature uint32

					Operation string // todo: insert / modify / remove / unmodiffied
					Buf       []byte
				}
			*/

			// mark unchanged chunks
			if delta[index].Operation == "" {
				delta[index].Operation = "s"
			}

			// mark removed chunks
			// rd.rh.winStart

			// mark removed chunks

			// mark inserted chunks

		}

		b, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			fmt.Println("Error: Failed to read byte: %v", err)
			return nil, err
		}

		if err == io.EOF {
			fmt.Println("End of file")
			break
		}

		rd.hash.RollHash(b)
	}
}

func initDelta(signatures []uint32) Delta {
	var delta Delta
	for i, sig := range signatures {
		delta = append(delta, Chunk{
			Index:     int64(i),
			Signature: sig,
		})
	}
	return delta
}
