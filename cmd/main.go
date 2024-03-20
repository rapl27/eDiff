package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rapl27/eDiff/delta"
	hash "github.com/rapl27/eDiff/rollinghash"
)

const (
	defaultChunkSize = 1024
)

/*
Usage: ./main [OPTIONS]

Options:

	-filename string
	      The file for which delta will be computed
	-chunkSize int
	      The file chunk size (optional)
*/

func main() {
	// Parse command arguments
	oldFilename := flag.String("oldFilename", "", "The original file")
	newFilename := flag.String("newFilename", "", "The modified file")
	chunkSize := flag.Int64("chunkSize", defaultChunkSize, "The file chunk size (optional)")

	flag.Parse()
	if *oldFilename == "" || *newFilename == "" {
		fmt.Println("Error: filename is required.")
		return
	}

	fileInfo, err := os.Stat(*oldFilename)
	if err != nil {
		fmt.Println("Error: Failed to get file size: ", err)
		return
	}
	fileSize := fileInfo.Size()

	// Validate chunk size
	if *chunkSize >= fileSize {
		fmt.Printf("Resizing chunk size to: %v", fileSize/2)
		*chunkSize = fileSize / 2
	}

	// Compute delta
	delta := runDelta(*oldFilename, *newFilename, *chunkSize)

	// Print delta and write to file
	deltaFile, err := os.Create("output.diff")
	if err != nil {
		fmt.Println("Error: Failed to create delta file: ", err)
		return
	}
	fmt.Println("Delta:")
	for _, d := range delta {
		txt := fmt.Sprintf("%v | % v | %v", d.Offset, d.Operation, string(d.Data))
		fmt.Println(txt)
		_, err := fmt.Fprintln(deltaFile, txt)
		if err != nil {
			fmt.Println("Error: Failed to write delta file: ", err)
			return
		}
	}
}

func runDelta(oldFilename, newFilename string, chunkSize int64) []delta.Delta {
	// Generate chunks signatures for the original file
	rh := hash.NewRollingHash(chunkSize).(*hash.RollingHash)

	oldFile, err := os.Open(oldFilename)
	if err != nil {
		fmt.Printf("Error: Failed to open file [%s]: %v", oldFilename, err)
		return []delta.Delta{}
	}
	defer oldFile.Close()

	chunk := make([]byte, chunkSize)
	var oldSigs []uint32
	for {
		n, err := oldFile.Read(chunk)
		if err != nil && err != io.EOF {
			fmt.Printf("Error: Failed to read chunk: %v", err)
			return []delta.Delta{}
		}

		if n == 0 {
			break
		}

		rh.Reset()
		rh.Write(chunk)
		oldSigs = append(oldSigs, rh.Signature())
	}

	// Compute delta
	rd, err := delta.NewRollingDiffer(chunkSize)
	if err != nil {
		fmt.Printf("Error: Failed to create differ instance: %v", err)
		return []delta.Delta{}
	}

	delta, err := rd.Delta(oldSigs, newFilename)
	if err != nil {
		fmt.Printf("Error: Failed to compute file delta: %v", err)
	}

	return delta
}
