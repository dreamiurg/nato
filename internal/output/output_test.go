package output

import (
	"os"
	"testing"
)

func TestNewFormatterNoColor(t *testing.T) {
	// With noColor=true, colors should be disabled
	f := NewFormatter(true)
	if !f.noColor {
		t.Error("NewFormatter(true) should have noColor=true")
	}
}

func TestNewFormatterRespectsNoColorEnv(t *testing.T) {
	// Set NO_COLOR env var
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")

	f := NewFormatter(false)
	if !f.noColor {
		t.Error("NewFormatter should respect NO_COLOR env var")
	}
}

func TestFormatterHasColors(t *testing.T) {
	f := NewFormatter(true)
	if len(f.colors) == 0 {
		t.Error("Formatter should have colors defined")
	}
}
