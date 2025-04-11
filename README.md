# Scopy - Smart Copy

Scopy is a command-line tool written in Go that allows you to intelligently copy the content of files with specific extensions, respecting exclusion settings and custom formats.

## Features

- Recursive directory processing
- File extension filtering
- File/directory exclusion by patterns
- File size limit
- Custom header formatting
- Detailed processing statistics
- Support for different comment formats
- Intuitive command-line interface using Cobra

## Installation

### Prerequisites

- Go 1.21 or higher

### Installation via go install

```bash
go install github.com/dakoctba/scopy@latest
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/dakoctba/scopy.git
cd scopy
```

2. Build the project:
```bash
go build
```

3. (Optional) Install the binary:
```bash
go install
```

## Usage

```bash
scopy [options] extension1 extension2 ...
```

### Options

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--header-format` | `-f` | Format of the header preceding each file (default: "// file: %s") | `--header-format "/* %s */"` |
| `--exclude` | `-e` | Patterns to exclude files/directories (comma-separated) | `--exclude "vendor,dist"` |
| `--list-only` | `-l` | Only list files that would be copied (default: true) | `--list-only=false` |
| `--max-size` | `-s` | Maximum size of files to include | `--max-size 500KB` |
| `--strip-comments` | `-c` | Remove comments from code files | `--strip-comments` |

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `version` | Display the application version | `scopy version` |

### Examples

```bash
# List .go and .js files (default behavior)
scopy go js

# Show content of .go and .js files
scopy --list-only=false go js

# Customize header format
scopy -f "/* %s */" go

# Ignore vendor and dist directories
scopy -e "vendor,dist" go js

# Ignore .go files larger than 500KB
scopy -s 500KB go

# Remove comments from copied files
scopy -c go js
```

## Statistics

At the end of execution, Scopy displays detailed statistics about the processed files:

- Total number of files processed
- Number of files per extension
- Total size in bytes
- Total number of lines copied

## Return Codes

| Code | Description |
|------|-------------|
| 0 | Successful execution |
| 1 | Usage error (invalid arguments) |
| 2 | Error reading/processing files |

## Project Structure

```
.
├── cmd/
│   └── root.go      # Main command and configurations
├── pkg/
│   ├── processor.go # File processing logic
│   └── processor_test.go # Unit tests
├── main.go          # Entry point
├── go.mod           # Dependency management
└── README.md        # Documentation
```

## Development

### Development Requirements

- Go 1.21 or higher
- Git

### Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/dakoctba/scopy.git
cd scopy
```

2. Install dependencies:
```bash
go mod tidy
```

### Running Tests

```bash
go test ./pkg/...
```

### Building the Project

```bash
go build
```

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- **Jackson** - [dakoctba](https://github.com/dakoctba)

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - Library for creating CLI applications in Go
- [Go](https://golang.org/) - Programming language
