package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmdHelp(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--help"})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	_ = cmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "nato") {
		t.Errorf("help output should contain 'nato', got: %s", output)
	}
}

func TestListAlphabets(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--list-alphabets"})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
