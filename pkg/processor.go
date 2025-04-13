package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// file: /Users/jackson/workspace/meus_projetos/scopy/pkg/processor.go
// Config contains the settings for file processing
type Config struct {
	HeaderFormat    string
	ExcludePatterns []string
	MaxSize         int64
	StripComments   bool
	Extensions      []string
	OutputToMemory  bool
	IncludeDotFiles bool // Incluir arquivos que começam com ponto (.)
	FollowSymlinks  bool // Seguir links simbólicos
}

// Processor is responsible for processing files
type Processor struct {
	config    Config
	stats     Stats
	gitIgnore *GitIgnore
	output    strings.Builder
}

// Stats contains the processing statistics
type Stats struct {
	TotalFiles      int
	FilesByExt      map[string]int
	TotalBytes      int64
	TotalLines      int
	CommentsRemoved int
}

// NewProcessor creates a new Processor instance
func NewProcessor(config Config) *Processor {
	return &Processor{
		config:    config,
		stats:     Stats{FilesByExt: make(map[string]int)},
		gitIgnore: NewGitIgnore(),
		output:    strings.Builder{},
	}
}

// Process starts the file processing
func (p *Processor) Process(baseDir string) error {
	// Try to load .gitignore
	gitIgnorePath := filepath.Join(baseDir, ".gitignore")
	if _, err := os.Stat(gitIgnorePath); err == nil {
		if err := p.gitIgnore.Load(gitIgnorePath); err != nil {
			return fmt.Errorf("error loading .gitignore: %v", err)
		}
	}

	// Reset total lines count before processing
	p.stats.TotalLines = 0

	// First, count total files to track the last file
	totalFiles := 0
	fileCounter := 0

	// Configuração para a caminhada no sistema de arquivos
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Se não seguimos links simbólicos e este for um erro de link simbólico, ignore
			if !p.config.FollowSymlinks && os.IsNotExist(err) {
				return nil
			}
			return err
		}

		// Se for um diretório, verifique se deve ser ignorado
		if info.IsDir() {
			// Ignora diretórios que começam com . a menos que includeDotFiles esteja ativado
			baseName := filepath.Base(path)
			if !p.config.IncludeDotFiles && strings.HasPrefix(baseName, ".") && path != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Ignora arquivos que começam com . a menos que includeDotFiles esteja ativado
		baseName := filepath.Base(path)
		if !p.config.IncludeDotFiles && strings.HasPrefix(baseName, ".") {
			return nil
		}

		// Skip files that should be excluded
		if p.gitIgnore.ShouldIgnore(path) || p.shouldExclude(path) {
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

		totalFiles++
		return nil
	}

	// Primeira passagem: contagem de arquivos elegíveis
	if p.config.FollowSymlinks {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			// Implementação personalizada para seguir links simbólicos
			if err != nil {
				return err
			}

			// Se for um link simbólico
			if info.Mode()&os.ModeSymlink != 0 {
				// Resolve o link simbólico
				realPath, err := os.Readlink(path)
				if err != nil {
					return nil // Ignora erro ao ler o link
				}

				// Se for um caminho relativo, torna-o absoluto
				if !filepath.IsAbs(realPath) {
					realPath = filepath.Join(filepath.Dir(path), realPath)
				}

				// Obtém informações sobre o destino do link
				destInfo, err := os.Stat(realPath)
				if err != nil {
					return nil // Ignora erro ao acessar o destino do link
				}

				// Se o destino for um diretório, processa-o recursivamente
				if destInfo.IsDir() {
					return filepath.Walk(realPath, walkFunc)
				}

				// Se for um arquivo, processa-o normalmente usando o destino do link
				return walkFunc(realPath, destInfo, nil)
			}

			// Para arquivos e diretórios normais, usa a função de caminhada padrão
			return walkFunc(path, info, err)
		})
		if err != nil {
			return err
		}
	} else {
		err := filepath.Walk(baseDir, walkFunc)
		if err != nil {
			return err
		}
	}

	// Segunda passagem: processamento de arquivos
	processFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Se não seguimos links simbólicos e este for um erro de link simbólico, ignore
			if !p.config.FollowSymlinks && os.IsNotExist(err) {
				return nil
			}
			return err
		}

		// Ignora diretórios
		if info.IsDir() {
			// Ignora diretórios que começam com . a menos que includeDotFiles esteja ativado
			baseName := filepath.Base(path)
			if !p.config.IncludeDotFiles && strings.HasPrefix(baseName, ".") && path != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Ignora arquivos que começam com . a menos que includeDotFiles esteja ativado
		baseName := filepath.Base(path)
		if !p.config.IncludeDotFiles && strings.HasPrefix(baseName, ".") {
			return nil
		}

		// Check if file should be excluded by .gitignore
		if p.gitIgnore.ShouldIgnore(path) {
			return nil
		}

		// Check if file should be excluded by patterns
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

		// Process file
		fileCounter++
		isLastFile := fileCounter == totalFiles
		if err := p.processFile(path, isLastFile); err != nil {
			return err
		}

		return nil
	}

	// Segunda passagem: processa os arquivos
	if p.config.FollowSymlinks {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			// Implementação personalizada para seguir links simbólicos
			if err != nil {
				return err
			}

			// Se for um link simbólico
			if info.Mode()&os.ModeSymlink != 0 {
				// Resolve o link simbólico
				realPath, err := os.Readlink(path)
				if err != nil {
					return nil // Ignora erro ao ler o link
				}

				// Se for um caminho relativo, torna-o absoluto
				if !filepath.IsAbs(realPath) {
					realPath = filepath.Join(filepath.Dir(path), realPath)
				}

				// Obtém informações sobre o destino do link
				destInfo, err := os.Stat(realPath)
				if err != nil {
					return nil // Ignora erro ao acessar o destino do link
				}

				// Se o destino for um diretório, processa-o recursivamente
				if destInfo.IsDir() {
					return filepath.Walk(realPath, processFunc)
				}

				// Se for um arquivo, processa-o normalmente usando o destino do link
				return processFunc(realPath, destInfo, nil)
			}

			// Para arquivos e diretórios normais, usa a função de processamento padrão
			return processFunc(path, info, err)
		})
		return err
	} else {
		return filepath.Walk(baseDir, processFunc)
	}
}

// GetStats returns the processing statistics
func (p *Processor) GetStats() Stats {
	return p.stats
}

// GetOutput returns the output stored in memory
func (p *Processor) GetOutput() string {
	return p.output.String()
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

func (p *Processor) processFile(path string, isLastFile bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Format header
	header := fmt.Sprintf(p.config.HeaderFormat+"\n", path)

	// Count the header line
	p.stats.TotalLines++

	// Write header
	if p.config.OutputToMemory {
		p.output.WriteString(header)
	} else {
		fmt.Print(header)
	}

	// Process content line by line to reduce memory usage
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// If strip comments is enabled, skip lines that are comments
		if p.config.StripComments && IsLineComment(line) {
			// Count removed comment lines
			p.stats.CommentsRemoved++
			continue
		}

		// Write the line to output
		if p.config.OutputToMemory {
			p.output.WriteString(line + "\n")
		} else {
			fmt.Println(line)
		}
		p.stats.TotalLines++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Add blank line between files, but not after the last file
	if !isLastFile {
		if p.config.OutputToMemory {
			p.output.WriteString("\n")
		} else {
			fmt.Println()
		}
		p.stats.TotalLines++
	}

	return nil
}
