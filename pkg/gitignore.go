// file: /Users/jackson/workspace/meus_projetos/scopy/pkg/gitignore.go
package pkg

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GitIgnore represents a .gitignore file parser
type GitIgnore struct {
	patterns []string
}

// NewGitIgnore creates a new GitIgnore instance
func NewGitIgnore() *GitIgnore {
	return &GitIgnore{
		patterns: make([]string, 0),
	}
}

// Load loads patterns from a .gitignore file
func (g *GitIgnore) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		g.patterns = append(g.patterns, line)
	}

	return scanner.Err()
}

// ShouldIgnore checks if a path should be ignored based on .gitignore patterns
func (g *GitIgnore) ShouldIgnore(path string) bool {
	for _, pattern := range g.patterns {
		// Convert pattern to absolute path if it's relative
		absPattern := pattern
		if !filepath.IsAbs(pattern) {
			absPattern = filepath.Join(filepath.Dir(path), pattern)
		}

		// Check if the path matches the pattern
		matched, err := filepath.Match(absPattern, path)
		if err == nil && matched {
			return true
		}

		// Check if the path contains the pattern as a directory
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}
