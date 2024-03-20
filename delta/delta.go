package delta

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/iancoleman/orderedmap"
	hash "github.com/rapl27/eDiff/rollinghash"
)

type RollingDiffer struct {
	hash      hash.RollingHash
	chunkSize int64
}

type Delta struct {
	Offset    int64
	Operation string // I:inserted M:modified R:removed U:unmodiffied
	Data      []byte
}

func NewRollingDiffer(chunkSize int64) (Differ, error) {
	rh := hash.NewRollingHash(chunkSize).(*hash.RollingHash)

	return &RollingDiffer{
		hash:      *rh,
		chunkSize: chunkSize,
	}, nil
}

func (rd *RollingDiffer) Delta(oldSigs []uint32, newFilename string) ([]Delta, error) {
	var delta []Delta

	// Open file
	newFile, err := os.Open(newFilename)
	if err != nil {
		fmt.Printf("Error: Failed to open file [%s]: %v", newFilename, err)
		return nil, err
	}
	defer newFile.Close()

	// Initialize chunk index to signature map
	chunksToSigMap := initSignaturesMap(oldSigs)

	// Initialize rolling hash with the first chunk in input file
	reader := bufio.NewReader(newFile)
	buf := make([]byte, rd.chunkSize)
	n, err := io.ReadAtLeast(reader, buf, int(rd.chunkSize))
	if err != nil || n < int(rd.chunkSize) {
		fmt.Printf("Error: Failed to read from [%s]: %v", newFilename, err)
		return nil, err
	}
	rd.hash.Write(buf)

	var unmatchedBytes []byte
	for {
		// Search for signature match
		matchIndex, unmatchIndexes := signatureMatch(rd.hash.Signature(), chunksToSigMap)
		if matchIndex != -1 {
			if len(unmatchedBytes) >= int(rd.chunkSize) {
				unmatchedBytes = unmatchedBytes[rd.chunkSize:]
			}

			// mark removed chunks
			if len(unmatchedBytes) == 0 && len(unmatchIndexes) != 0 {
				for _, i := range unmatchIndexes {
					chunksToSigMap.Delete(strconv.Itoa(i))
					delta = append(delta, Delta{
						Offset:    int64(i),
						Operation: "R",
					})
				}
			}

			// mark modified and inserted chunks
			for i := 0; i < len(unmatchedBytes); i += int(rd.chunkSize) {
				end := i + int(rd.chunkSize)
				if end > len(unmatchedBytes) {
					end = len(unmatchedBytes)
				}
				newChunk := unmatchedBytes[i:end]

				if len(unmatchIndexes) != 0 {
					delta = append(delta, Delta{
						Offset:    int64(unmatchIndexes[0]),
						Operation: "M",
						Data:      newChunk,
					})
					unmatchIndexes = unmatchIndexes[1:]
				} else {
					delta = append(delta, Delta{
						Offset:    int64(rd.hash.Offset()) / rd.chunkSize,
						Operation: "I",
						Data:      newChunk,
					})
				}
			}

			// mark unmodified chunks
			{
				chunksToSigMap.Delete(strconv.Itoa(matchIndex))
				delta = append(delta, Delta{
					Offset:    int64(matchIndex),
					Operation: "U",
				})
			}
			unmatchedBytes = []byte{}
		}

		// Read next byte
		b, err := reader.ReadByte()
		if err != nil && err != io.EOF {
			fmt.Printf("Error: Failed to read byte: %v", err)
			return nil, err
		}

		if err == io.EOF {
			fmt.Println("End of file")
			break
		}

		// Roll the hash
		byteOut := rd.hash.RollHash(b)
		unmatchedBytes = append(unmatchedBytes, byteOut)
	}

	return delta, nil
}

// Create map between chunk index and signature
func initSignaturesMap(signatures []uint32) *orderedmap.OrderedMap {
	m := orderedmap.New()
	for i, sig := range signatures {
		m.Set(strconv.Itoa(i+1), sig)
	}
	return m
}

// Search for signature match
// Return chunk index that matched the signature, and an index slice for chunks prior to the match
func signatureMatch(newSig uint32, oldSigs *orderedmap.OrderedMap) (int, []int) {
	var unmatchIndex []int
	for _, key := range oldSigs.Keys() {
		sig, _ := oldSigs.Get(key)
		if newSig == sig {
			index, _ := strconv.Atoi(key)
			return index, unmatchIndex
		} else {
			index, _ := strconv.Atoi(key)
			unmatchIndex = append(unmatchIndex, index)
		}
	}
	return -1, unmatchIndex
}
