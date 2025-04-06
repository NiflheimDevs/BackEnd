package utils

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RemoveUnordered[T any](slice []T, index *int) []T {
	if *index < 0 || *index >= len(slice) {
		return slice
	}
	slice[*index] = slice[len(slice)-1]
	*index = *index - 1
	return slice[:len(slice)-1]
}
