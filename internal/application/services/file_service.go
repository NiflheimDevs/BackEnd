package services

type FileService interface {
	GetUserProfileName(userid int, wantHighQual bool) string
	GetProfilePhotoURL(userid int, wantHighQual bool) string
	GetResumeURL(userid int) string
	UploadProfilePhoto(data []byte, userid int) string
	DeleteProfilePhoto(userid int)
	UploadResume(data []byte, userid int)
	DeleteResume(userid int)
}
