package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateInvitationCode generates a unique invitation code for classroom
// Format: CLS-XXXXXXXX (e.g., CLS-A7K9M2X5)
func GenerateInvitationCode(prefix string) string {
	// Always use "CLS" prefix for classroom codes
	randomPart := generateRandomString(8)
	return fmt.Sprintf("CLS-%s", randomPart)
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}
