package helpers

import (
	"math/rand"
	"time"
)

func GenerateID() string {
	rand.Seed(time.Now().UnixNano())
	// Generate a 4-word string
	const wordLength = 4
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, wordLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
