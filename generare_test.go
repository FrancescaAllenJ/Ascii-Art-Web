package asciiart

import (
	"strings"
	"testing"
)

func TestGenerateASCII_ValidInput(t *testing.T) {
	output, err := GenerateASCII("H", "standard")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !strings.Contains(output, "| |  | | (_)") {
		t.Errorf("Expected ASCII art for 'H', got:\n%s", output)
	}
}

func TestConvert_MultiLine(t *testing.T) {
	input := "Hello\nWorld"
	output, err := Convert(input, "standard")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !strings.Contains(output, " __          __") || !strings.Contains(output, "| |  | |        | | | |") {
		t.Errorf("Expected ASCII output for 'Hello\\nWorld', got:\n%s", output)
	}
}

func TestConvert_EmptyInput(t *testing.T) {
	output, err := Convert("", "standard")
	if err != nil {
		t.Errorf("Expected no error for empty input, got: %v", err)
	}
	if output != "\n" {
		t.Errorf("Expected '\\n' for empty input, got: %q", output)
	}
}
