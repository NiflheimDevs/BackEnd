package services

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
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

func (fs *FileService) UploadProfilePhoto(data []byte, userid int) string {
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
	var webpBuffer, jpegBuffer bytes.Buffer
	resizedImg := resize.Resize(512, 512, img, resize.Lanczos2)

	err = webp.Encode(&webpBuffer, resizedImg, &webp.Options{Lossless: false, Quality: 85})
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.UNPROCESSABLE,
		})
	}
	outputName := fmt.Sprintf("userprofile%d_low.webp", userid)
	fs.FileStorage.StorageFile(webpBuffer.Bytes(), outputName)

	options := &jpeg.EncoderOptions{
		Quality:         85,
		ProgressiveMode: true,
	}

	err = jpeg.Encode(&jpegBuffer, img, options)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.UNPROCESSABLE,
		})
	}
	outputName = fmt.Sprintf("userprofile%d_high.jpeg", userid)
	fs.FileStorage.StorageFile(jpegBuffer.Bytes(), outputName)
	return outputName
}

func (fs *FileService) DeleteProfilePhoto(userid int) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	target := fmt.Sprintf("userprofile%d_low.webp", userid)
	err := fs.FileStorage.DeleteFile(target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.UNPROCESSABLE,
			Errors: []enums.SpecificError{enums.MISSING_FILE},
		})
	}

	target = fmt.Sprintf("userprofile%d_high.jpeg", userid)

	err = fs.FileStorage.DeleteFile(target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.UNPROCESSABLE,
			Errors: []enums.SpecificError{enums.MISSING_FILE},
		})
	}
}
