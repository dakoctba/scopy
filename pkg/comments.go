// file: /Users/jackson/workspace/meus_projetos/scopy/pkg/comments.go
package pkg

import (
	"strings"
)

// Common line comment markers across different programming languages
var commonLineCommentMarkers = []string{
	"//",
	"#",
	"--",
	";",
	"%",
}

// IsLineComment checks if a line starts with a common comment marker
// Spaces at the beginning of the line are ignored
func IsLineComment(line string) bool {
	// Ensure we're only removing lines that BEGIN with a comment
	// by checking the trimmed line starts with a comment marker
	trimmedLine := strings.TrimSpace(line)

	// If the line is empty, it's not a comment
	if trimmedLine == "" {
		return false
	}

	// Special case for C-style block comments on a single line
	if strings.HasPrefix(trimmedLine, "/*") && strings.HasSuffix(trimmedLine, "*/") {
		return true
	}

	// Check if the line starts with any of the common comment markers
	for _, marker := range commonLineCommentMarkers {
		if strings.HasPrefix(trimmedLine, marker) {
			// This ensures that comments in the middle or at the end of a line are not detected
			return true
		}
	}

	return false
}
