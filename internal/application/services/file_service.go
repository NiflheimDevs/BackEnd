package services

type FileService interface {
	GetProfilePhotoURL(userid int, wantHighQual bool) string
	GetTeamProfilePhotoURL(teamid int64, wantHighQual bool) string
	GetResumeURL(userid int) string
	UploadProfilePhoto(data []byte, userid int)
	DeleteProfilePhoto(userid int)
	UploadTeamProfilePhoto(data []byte, teamid int64)
	DeleteTeamProfilePhoto(teamid int64)
	UploadResume(data []byte, userid int)
	DeleteResume(userid int)
}
