package cdcservices

type TeamCdc interface {
	TeamCapturer(data []byte) error
}
