package utils

import (
	"strings"
)

// CheckFileType checks if the provided file type is good
func CheckFileType(typeStr string) bool {
	lower := strings.ToLower(typeStr)

	switch lower {
	case "file":
		fallthrough
	case "directory":
		return true
	default:
		return false
	}
}
