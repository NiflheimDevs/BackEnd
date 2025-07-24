package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func Contains[T comparable](slice []T, item T) bool {
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

func GenerateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
