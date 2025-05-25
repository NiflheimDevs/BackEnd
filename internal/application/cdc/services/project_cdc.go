package cdcservices

type ProjectCdc interface {
	ProjectCapturer(data []byte) error
	ProjectTagCapturer(data []byte) error
}
