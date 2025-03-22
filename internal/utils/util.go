package utils

type Utils struct{}

func NewUtils() *Utils {
	return &Utils{}
}

func (u *Utils) Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
