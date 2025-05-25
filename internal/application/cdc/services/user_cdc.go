package cdcservices

type UserCdc interface {
	UserCapturer(data []byte) error
	UserTagCapturer(data []byte) error
}
