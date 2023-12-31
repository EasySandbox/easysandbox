package generatepassword

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateRandomPassword(length uint8) (string, error) {

	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	// Define the characters to be used in the password
	charset := "abcdefghjkmnopqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789!@#$%^&*()_+{}[]:;<>,.?/~"

	// Initialize a random password slice
	password := make([]byte, length)

	// Calculate the maximum index of the character set
	maxIndex := big.NewInt(int64(len(charset)))

	// Generate a random password by selecting characters from the charset
	for i := uint8(0); i < length; i++ {
		index, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", fmt.Errorf("failed to generate password: %w", err)
		}
		password[i] = charset[index.Int64()]
	}

	return string(password), nil
}