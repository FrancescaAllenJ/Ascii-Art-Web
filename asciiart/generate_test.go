package asciiart

import (
	"strings"
	"testing"
)

func TestGenerateASCII_ValidInput(t *testing.T) {
	ascii, err := GenerateASCII("Hi", "standard")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !strings.Contains(ascii, "H") {
		t.Errorf("Expected output to contain ASCII for 'H', got: %s", ascii)
	}
}

func TestGenerateASCII_InvalidChar(t *testing.T) {
	_, err := GenerateASCII("Hello ☺", "standard")
	if err == nil {
		t.Errorf("Expected error for unsupported character, got nil")
	}
}

func TestConvert_MultiLine(t *testing.T) {
	text := "Hello\nWorld"
	ascii, err := Convert(text, "standard")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !strings.Contains(ascii, "Hello") && !strings.Contains(ascii, "World") {
		t.Errorf("Expected output to include both 'Hello' and 'World', got: %s", ascii)
	}
}

func TestConvert_EmptyInput(t *testing.T) {
	ascii, err := Convert("", "standard")
	if err != nil {
		t.Errorf("Expected no error for empty string, got: %v", err)
	}
	if ascii != "" {
		t.Errorf("Expected empty output for empty input, got: %q", ascii)
	}
}
