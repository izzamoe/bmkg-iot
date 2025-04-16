package utils

import (
	"fmt"
	"time"
)

// generateRandomID creates a simple unique ID for the client
func GenerateRandomID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
