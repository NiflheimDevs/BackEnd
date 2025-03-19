package storage

import (
	"os"

	"github.com/niflheimdevs/backend/internal/exceptions"
)

type FileStorage struct {
}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) SotorageFile(data []byte, outputName string) error {
	outputFile, err := os.Create(outputName)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.SERVER_ERROR,
		})
	}
	defer outputFile.Close()

	_, err = outputFile.Write(outputBuffer.Bytes())
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.SERVER_ERROR,
		})
	}
}
