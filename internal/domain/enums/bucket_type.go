package enums

type BucketType uint

const (
	ProfilePic BucketType = iota + 1
	Resume
)

func (bt BucketType) String() string {
	switch bt {
	case ProfilePic:
		return "profilePic"
	case Resume:
		return "resume"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		ProfilePic,
		Resume,
	}
}
