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
	// after removing leading whitespace
	for _, marker := range commonLineCommentMarkers {
		if strings.HasPrefix(trimmedLine, marker) {
			return true
		}
	}

	return false
}
