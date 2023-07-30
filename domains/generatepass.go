package domains

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func generateRandomPassword(length int) (string, error) {
	// Define the characters to be used in the password
	charset := "abcdefghjkmnopqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789!@#$%^&*()_+{}[]:;<>,.?/~"

	// Initialize a random password slice
	password := make([]byte, length)

	// Calculate the maximum index of the character set
	maxIndex := big.NewInt(int64(len(charset)))

	// Generate a random password by selecting characters from the charset
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", fmt.Errorf("Failed to generate password: %w", err)
		}
		password[i] = charset[index.Int64()]
	}

	return string(password), nil
}
