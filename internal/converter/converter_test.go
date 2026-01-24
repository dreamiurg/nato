package converter

import (
	"testing"

	"github.com/dreamiurg/nato/internal/alphabet"
)

func TestConvertSingleWord(t *testing.T) {
	alpha, _ := alphabet.Get("nato")

	tests := []struct {
		name     string
		input    string
		expected [][]string // Each inner slice is a word's phonetics
	}{
		{
			name:     "simple word",
			input:    "hello",
			expected: [][]string{{"Hotel", "Echo", "Lima", "Lima", "Oscar"}},
		},
		{
			name:     "uppercase word",
			input:    "HELLO",
			expected: [][]string{{"Hotel", "Echo", "Lima", "Lima", "Oscar"}},
		},
		{
			name:     "mixed case word",
			input:    "HeLLo",
			expected: [][]string{{"Hotel", "Echo", "Lima", "Lima", "Oscar"}},
		},
		{
			name:     "single letter",
			input:    "a",
			expected: [][]string{{"Alfa"}},
		},
		{
			name:     "digits",
			input:    "123",
			expected: [][]string{{"One", "Two", "Three"}},
		},
		{
			name:     "mixed letters and digits",
			input:    "a1b2",
			expected: [][]string{{"Alfa", "One", "Bravo", "Two"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			words := Convert(tc.input, alpha)
			if len(words) != len(tc.expected) {
				t.Fatalf("Convert(%q) returned %d words, want %d", tc.input, len(words), len(tc.expected))
			}

			for i, word := range words {
				if len(word.Results) != len(tc.expected[i]) {
					t.Errorf("word %d has %d results, want %d", i, len(word.Results), len(tc.expected[i]))
					continue
				}
				for j, result := range word.Results {
					if result.Phonetic != tc.expected[i][j] {
						t.Errorf("word %d result %d = %q, want %q", i, j, result.Phonetic, tc.expected[i][j])
					}
				}
			}
		})
	}
}

func TestConvertMultipleWords(t *testing.T) {
	alpha, _ := alphabet.Get("nato")

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:  "two words",
			input: "hi there",
			expected: [][]string{
				{"Hotel", "India"},
				{"Tango", "Hotel", "Echo", "Romeo", "Echo"},
			},
		},
		{
			name:  "multiple spaces",
			input: "a  b",
			expected: [][]string{
				{"Alfa"},
				{"Bravo"},
			},
		},
		{
			name:  "leading space",
			input: " hi",
			expected: [][]string{
				{"Hotel", "India"},
			},
		},
		{
			name:  "trailing space",
			input: "hi ",
			expected: [][]string{
				{"Hotel", "India"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			words := Convert(tc.input, alpha)
			if len(words) != len(tc.expected) {
				t.Fatalf("Convert(%q) returned %d words, want %d", tc.input, len(words), len(tc.expected))
			}

			for i, word := range words {
				if len(word.Results) != len(tc.expected[i]) {
					t.Errorf("word %d has %d results, want %d", i, len(word.Results), len(tc.expected[i]))
					continue
				}
				for j, result := range word.Results {
					if result.Phonetic != tc.expected[i][j] {
						t.Errorf("word %d result %d = %q, want %q", i, j, result.Phonetic, tc.expected[i][j])
					}
				}
			}
		})
	}
}

func TestConvertSpecialCharacters(t *testing.T) {
	alpha, _ := alphabet.Get("nato")

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:     "with punctuation",
			input:    "hi!",
			expected: [][]string{{"Hotel", "India"}},
		},
		{
			name:     "with hyphen",
			input:    "a-b",
			expected: [][]string{{"Alfa", "Bravo"}},
		},
		{
			name:     "only punctuation",
			input:    "!!!",
			expected: [][]string{},
		},
		{
			name:     "empty string",
			input:    "",
			expected: [][]string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			words := Convert(tc.input, alpha)
			if len(words) != len(tc.expected) {
				t.Fatalf("Convert(%q) returned %d words, want %d", tc.input, len(words), len(tc.expected))
			}
		})
	}
}

func TestConvertDima(t *testing.T) {
	alpha, _ := alphabet.Get("nato")
	words := Convert("dima", alpha)

	expected := []string{"Delta", "India", "Mike", "Alfa"}

	if len(words) != 1 {
		t.Fatalf("Convert(\"dima\") returned %d words, want 1", len(words))
	}

	if len(words[0].Results) != len(expected) {
		t.Fatalf("word has %d results, want %d", len(words[0].Results), len(expected))
	}

	for i, result := range words[0].Results {
		if result.Phonetic != expected[i] {
			t.Errorf("result %d = %q, want %q", i, result.Phonetic, expected[i])
		}
	}
}
