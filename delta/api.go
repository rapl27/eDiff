package delta

type Differ interface {
	Delta(oldSigs []uint32, newFilename string) (Delta, error)
}
