package services

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/repositories/storage"
	"github.com/pixiv/go-libjpeg/jpeg"
	_ "golang.org/x/image/webp"
)

type FileService struct {
	FileStorage *storage.FileStorage
}

func NewFileService(fileStorage *storage.FileStorage) *FileService {
	return &FileService{
		FileStorage: fileStorage,
	}
}

func (fs *FileService) UploadProfilePhoto(data []byte, userid int) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}
	imgReader := bytes.NewReader(data)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.UNPROCESSABLE,
			Errors: []enums.SpecificError{
				enums.FORMAT_NOT_SUPPORTED,
			},
		})
	}
	var outputBuffer bytes.Buffer

	options := &jpeg.EncoderOptions{
		Quality:         85,
		ProgressiveMode: true,
	}

	err = jpeg.Encode(&outputBuffer, img, options)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.UNPROCESSABLE,
		})
	}

}
