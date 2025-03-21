package storage

import (
	"os"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type FileStorage struct {
	Constants *bootstrap.Constants
}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) SotorageFile(data []byte, outputName string) {

	outputFile, err := os.Create(fs.Constants.StorageDir + outputName)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
		})
	}
	defer outputFile.Close()

	_, err = outputFile.Write(data)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
		})
	}
}
