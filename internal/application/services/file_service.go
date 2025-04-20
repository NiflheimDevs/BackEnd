package services

type FileService interface {
	GetUserProfileName(userid int, wantHighQual bool) string
	GetObjectURL(userid int, wantHighQual bool) string
	UploadProfilePhoto(data []byte, userid int) string
	DeleteProfilePhoto(userid int)
}
