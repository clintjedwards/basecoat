package utils

import (
	"math/rand"
	"time"
)

// GenerateRandString generates a variable length string; can be used for ids
func GenerateRandString(length int) []byte {

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return b
}
