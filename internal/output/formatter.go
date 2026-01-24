package output

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/dreamiurg/nato/internal/converter"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// ColorMode determines how colors are applied to output.
type ColorMode string

const (
	ColorModeWord   ColorMode = "word"   // Color entire phonetic word
	ColorModeLetter ColorMode = "letter" // Color only first character
)

// Formatter handles output formatting with optional colors.
type Formatter struct {
	noColor bool
	isTTY   bool
	colors  []*color.Color
	mode    ColorMode
}

// NewFormatter creates a new output formatter.
func NewFormatter(noColor bool, mode ColorMode) *Formatter {
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
		mode:    mode,
		colors: []*color.Color{
			color.New(color.FgRed),
			color.New(color.FgGreen),
			color.New(color.FgYellow),
			color.New(color.FgBlue),
			color.New(color.FgMagenta),
			color.New(color.FgCyan),
			color.New(color.FgHiRed),
			color.New(color.FgHiGreen),
		},
	}
}

// letterColorIndex returns a deterministic color index for a given rune.
// Same letter always maps to same color index.
func (f *Formatter) letterColorIndex(r rune) int {
	upper := unicode.ToUpper(r)

	// Map A-Z to 0-25, 0-9 to 26-35
	var index int
	if upper >= 'A' && upper <= 'Z' {
		index = int(upper - 'A')
	} else if upper >= '0' && upper <= '9' {
		index = 26 + int(upper-'0')
	} else {
		index = 0
	}

	return index % len(f.colors)
}

// Print outputs the converted words to stdout.
func (f *Formatter) Print(words []converter.Word) {
	for _, word := range words {
		var phonetics []string
		var prevColorIdx = -1
		var prevRune rune

		for _, result := range word.Results {
			if f.noColor {
				phonetics = append(phonetics, result.Phonetic)
				continue
			}

			// Get base color index for this letter
			colorIdx := f.letterColorIndex(result.Original)

			// Contrast-aware: if adjacent different letters would have same color, shift
			if prevColorIdx >= 0 && colorIdx == prevColorIdx && result.Original != prevRune {
				colorIdx = (colorIdx + 1) % len(f.colors)
			}

			c := f.colors[colorIdx]

			// Apply color based on mode
			var colored string
			if f.mode == ColorModeLetter && len(result.Phonetic) > 0 {
				// Color only first character
				first := string(result.Phonetic[0])
				rest := result.Phonetic[1:]
				colored = c.Sprint(first) + rest
			} else {
				// Color entire word (default)
				colored = c.Sprint(result.Phonetic)
			}
			phonetics = append(phonetics, colored)

			prevColorIdx = colorIdx
			prevRune = result.Original
		}
		fmt.Println(strings.Join(phonetics, " "))
	}
}
