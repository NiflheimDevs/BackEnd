package storage

import (
	"os"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type FileStorage struct {
	Constants *bootstrap.Constants
}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) StoreFile(data []byte, outputName string) {

	outputFile, err := os.Create(fs.Constants.StorageDir + outputName)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
	defer outputFile.Close()

	_, err = outputFile.Write(data)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
}

func (fs *FileStorage) DeleteFile(target string) error {
	err := os.Remove(fs.Constants.StorageDir + target)
	return err
}
