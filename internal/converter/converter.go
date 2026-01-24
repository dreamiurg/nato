package converter

import (
	"unicode"

	"github.com/dreamiurg/nato/internal/alphabet"
)

// Convert transforms input text to phonetic words using the given alphabet.
// Returns a slice of Words, where each Word represents one input word.
func Convert(text string, alpha *alphabet.Alphabet) []Word {
	var words []Word
	var currentWord Word

	for _, r := range text {
		// Normalize to uppercase for lookup
		upper := unicode.ToUpper(r)

		if unicode.IsSpace(r) {
			// Word boundary - save current word if not empty
			if len(currentWord.Results) > 0 {
				words = append(words, currentWord)
				currentWord = Word{}
			}
			continue
		}

		// Try to find phonetic word
		if phonetic, ok := alpha.Lookup(upper); ok {
			currentWord.Results = append(currentWord.Results, ConversionResult{
				Original:  r,
				Phonetic:  phonetic,
				IsSpace:   false,
				IsUnknown: false,
			})
		}
		// Non-letter/digit characters are silently skipped per spec
	}

	// Don't forget the last word
	if len(currentWord.Results) > 0 {
		words = append(words, currentWord)
	}

	return words
}
