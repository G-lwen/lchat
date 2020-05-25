package utils

import (
	"math/rand"
	"time"
)

const (
	NumberCode = "0123456789"
	NumberLowerCharasterCode = "0123456789abcdefghizklmnopqrstuvwxyz"
	NumberCharasterCode = "0123456789abcdefghizklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandomCodeGenerator(n int, code string) string {
	letters := []byte(code)
	l := len(letters)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[r.Intn(l)]
	}
	return string(result)
}