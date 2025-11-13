package utils

import (
	"regexp"
)

// ExtractVirtualKey extracts the virtual key from the Authorization header
// Returns the virtual key and a boolean indicating if extraction was successful
func ExtractVirtualKey(authHeader string) (string, bool) {
	if authHeader == "" {
		return "", false
	}

	re := regexp.MustCompile(`^Bearer ([\w-]+)$`)
	matches := re.FindStringSubmatch(authHeader)

	// If the regex yields a different result than 2, the Authorization header format is wrong
	if len(matches) != 2 {
		return "", false
	}

	return matches[1], true
}
