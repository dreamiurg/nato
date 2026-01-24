package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/dreamiurg/nato/internal/converter"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// Formatter handles output formatting with optional colors.
type Formatter struct {
	noColor bool
	isTTY   bool
	colors  []*color.Color
}

// NewFormatter creates a new output formatter.
func NewFormatter(noColor bool) *Formatter {
	isTTY := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		noColor = true
	}

	// Disable colors if not a TTY
	if !isTTY {
		noColor = true
	}

	return &Formatter{
		noColor: noColor,
		isTTY:   isTTY,
		colors: []*color.Color{
			color.New(color.FgCyan),
			color.New(color.FgYellow),
			color.New(color.FgGreen),
			color.New(color.FgMagenta),
		},
	}
}

// Print outputs the converted words to stdout.
func (f *Formatter) Print(words []converter.Word) {
	for _, word := range words {
		var phonetics []string
		for i, result := range word.Results {
			if f.noColor {
				phonetics = append(phonetics, result.Phonetic)
			} else {
				// Alternate colors for visual distinction
				c := f.colors[i%len(f.colors)]
				phonetics = append(phonetics, c.Sprint(result.Phonetic))
			}
		}
		fmt.Println(strings.Join(phonetics, " "))
	}
}
