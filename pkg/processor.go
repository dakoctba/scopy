package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Config contém as configurações para o processamento de arquivos
type Config struct {
	HeaderFormat    string
	ExcludePatterns []string
	ListOnly        bool
	MaxSize         int64
	StripComments   bool
	Extensions      []string
}

// Processor é responsável por processar os arquivos
type Processor struct {
	config Config
	stats  Stats
}

// Stats contém as estatísticas do processamento
type Stats struct {
	TotalFiles int
	FilesByExt map[string]int
	TotalBytes int64
	TotalLines int
}

// NewProcessor cria uma nova instância do Processor
func NewProcessor(config Config) *Processor {
	return &Processor{
		config: config,
		stats: Stats{
			FilesByExt: make(map[string]int),
		},
	}
}

// Process inicia o processamento dos arquivos
func (p *Processor) Process(baseDir string) error {
	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignora diretórios
		if info.IsDir() {
			return nil
		}

		// Verifica se o arquivo deve ser excluído
		if p.shouldExclude(path) {
			return nil
		}

		// Verifica a extensão do arquivo
		ext := strings.ToLower(filepath.Ext(path))
		if !p.hasValidExtension(ext) {
			return nil
		}

		// Verifica o tamanho máximo
		if p.config.MaxSize > 0 && info.Size() > p.config.MaxSize {
			return nil
		}

		// Atualiza estatísticas
		p.stats.TotalFiles++
		p.stats.FilesByExt[ext]++
		p.stats.TotalBytes += info.Size()

		// Se for apenas listar, imprime o caminho e retorna
		if p.config.ListOnly {
			fmt.Println(path)
			return nil
		}

		// Processa o arquivo
		return p.processFile(path)
	})
}

// GetStats retorna as estatísticas do processamento
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

	// Remove o ponto da extensão se presente
	ext = strings.TrimPrefix(ext, ".")

	for _, validExt := range p.config.Extensions {
		// Remove o ponto da extensão válida se presente
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

	// Imprime o cabeçalho
	fmt.Printf(p.config.HeaderFormat+"\n", path)

	// Copia o conteúdo do arquivo
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return err
	}

	fmt.Println() // Adiciona uma linha em branco entre arquivos
	return nil
}
