package pkg

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppConfig represents the program configuration
type AppConfig struct {
	InputFile    string
	OutputFile   string
	HeaderFormat string
	ShowHelp     bool
	ShowVersion  bool
	Version      string
	BuildTime    string
	GitCommit    string
}

// ParseFlags parses command line arguments
func (c *AppConfig) ParseFlags() {
	flag.StringVar(&c.InputFile, "i", "", "Input file path")
	flag.StringVar(&c.OutputFile, "o", "", "Output file path")
	flag.StringVar(&c.HeaderFormat, "f", "file: %s", "Header format (default: 'file: %s')")
	flag.BoolVar(&c.ShowHelp, "h", false, "Show help")
	flag.BoolVar(&c.ShowVersion, "v", false, "Show version")
	flag.Parse()
}

// Validate validates the configuration
func (c *AppConfig) Validate() error {
	if c.ShowHelp || c.ShowVersion {
		return nil
	}

	if c.InputFile == "" {
		return fmt.Errorf("input file is required")
	}

	if c.OutputFile == "" {
		// If not specified, use the same name as input file with .txt extension
		ext := filepath.Ext(c.InputFile)
		base := strings.TrimSuffix(c.InputFile, ext)
		c.OutputFile = base + ".txt"
	}

	return nil
}

// PrintUsage displays the help message
func (c *AppConfig) PrintUsage() {
	fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
	fmt.Println("Options:")
	fmt.Println("  -i <file>    Input file path")
	fmt.Println("  -o <file>    Output file path (optional)")
	fmt.Printf("  -f <format>  Header format (default: 'file: %%s')\n")
	fmt.Println("  -h           Show this help message")
	fmt.Println("  -v           Show version information")
	fmt.Println("\nExample:")
	fmt.Printf("  %s -i input.txt -o output.txt -f 'file: %%s'\n", os.Args[0])
}
