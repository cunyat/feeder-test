package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var numbers = []byte("0123456789")

var rsource = rand.NewSource(time.Now().UnixNano())

func GenerateSKUs(count int) []string {
	skus := make([]string, count)
	for i := range skus {
		skus[i] = fmt.Sprintf("%s-%s\n", pick(letters, 4), pick(numbers, 4))
	}

	return skus
}

func pick(source []byte, count int) string {
	b := make([]byte, count)

	for i := range b {
		b[i] = source[randInt(len(source))]
	}

	return string(b)
}

func randInt(max int) int64 {
	return rsource.Int63() % int64(max)
}
