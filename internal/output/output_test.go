package output

import (
	"os"
	"testing"
)

func TestNewFormatterNoColor(t *testing.T) {
	// With noColor=true, colors should be disabled
	f := NewFormatter(true, ColorModeWord)
	if !f.noColor {
		t.Error("NewFormatter(true) should have noColor=true")
	}
}

func TestNewFormatterRespectsNoColorEnv(t *testing.T) {
	// Set NO_COLOR env var
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	f := NewFormatter(false, ColorModeWord)
	if !f.noColor {
		t.Error("NewFormatter should respect NO_COLOR env var")
	}
}

func TestFormatterHasColors(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)
	if len(f.colors) == 0 {
		t.Error("Formatter should have colors defined")
	}
}

func TestFormatterHasMinimumColors(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)
	if len(f.colors) < 6 {
		t.Errorf("Formatter should have at least 6 colors for contrast, got %d", len(f.colors))
	}
}

func TestColorModeDefaults(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)
	if f.mode != ColorModeWord {
		t.Errorf("expected ColorModeWord, got %v", f.mode)
	}

	f2 := NewFormatter(true, ColorModeLetter)
	if f2.mode != ColorModeLetter {
		t.Errorf("expected ColorModeLetter, got %v", f2.mode)
	}
}

func TestLetterColorIndexConsistency(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)

	// Same letter should always return same index
	idx1 := f.letterColorIndex('A')
	idx2 := f.letterColorIndex('A')
	if idx1 != idx2 {
		t.Error("letterColorIndex should be deterministic")
	}

	// Different letters should potentially differ
	idxA := f.letterColorIndex('A')
	idxB := f.letterColorIndex('B')
	if idxA == idxB {
		t.Error("A and B should have different color indices")
	}
}

func TestLetterColorIndexCaseInsensitive(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)
	if f.letterColorIndex('a') != f.letterColorIndex('A') {
		t.Error("letterColorIndex should be case-insensitive")
	}
	if f.letterColorIndex('z') != f.letterColorIndex('Z') {
		t.Error("letterColorIndex should be case-insensitive for Z")
	}
}

func TestLetterColorIndexDigits(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)

	// Digits should have deterministic indices
	idx0a := f.letterColorIndex('0')
	idx0b := f.letterColorIndex('0')
	if idx0a != idx0b {
		t.Error("digit color index should be deterministic")
	}

	// Different digits should differ (at least some)
	idx0 := f.letterColorIndex('0')
	idx1 := f.letterColorIndex('1')
	if idx0 == idx1 {
		t.Error("0 and 1 should have different color indices")
	}
}

func TestLetterColorIndexRange(t *testing.T) {
	f := NewFormatter(true, ColorModeWord)

	// All indices should be within palette range
	for r := 'A'; r <= 'Z'; r++ {
		idx := f.letterColorIndex(r)
		if idx < 0 || idx >= len(f.colors) {
			t.Errorf("color index for %c out of range: %d", r, idx)
		}
	}
	for r := '0'; r <= '9'; r++ {
		idx := f.letterColorIndex(r)
		if idx < 0 || idx >= len(f.colors) {
			t.Errorf("color index for %c out of range: %d", r, idx)
		}
	}
}
