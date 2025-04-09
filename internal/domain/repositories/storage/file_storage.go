package storage

type FileStorage interface {
	StoreFile(data []byte, outputName string)

	DeleteFile(target string) error
}
