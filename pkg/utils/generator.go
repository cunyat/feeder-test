package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var numbers = []byte("0123456789")

var rsource = rand.NewSource(time.Now().UnixNano())

// GenerateSKU generates a valid sku for the application
func GenerateSKU() string {
	return fmt.Sprintf("%s-%s\n", pick(letters, 4), pick(numbers, 4))
}

// GenerateSKUs generates many skus
func GenerateSKUs(count int) []string {
	skus := make([]string, count)
	for i := range skus {
		skus[i] = GenerateSKU()
	}

	return skus
}

// pick return a random element from the given source
func pick(source []byte, count int) string {
	b := make([]byte, count)

	for i := range b {
		b[i] = source[randInt(len(source))]
	}

	return string(b)
}

// randInt return a random integer
func randInt(max int) int64 {
	return rsource.Int63() % int64(max)
}
