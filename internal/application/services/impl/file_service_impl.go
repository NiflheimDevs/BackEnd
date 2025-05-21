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
	"github.com/niflheimdevs/backend/internal/utils"
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

func (fs *FileService) GetProfilePhotoURL(userid int, wantHighQual bool) string {
	return fs.S3Storage.GetPresignedURL(enums.ProfilePic, fs.getUserProfileName(userid, wantHighQual), 24*time.Hour)
}

func (fs *FileService) GetTeamProfilePhotoURL(teamid int64, wantHighQual bool) string {
	return fs.S3Storage.GetPresignedURL(enums.TeamProfile, fs.getTeamProfileName(teamid, wantHighQual), 24*time.Hour)
}

func (fs *FileService) GetResumeURL(userid int) string {
	outputName := fmt.Sprintf("resume%d.pdf", userid)
	objects := fs.S3Storage.GetObjectList(enums.Resume)
	if utils.Contains(objects, outputName) {
		return fs.S3Storage.GetPresignedURL(enums.Resume, outputName, 24*time.Hour)
	} else {
		return ""
	}
}

func (fs *FileService) createHighQualPhoto(img image.Image) *bytes.Buffer {

	var jpegBuffer bytes.Buffer
	rgbaImg := image.NewRGBA(img.Bounds())
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)

	err := pkg.ImageEncode(&jpegBuffer, rgbaImg, 2)
	if err != nil {
		log.Println("jpeg", err)
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
		})
	}
	return &jpegBuffer
}

func (fs *FileService) UploadTeamProfilePhoto(data []byte, teamid int64) {
	img := fs.createImageInterface(data)

	webpBuffer := fs.createLowQualPhoto(img)
	jpegBuffer := fs.createHighQualPhoto(img)

	outputName := fs.getTeamProfileName(teamid, false)
	fs.S3Storage.UploadObject(enums.TeamProfile, outputName, webpBuffer.Bytes())

	outputName = fs.getTeamProfileName(teamid, true)
	fs.S3Storage.UploadObject(enums.TeamProfile, outputName, jpegBuffer.Bytes())
}

func (fs *FileService) UploadProfilePhoto(data []byte, userid int) {
	// ! move this if out
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	img := fs.createImageInterface(data)

	webpBuffer := fs.createLowQualPhoto(img)
	jpegBuffer := fs.createHighQualPhoto(img)

	outputName := fs.getUserProfileName(userid, false)
	fs.S3Storage.UploadObject(enums.ProfilePic, outputName, webpBuffer.Bytes())

	outputName = fs.getUserProfileName(userid, true)
	fs.S3Storage.UploadObject(enums.ProfilePic, outputName, jpegBuffer.Bytes())
}

func (fs *FileService) DeleteTeamProfilePhoto(teamid int64) {

	target := fs.getTeamProfileName(teamid, false)
	err := fs.S3Storage.DeleteObject(enums.TeamProfile, target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}

	target = fs.getTeamProfileName(teamid, true)
	err = fs.S3Storage.DeleteObject(enums.TeamProfile, target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}
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

	target := fs.getUserProfileName(userid, false)
	err := fs.S3Storage.DeleteObject(enums.ProfilePic, target)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.UNPROCESSABLE,
			Errors: []exceptions.SpecificError{exceptions.MISSING_FILE},
		})
	}

	target = fs.getUserProfileName(userid, true)
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

func (fs *FileService) getUserProfileName(userid int, wantHighQual bool) string {
	if !wantHighQual {
		return fmt.Sprintf("userprofile%d_low.webp", userid)
	} else {
		return fmt.Sprintf("userprofile%d_high.jpeg", userid)
	}
}

func (fs *FileService) getTeamProfileName(teamid int64, wantHighQual bool) string {
	if !wantHighQual {
		return fmt.Sprintf("teamprofile%d_low.webp", teamid)
	} else {
		return fmt.Sprintf("teamprofile%d_high.jpeg", teamid)
	}
}

func (fs *FileService) createImageInterface(data []byte) image.Image {
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
	return img
}

func (fs *FileService) createLowQualPhoto(img image.Image) *bytes.Buffer {
	var webpBuffer bytes.Buffer

	const threshold = 512

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newWidth, newHeight := width, height
	if width > threshold || height > threshold {
		newWidth, newHeight = threshold, threshold
	} else {
		newWidth, newHeight = width/2, height/2
	}

	resizedImg := pkg.Resize(uint(newWidth), uint(newHeight), img)

	err := pkg.ImageEncode(&webpBuffer, resizedImg, 1)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.UNPROCESSABLE,
		})
	}

	return &webpBuffer
}
