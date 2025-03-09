package util

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func randomChar() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	return string(chars[index.Int64()])
}

func ensureMinLength(s string, minLength int) string {
	for len(s) < minLength {
		s += randomChar()
	}
	return s
}

func shuffle(s string) string {
	runes := []rune(s)
	rng := rand.Reader
	for i := len(runes) - 1; i > 0; i-- {
		jBig, _ := rand.Int(rng, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func HashStringData(input string) string {
	hash := sha256.Sum256([]byte(input))
	hashStr := fmt.Sprintf("%x", hash)
	shortHash := ensureMinLength(hashStr[:7], 7)
	newHash := shortHash + randomChar()

	return shuffle(newHash)
}
