package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dakoctba/scopy/pkg"
	"github.com/spf13/cobra"
)

var (
	// Versão padrão que será substituída pelo GoReleaser durante a compilação
	// através da flag -X github.com/dakoctba/scopy/cmd.version
	version   = "unknown"
	buildTime = "unknown"
	gitCommit = "unknown"

	// Flags
	headerFormat    string
	excludePatterns string
	maxSize         string
	stripComments   bool
	includeDotFiles bool
	followSymlinks  bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "scopy [extensions...]",
	Short: "Smart Copy - Copy content from files with specific extensions",
	Long: `Scopy is a command line tool that allows copying content
from files with specific extensions intelligently, respecting
exclusion settings and custom formats.`,
	Example: `  scopy go js                               # Copy .go and .js files
  scopy --header-format "/* %s */" go       # Customize header format
  scopy --exclude "vendor,dist" go js       # Ignore vendor and dist directories
  scopy --max-size 500KB go                 # Ignore .go files larger than 500KB
  scopy --strip-comments go js              # Remove comments from copied files
  scopy --all go                            # Include dot files (hidden files)
  scopy --follow go                         # Follow symbolic links`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Convert maximum size to bytes
		var maxSizeBytes int64
		if maxSize != "" {
			var err error
			maxSizeBytes, err = parseSize(maxSize)
			if err != nil {
				return fmt.Errorf("error parsing maximum size: %v", err)
			}
		}

		// Check if output is being redirected
		isRedirected := false
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
			isRedirected = true
		}

		// Configure processor
		config := pkg.Config{
			HeaderFormat:    headerFormat,
			ExcludePatterns: strings.Split(excludePatterns, ","),
			MaxSize:         maxSizeBytes,
			StripComments:   stripComments,
			Extensions:      args,
			OutputToMemory:  !isRedirected, // Store in memory if NOT redirected
			IncludeDotFiles: includeDotFiles,
			FollowSymlinks:  followSymlinks,
		}

		processor := pkg.NewProcessor(config)
		err := processor.Process(".")
		if err != nil {
			return fmt.Errorf("error processing files: %v", err)
		}

		// If not redirected, copy to clipboard
		if !isRedirected {
			output := processor.GetOutput()
			if err := clipboard.WriteAll(output); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not copy to clipboard: %v\n", err)
			} else {
				fmt.Fprintln(os.Stderr, "Content copied to clipboard!")
			}
		}

		// Display statistics to stderr
		stats := processor.GetStats()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Total files: %d\n", stats.TotalFiles)
		fmt.Fprintf(os.Stderr, "Files by extension:\n")
		for ext, count := range stats.FilesByExt {
			fmt.Fprintf(os.Stderr, "  %s: %d\n", ext, count)
		}
		fmt.Fprintf(os.Stderr, "Total bytes: %d\n", stats.TotalBytes)
		fmt.Fprintf(os.Stderr, "Total lines: %d\n", stats.TotalLines)

		// Show comment removal statistics if strip-comments was enabled
		if stripComments && stats.CommentsRemoved > 0 {
			fmt.Fprintf(os.Stderr, "Removed lines (comments): %d\n", stats.CommentsRemoved)
		}

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
	rootCmd.Flags().StringVarP(&excludePatterns, "exclude", "e", "", "Patterns to exclude files/directories (comma-separated)")
	rootCmd.Flags().StringVarP(&maxSize, "max-size", "s", "", "Maximum size of files to be included")
	rootCmd.Flags().BoolVarP(&stripComments, "strip-comments", "c", false, "Remove comments from code files")

	rootCmd.Flags().BoolVarP(&includeDotFiles, "all", "a", false, "Include files & directories beginning with a dot (.)")
	rootCmd.Flags().BoolVarP(&followSymlinks, "follow", "F", false, "Follow symbolic links")

	rootCmd.Flags().BoolP("version", "v", false, "Show version number")

	rootCmd.SetVersionTemplate(`{{.Name}} version {{.Version}}
build time: ` + buildTime + `
git commit: ` + gitCommit + `
`)
	rootCmd.Version = version

	// Permitir que o comando completion apareça
	rootCmd.CompletionOptions.DisableDefaultCmd = false

	// Configurar um template de ajuda personalizado que não liste o comando help
	const customHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

	// Sobrescrever template para omitir o comando help, mas mostrar completion
	const customUsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (and (not .Hidden) (ne .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

	// Aplicar os templates personalizados
	rootCmd.SetHelpTemplate(customHelpTemplate)
	rootCmd.SetUsageTemplate(customUsageTemplate)

	// Remover o comando help padrão
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "hidden-help",
		Hidden: true,
	})

	// Adicionar comando help oculto (para funcionar se alguém digitar 'help')
	helpCmd := &cobra.Command{
		Use:    "help",
		Short:  "Help about any command",
		Hidden: true,
		Run: func(c *cobra.Command, args []string) {
			rootCmd.Help()
		},
	}
	rootCmd.AddCommand(helpCmd)

	// Adicionar o comando version, mas oculto
	versionCmd := &cobra.Command{
		Use:    "version",
		Short:  "Display application version",
		Hidden: true, // Não aparece na ajuda
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("scopy version %s\n", version)
			fmt.Printf("build time: %s\n", buildTime)
			fmt.Printf("git commit: %s\n", gitCommit)
		},
	}
	rootCmd.AddCommand(versionCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
