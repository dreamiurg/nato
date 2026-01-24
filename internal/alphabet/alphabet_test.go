package alphabet

import (
	"testing"
)

func TestNATOAlphabetLetters(t *testing.T) {
	alpha, err := Get("nato")
	if err != nil {
		t.Fatalf("failed to get NATO alphabet: %v", err)
	}

	tests := []struct {
		input    rune
		expected string
	}{
		{'A', "Alfa"},
		{'B', "Bravo"},
		{'C', "Charlie"},
		{'D', "Delta"},
		{'E', "Echo"},
		{'F', "Foxtrot"},
		{'G', "Golf"},
		{'H', "Hotel"},
		{'I', "India"},
		{'J', "Juliett"},
		{'K', "Kilo"},
		{'L', "Lima"},
		{'M', "Mike"},
		{'N', "November"},
		{'O', "Oscar"},
		{'P', "Papa"},
		{'Q', "Quebec"},
		{'R', "Romeo"},
		{'S', "Sierra"},
		{'T', "Tango"},
		{'U', "Uniform"},
		{'V', "Victor"},
		{'W', "Whiskey"},
		{'X', "Xray"},
		{'Y', "Yankee"},
		{'Z', "Zulu"},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got, ok := alpha.Lookup(tc.input)
			if !ok {
				t.Errorf("Lookup(%c) returned not found", tc.input)
			}
			if got != tc.expected {
				t.Errorf("Lookup(%c) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestNATOAlphabetDigits(t *testing.T) {
	alpha, err := Get("nato")
	if err != nil {
		t.Fatalf("failed to get NATO alphabet: %v", err)
	}

	tests := []struct {
		input    rune
		expected string
	}{
		{'0', "Zero"},
		{'1', "One"},
		{'2', "Two"},
		{'3', "Three"},
		{'4', "Four"},
		{'5', "Five"},
		{'6', "Six"},
		{'7', "Seven"},
		{'8', "Eight"},
		{'9', "Niner"},
	}

	for _, tc := range tests {
		t.Run(string(tc.input), func(t *testing.T) {
			got, ok := alpha.Lookup(tc.input)
			if !ok {
				t.Errorf("Lookup(%c) returned not found", tc.input)
			}
			if got != tc.expected {
				t.Errorf("Lookup(%c) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestGetUnknownAlphabet(t *testing.T) {
	_, err := Get("unknown")
	if err == nil {
		t.Error("expected error for unknown alphabet")
	}
}

func TestGetCaseInsensitive(t *testing.T) {
	tests := []string{"nato", "NATO", "Nato", "nAtO"}
	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			alpha, err := Get(name)
			if err != nil {
				t.Errorf("Get(%q) failed: %v", name, err)
			}
			if alpha.Name != "nato" {
				t.Errorf("Get(%q).Name = %q, want %q", name, alpha.Name, "nato")
			}
		})
	}
}

func TestList(t *testing.T) {
	alphabets := List()
	if len(alphabets) == 0 {
		t.Error("List() returned empty slice")
	}

	// Verify NATO is in the list
	found := false
	for _, a := range alphabets {
		if a.Name == "nato" {
			found = true
			break
		}
	}
	if !found {
		t.Error("NATO alphabet not found in List()")
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Error("Names() returned empty slice")
	}

	// Verify nato is in the list
	found := false
	for _, name := range names {
		if name == "nato" {
			found = true
			break
		}
	}
	if !found {
		t.Error("nato not found in Names()")
	}
}
