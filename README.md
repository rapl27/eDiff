# eDiff
Rolling Hash Algorithm

Coomand line usage: ./main [OPTIONS]

Options:
	-oldFilename string
	      The file for which delta will be computed
	-newFilename string
	      The file for which delta will be computed
	-chunkSize int
	      The file chunk size (optional)

Delta output:
Chunk index | Operation (I:inserted M:modified R:removed U:unmodiffied) | Data
1 U
2 M "text"
3 R