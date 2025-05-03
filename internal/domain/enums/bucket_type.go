package enums

type BucketType uint

const (
	ProfilePic BucketType = iota + 1
	TeamProfile
	Resume
)

func (bt BucketType) String() string {
	switch bt {
	case ProfilePic:
		return "profilePic"
	case Resume:
		return "resume"
	case TeamProfile:
		return "teamProfile"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		ProfilePic,
		TeamProfile,
		Resume,
	}
}
