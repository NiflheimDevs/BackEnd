package enums

type BucketType uint

const (
	ProfilePic BucketType = iota + 1
)

func (bt BucketType) String() string {
	switch bt {
	case ProfilePic:
		return "profilePic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		ProfilePic,
	}
}
