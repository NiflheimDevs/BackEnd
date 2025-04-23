package servicesimpl

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"time"

	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/repositories/storage"
	"github.com/niflheimdevs/backend/pkg"
)

type FileService struct {
	S3Storage storage.S3Storage
}

func NewFileService(s3storage storage.S3Storage) *FileService {
	return &FileService{
		S3Storage: s3storage,
	}
}

func (fs *FileService) GetUserProfileName(userid int, wantHighQual bool) string {
	if !wantHighQual {
		return fmt.Sprintf("userprofile%d_low.webp", userid)
	} else {
		return fmt.Sprintf("userprofile%d_high.jpeg", userid)
	}

}

func (fs *FileService) GetProfilePhotoURL(userid int, wantHighQual bool) string {
	return fs.S3Storage.GetPresignedURL(enums.ProfilePic, fs.GetUserProfileName(userid, wantHighQual), 8*time.Hour)
}

func (fs *FileService) GetResumeURL(userid int) string {
	outputName := fmt.Sprintf("resume%d.pdf", userid)
	return fs.S3Storage.GetPresignedURL(enums.Resume, outputName, 8*time.Hour)
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
	resizedImg := pkg.Resize(512, 512, img)

	err = pkg.ImageEncode(&webpBuffer, resizedImg, 1)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
		})
	}
	outputName := fs.GetUserProfileName(userid, false)
	fs.S3Storage.UploadObject(enums.ProfilePic, outputName, webpBuffer.Bytes())

	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)

	err = pkg.ImageEncode(&jpegBuffer, rgbaImg, 2)
	if err != nil {
		log.Println("jpeg", err)
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
		})
	}
	outputName = fs.GetUserProfileName(userid, true)
	fs.S3Storage.UploadObject(enums.ProfilePic, outputName, jpegBuffer.Bytes())
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
	err := fs.S3Storage.DeleteObject(enums.ProfilePic, target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}

	target = fs.GetUserProfileName(userid, true)
	err = fs.S3Storage.DeleteObject(enums.ProfilePic, target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}
}

func (fs *FileService) UploadResume(data []byte, userid int) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	outputName := fmt.Sprintf("resume%d.pdf", userid)

	fs.S3Storage.UploadObject(enums.Resume, outputName, data)
}

func (fs *FileService) DeleteResume(userid int) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	outputName := fmt.Sprintf("resume%d.pdf", userid)

	err := fs.S3Storage.DeleteObject(enums.Resume, outputName)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}
}
