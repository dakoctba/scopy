package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dakoctba/scopy/pkg"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"

	// Flags
	headerFormat    string
	excludePatterns string
	listOnly        bool
	maxSize         string
	stripComments   bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "scopy [extensões...]",
	Short: "Smart Copy - Copia conteúdo de arquivos com extensões específicas",
	Long: `Scopy é uma ferramenta de linha de comando que permite copiar o conteúdo
de arquivos com extensões específicas de forma inteligente, respeitando
configurações de exclusão e formatos personalizados.`,
	Example: `  scopy go js                               # Copia arquivos .go e .js
  scopy --header-format "/* %s */" go       # Customiza o formato do cabeçalho
  scopy --exclude "vendor,dist" go js       # Ignora diretórios vendor e dist
  scopy --list-only go                      # Lista apenas arquivos .go sem mostrar conteúdo
  scopy --max-size 500KB go                 # Ignora arquivos .go maiores que 500KB
  scopy --strip-comments go js              # Remove comentários dos arquivos copiados`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Convert maximum size to bytes
		var maxSizeBytes int64
		if maxSize != "" {
			var err error
			maxSizeBytes, err = parseSize(maxSize)
			if err != nil {
				return fmt.Errorf("erro ao analisar tamanho máximo: %v", err)
			}
		}

		// Configure processor
		config := pkg.Config{
			HeaderFormat:    headerFormat,
			ExcludePatterns: strings.Split(excludePatterns, ","),
			ListOnly:        listOnly,
			MaxSize:         maxSizeBytes,
			StripComments:   stripComments,
			Extensions:      args,
		}

		processor := pkg.NewProcessor(config)
		err := processor.Process(".")
		if err != nil {
			return fmt.Errorf("erro ao processar arquivos: %v", err)
		}

		// Display statistics
		stats := processor.GetStats()
		fmt.Printf("\nEstatísticas:\n")
		fmt.Printf("Total de arquivos: %d\n", stats.TotalFiles)
		fmt.Printf("Arquivos por extensão:\n")
		for ext, count := range stats.FilesByExt {
			fmt.Printf("  %s: %d\n", ext, count)
		}
		fmt.Printf("Total de bytes: %d\n", stats.TotalBytes)
		fmt.Printf("Total de linhas: %d\n", stats.TotalLines)

		return nil
	},
}

func parseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(sizeStr)
	var multiplier int64 = 1

	if strings.HasSuffix(sizeStr, "KB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "KB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return size * multiplier, nil
}

func init() {
	rootCmd.Flags().StringVarP(&headerFormat, "header-format", "f", "// file: %s", "Format of the header that precedes each file")
	rootCmd.Flags().StringVarP(&excludePatterns, "exclude", "e", "", "Padrões para excluir arquivos/diretórios (separados por vírgula)")
	rootCmd.Flags().BoolVarP(&listOnly, "list-only", "l", false, "Apenas lista os arquivos que seriam copiados")
	rootCmd.Flags().StringVarP(&maxSize, "max-size", "s", "", "Tamanho máximo dos arquivos a serem incluídos")
	rootCmd.Flags().BoolVarP(&stripComments, "strip-comments", "c", false, "Remove comentários dos arquivos de código")

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Exibe a versão do aplicativo",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("scopy versão %s\n", version)
		},
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
