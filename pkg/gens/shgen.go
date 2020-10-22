package gens

import (
	"crypto/rand"
	"fmt"
)

func GenerateShortUrl() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)

	return fmt.Sprintf("%x", bytes)
}

func GenerateId() string {
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)

	return fmt.Sprintf("%x", bytes)
}
