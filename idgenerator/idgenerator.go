package idgenerator

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"
)

type idgenerator struct{}

func New() *idgenerator {
	return &idgenerator{}
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func base62Encode(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	encoded := ""
	for num > 0 {
		remainder := num % 62
		encoded = string(base62Chars[remainder]) + encoded
		num = num / 62
	}

	return encoded
}

func (i *idgenerator) Generate(ctx context.Context) (id string) {

	// Get current time in nanoseconds
	nanoTime := time.Now().UnixNano()

	// Generate a random number between 0 and 100000
	randomNum, err := rand.Int(rand.Reader, big.NewInt(100000))
	if err != nil {
		return ""
	}

	// Encode the time and random number to base62
	timeEncoded := base62Encode(nanoTime)
	randomEncoded := base62Encode(randomNum.Int64())

	// Concatenate the encoded time and random number
	id = randomEncoded + timeEncoded

	return id
}

func (i *idgenerator) GenerateUUID(ctx context.Context) (id string) {
	id = i.Generate(ctx) + i.Generate(ctx)

	return id
}
