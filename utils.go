package main

import (
	"math/rand"
	"strings"
	"time"
)

// Function to generate a random string of specified length
func generateRandomString(length int) string {
	// Define the characters that can appear in the string (both lowercase and uppercase)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Create a slice to store the random string
	var sb strings.Builder
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Generate a random string of the desired length
	for i := 0; i < length; i++ {
		// Randomly select a character from the charset and append to the string builder
		randomIndex := rng.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
