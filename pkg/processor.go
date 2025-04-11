package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// file: /Users/jackson/workspace/meus_projetos/scopy/pkg/processor.go
// Config contains the settings for file processing
type Config struct {
	HeaderFormat    string
	ExcludePatterns []string
	ListOnly        bool
	MaxSize         int64
	StripComments   bool
	Extensions      []string
}

// Processor is responsible for processing files
type Processor struct {
	config Config
	stats  Stats
}

// Stats contains the processing statistics
type Stats struct {
	TotalFiles int
	FilesByExt map[string]int
	TotalBytes int64
	TotalLines int
}

// NewProcessor creates a new Processor instance
func NewProcessor(config Config) *Processor {
	return &Processor{
		config: config,
		stats: Stats{
			FilesByExt: make(map[string]int),
		},
	}
}

// Process starts the file processing
func (p *Processor) Process(baseDir string) error {
	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore directories
		if info.IsDir() {
			return nil
		}

		// Check if file should be excluded
		if p.shouldExclude(path) {
			return nil
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(path))
		if !p.hasValidExtension(ext) {
			return nil
		}

		// Check maximum size
		if p.config.MaxSize > 0 && info.Size() > p.config.MaxSize {
			return nil
		}

		// Update statistics
		p.stats.TotalFiles++
		p.stats.FilesByExt[ext]++
		p.stats.TotalBytes += info.Size()

		// If only listing, print path and return
		if p.config.ListOnly {
			fmt.Println(path)
			return nil
		}

		// Process the file
		return p.processFile(path)
	})
}

// GetStats returns the processing statistics
func (p *Processor) GetStats() Stats {
	return p.stats
}

func (p *Processor) shouldExclude(path string) bool {
	for _, pattern := range p.config.ExcludePatterns {
		if pattern != "" && strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func (p *Processor) hasValidExtension(ext string) bool {
	if ext == "" {
		return false
	}

	// Remove dot from extension if present
	ext = strings.TrimPrefix(ext, ".")

	for _, validExt := range p.config.Extensions {
		// Remove dot from valid extension if present
		validExt = strings.TrimPrefix(validExt, ".")
		if strings.ToLower(validExt) == ext {
			return true
		}
	}
	return false
}

func (p *Processor) processFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Print header
	fmt.Printf(p.config.HeaderFormat+"\n", path)

	// Copy file content
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return err
	}

	fmt.Println() // Add blank line between files
	return nil
}
