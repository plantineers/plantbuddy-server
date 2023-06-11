// Author: Maximilian Floto
package utils

import (
	"crypto/sha256"
	"fmt"
)

// HashPassword hashes a password with a fixed salt.
func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+"plantbuddy_salt")))
}
