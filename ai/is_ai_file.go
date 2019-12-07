package ai

import "strings"

// IsAiFile ...
func IsAiFile(fileName string) bool {
	return strings.HasPrefix(fileName, "AI") &&
		(strings.HasSuffix(fileName, ".cc") || strings.HasSuffix(fileName, ".o"))
}
