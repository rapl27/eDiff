package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
		fmt.Println("Resizing chunk size to: %v", fileSize/2)
		*chunkSize = fileSize / 2
	}

	// Generate chunks signatures for the original file
	rh, err := hash.NewRollingHashh(*chunkSize)
	if err != nil {
		fmt.Println("Error: Failed to create the rolling hash: %v", err)
		return
	}
	oldFile, err := os.Open(*oldFilename)
	if err != nil {
		fmt.Println("Error: Failed to open file [%s]: %v", oldFilename, err)
		return
	}
	defer oldFile.Close()

	chunk := make([]byte, *chunkSize)
	var oldSigs []uint32
	for {
		n, err := oldFile.Read(chunk)
		if err != nil && err != io.EOF {
			fmt.Println("Error: Failed to read chunk: %v", err)
			return
		}

		if n == 0 {
			break
		}

		oldSigs = append(oldSigs, rh.Sign(chunk))
	}

	// Compute delta
	rd := delta.NewRollingDiffer(*chunkSize)
	delta, err := rd.Delta(*oldFilename, *newFilename)
	if err != nil {
		fmt.Println("Error: Failed to compute file delta: %v", err)
		return
	}

	// Print delta and write to file
	deltaFile, err := os.Create("output.diff")
	if err != nil {
		fmt.Println("Error: Failed to create delta file: ", err)
		return
	}
	fmt.Println("Delta:")
	for _, i := range delta {
		txt := fmt.Sprintf("%v | % v | %v", i.index, i.op, i.content)
		fmt.Println(txt)
		_, err := fmt.Fprintln(deltaFile, txt)
		if err != nil {
			fmt.Println("Error: Failed to write delta file: ", err)
			return
		}
	}
}
