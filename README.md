# Scopy - Smart Copy

Scopy is a command-line tool written in Go that allows you to intelligently copy the content of files with specific extensions, respecting exclusion settings and custom formats.

## Features

- Recursive directory processing
- File extension filtering
- File/directory exclusion by patterns
- Automatic .gitignore support
- File size limit
- Custom header formatting
- Automatic clipboard copy
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
make build
```

Ou usando Go diretamente:
```bash
go build
```

3. (Optional) Install the binary:
```bash
make install
```

Ou usando Go diretamente:
```bash
go install
```

## Development

O projeto inclui um Makefile para facilitar o processo de desenvolvimento:

```bash
# Compilar o projeto
make build

# Executar testes
make test

# Executar o aplicativo com argumentos
make run ARGS="go js --all"

# Limpar arquivos de build
make clean

# Instalar localmente
make install

# Desinstalar
make uninstall

# Exibir ajuda do Makefile
make help
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
| `--max-size` | `-s` | Maximum size of files to include | `--max-size 500KB` |
| `--strip-comments` | `-c` | Remove lines that start with comments from code files (default: false) | `--strip-comments` |
| `--all` | `-a` | Include files & directories beginning with a dot (.) | `--all` |
| `--follow` | `-F` | Follow symbolic links | `--follow` |
| `--version` | `-v` | Show version number | `--version` |

### Commands

| Command | Description | Example |
|---------|-------------|---------|
| `version` | Display detailed application version | `scopy version` |

### Examples

```bash
# Copy content of .go and .js files (default behavior)
scopy go js

# Customize header format
scopy -f "/* %s */" go

# Ignore vendor and dist directories
scopy -e "vendor,dist" go js

# Ignore .go files larger than 500KB
scopy -s 500KB go

# Remove comment lines from copied files
scopy -c go js

# Include hidden files (dot files)
scopy -a go js

# Follow symbolic links
scopy -F go js

# Show version information
scopy -v
# or
scopy version
```

## Output Behavior

Scopy has different output behaviors depending on how it's used:

1. **With redirection** (`scopy go js > output.txt`):
   - Content is written to stdout (redirected to file)
   - Statistics are shown in the terminal (stderr)

2. **Without redirection** (`scopy go js`):
   - Content is copied to clipboard
   - Only statistics are shown in the terminal
   - No content is displayed in the terminal

## Clipboard Support

When running Scopy without output redirection, the content of the files is automatically copied to your system's clipboard. This makes it easy to paste the content into any application.

The terminal will only display:
1. A confirmation message when content is copied to clipboard
2. Statistics about the processed files

## Gitignore Support

Scopy automatically reads and respects the `.gitignore` file in your project directory. This means that files and directories listed in your `.gitignore` will be automatically excluded from processing, including:

- `node_modules/`
- `.env` files
- Build artifacts
- Log files
- And any other patterns you've specified in your `.gitignore`

You can still use the `--exclude` flag to add additional patterns that should be ignored.

## Statistics

At the end of execution, Scopy displays statistics about the processed files directly to the terminal:

- Total number of files processed
- Number of files per extension
- Total size in bytes
- Total number of lines copied
- Number of comment lines removed (when `--strip-comments` is enabled)

The statistics are always displayed to stderr, ensuring they don't interfere with output redirection.

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

## Comment Stripping

When using the `--strip-comments` flag, Scopy will remove any line that starts with a common comment marker. By default, this feature is disabled.

The following comment markers are recognized:

- `//` (C, C++, Java, JavaScript, Go, etc.)
- `#` (Python, Ruby, Shell scripts, YAML, etc.)
- `--` (SQL, Lua, etc.)
- `;` (Assembly, INI files, etc.)
- `%` (LaTeX, Matlab, etc.)
- `/* ... */` (C-style block comments when contained on a single line)

Only complete comment lines are removed. Comments that appear in the middle or at the end of a line will be preserved:

```go
// This entire line will be removed
/* This entire line will also be removed */

func main() { // This comment will be kept because it's not at the start of the line
    fmt.Println("Hello") // This comment will also be kept
    doSomething() /* This inline comment is kept */
}
```

The comment stripping is independent of file extension - the same rules apply to all files.

## Environment Setup

Scopy uses environment variables for configuration. These can be set using a `.env` file:

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit the `.env` file and add your GitHub token:
```
GITHUB_TOKEN=your_github_token_here
```

The application will automatically load these variables when it starts.

## Releasing

Scopy uses [GoReleaser](https://goreleaser.com) to create releases for multiple platforms.

### Requirements

- [GoReleaser](https://goreleaser.com/install/)
- Git tag with version number
- GitHub token (configured in `.env` file)

### Creating a New Release

1. Create a new Git tag:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

2. Run the release using Make:
```bash
make release ARGS="--clean"
```

Ou, alternativamente:
```bash
bin/release.sh --clean
```

Este processo irá:
- Carregar o token do GitHub do arquivo `.env`
- Executar o GoReleaser para criar uma release no GitHub
- Gerar binários para macOS (Intel e Apple Silicon), Windows e Linux (amd64 e arm64)

### Testing a Release

Para testar o processo de release sem publicar:

```bash
make snapshot
```

Ou usando o script diretamente:
```bash
bin/release.sh --snapshot --clean
```

### Manual Release (Alternative)

Se você preferir executar o GoReleaser diretamente:

```bash
# Export your GitHub token
export GITHUB_TOKEN=your_github_token_here

# Run GoReleaser
goreleaser release
```
