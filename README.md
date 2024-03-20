
## Rolling Hash Algorithm

#### Command Line Usage
./main [OPTIONS]

#### Options

- `-oldFilename string`: The original file
- `-newFilename string`: The modified file
- `-chunkSize int`: The file chunk size (optional)

#### Delta Output

| Chunk Index | Operation (I:inserted M:modified R:removed U:unmodified) | Data   |
|-------------|-----------------------------------------------------------|--------|
| 1           | U                                                         |        |
| 2           | M                                                         | "text" |
| 3           | R                                                         |        |
