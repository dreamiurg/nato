package alphabet

// Alphabet represents a phonetic alphabet mapping.
type Alphabet struct {
	Name        string
	DisplayName string
	Description string
	Letters     map[rune]string
	Digits      map[rune]string
}

// Lookup returns the phonetic word for a character, or empty string if not found.
func (a *Alphabet) Lookup(r rune) (string, bool) {
	if word, ok := a.Letters[r]; ok {
		return word, true
	}
	if word, ok := a.Digits[r]; ok {
		return word, true
	}
	return "", false
}
