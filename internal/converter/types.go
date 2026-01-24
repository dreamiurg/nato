package converter

// ConversionResult represents the result of converting a single character.
type ConversionResult struct {
	Original  rune
	Phonetic  string
	IsSpace   bool
	IsUnknown bool
}

// Word represents a group of conversion results (one input word).
type Word struct {
	Results []ConversionResult
}
