package servicesimpl

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"github.com/niflheimdevs/backend/internal/domain/repositories/storage"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/pixiv/go-libjpeg/jpeg"
)

type FileService struct {
	FileStorage storage.FileStorage
}

func NewFileService(fileStorage storage.FileStorage) *FileService {
	return &FileService{
		FileStorage: fileStorage,
	}
}

func (fs *FileService) GetUserProfileName(userid int, wantHighQual bool) string {
	if !wantHighQual {
		return fmt.Sprintf("userprofile%d_low.webp", userid)
	} else {
		return fmt.Sprintf("userprofile%d_high.jpeg", userid)
	}

}

func (fs *FileService) UploadProfilePhoto(data []byte, userid int) string {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}
	imgReader := bytes.NewReader(data)

	img, _, err := image.Decode(imgReader)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{
				exceptions.FORMAT_NOT_SUPPORTED,
			},
		})
	}
	var webpBuffer, jpegBuffer bytes.Buffer
	resizedImg := resize.Resize(512, 512, img, resize.Lanczos2)

	err = webp.Encode(&webpBuffer, resizedImg, &webp.Options{Lossless: false, Quality: 85})

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
		})
	}
	outputName := fs.GetUserProfileName(userid, false)
	fs.FileStorage.StoreFile(webpBuffer.Bytes(), outputName)

	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)

	options := &jpeg.EncoderOptions{
		Quality:         85,
		ProgressiveMode: true,
	}

	err = jpeg.Encode(&jpegBuffer, rgbaImg, options)
	if err != nil {
		log.Println("jpeg", err)
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
		})
	}
	outputName = fs.GetUserProfileName(userid, true)
	fs.FileStorage.StoreFile(jpegBuffer.Bytes(), outputName)
	return outputName
}

func (fs *FileService) DeleteProfilePhoto(userid int) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	target := fs.GetUserProfileName(userid, false)
	err := fs.FileStorage.DeleteFile(target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}

	target = fs.GetUserProfileName(userid, true)

	err = fs.FileStorage.DeleteFile(target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}
}
