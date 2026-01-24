package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dreamiurg/nato/internal/alphabet"
	"github.com/dreamiurg/nato/internal/converter"
	"github.com/dreamiurg/nato/internal/output"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	alphabetName  string
	noColor       bool
	colorMode     string
	listAlphabets bool
	version       = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "nato [text...]",
	Short: "Convert text to phonetic alphabet",
	Long: `nato converts text to NATO phonetic alphabet (or other phonetic alphabets).

Examples:
  nato hello           # Outputs: Hotel Echo Lima Lima Oscar
  nato "Hello World"   # Converts multiple words
  nato -a lapd hello   # Uses LAPD alphabet
  echo "test" | nato   # Accepts piped input`,
	Version: version,
	RunE:    run,
}

func init() {
	rootCmd.Flags().StringVarP(&alphabetName, "alphabet", "a", "nato", "phonetic alphabet to use")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "disable colored output")
	rootCmd.Flags().StringVar(&colorMode, "color-mode", "word", "color mode: 'word' (full word) or 'letter' (first char only)")
	rootCmd.Flags().BoolVarP(&listAlphabets, "list-alphabets", "l", false, "list available alphabets")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if listAlphabets {
		return listAlphabetsCmd()
	}

	// Get input text from args or stdin
	text := strings.Join(args, " ")
	if text == "" {
		// Check if stdin has data (is a pipe)
		if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			text = strings.Join(lines, " ")
		}
	}

	// Show help if still no input
	if text == "" {
		return cmd.Help()
	}

	// Get alphabet
	alpha, err := alphabet.Get(alphabetName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\nAvailable alphabets: %s\n", err, strings.Join(alphabet.Names(), ", "))
		return err
	}

	// Convert text
	words := converter.Convert(text, alpha)

	// Format output
	formatter := output.NewFormatter(noColor, output.ColorMode(colorMode))
	formatter.Print(words)

	return nil
}

func listAlphabetsCmd() error {
	fmt.Println("Available alphabets:")
	for _, a := range alphabet.List() {
		fmt.Printf("  %-14s - %s\n", a.Name, a.Description)
	}
	return nil
}
