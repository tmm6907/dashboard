package utils

import (
	"fmt"
	"os"
)

func ReadDockerSecret(secretName string) (string, error) {
	// Docker mounts secrets inside /run/secrets/ by default
	secretPath := fmt.Sprintf("/run/secrets/%s", secretName)
	content, err := os.ReadFile(secretPath)
	if err != nil {
		return "", err
	}
	// Trim any whitespace (like newlines) from the secret content
	return string(content), nil
}
