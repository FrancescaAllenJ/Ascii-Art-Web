package asciiart_test

import (
	"strings"
	"testing"

	"learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"
)

// --- helpers ---

// normalize newlines to \n and split into lines (no trimming of spaces)
func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

// return the index of the first non-empty line and its content
func firstNonEmptyLine(lines []string) (int, string) {
	for i, ln := range lines {
		if ln != "" {
			return i, ln
		}
	}
	return -1, ""
}

// compare allowing both leading and trailing spaces around the meaningful shape
func equalIgnoringEdgeSpaces(actual, want string) bool {
	a := strings.TrimSpace(actual)
	w := strings.TrimSpace(want)
	if len(a) < len(w) {
		return false
	}
	return strings.HasPrefix(a, w)
}

// --- tests ---

func TestGenerateASCII_ValidInput(t *testing.T) {
	out, err := asciiart.GenerateASCII("Hello", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatalf("got empty output")
	}

	lines := splitLines(out)
	_, top := firstNonEmptyLine(lines)
	if top == "" {
		t.Fatalf("could not find first non-empty line; output:\n%s", out)
	}

	// Known top row for "Hello" (standard banner). Allow edge spaces.
	wantTop := " _    _          _   _"
	if !equalIgnoringEdgeSpaces(top, wantTop) {
		t.Fatalf("top row mismatch (ignoring edge spaces).\nwant: %q\ngot : %q\nfull output:\n%s", wantTop, top, out)
	}
}

func TestGenerateASCII_InvalidChar(t *testing.T) {
	_, err := asciiart.GenerateASCII("Hi🙂", "standard")
	if err == nil {
		t.Fatalf("expected error for unsupported rune (non-ASCII)")
	}
}

func TestConvert_MultiLine(t *testing.T) {
	out, err := asciiart.Convert("Hello\nWorld", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatalf("got empty output")
	}

	// Sanity: expect a blank line somewhere between blocks
	norm := strings.ReplaceAll(out, "\r\n", "\n")
	if !strings.Contains(norm, "\n\n") {
		t.Fatalf("expected a blank separator line between blocks; got:\n%s", out)
	}

	lines := splitLines(out)

	// First block top row
	startIdx, firstTop := firstNonEmptyLine(lines)
	wantHelloTop := " _    _          _   _"
	if !equalIgnoringEdgeSpaces(firstTop, wantHelloTop) {
		t.Fatalf("first block top row mismatch (ignoring edge spaces).\nwant: %q\ngot : %q\nfull output:\n%s", wantHelloTop, firstTop, out)
	}

	// Second block top row: skip exactly 8 rows of the first block,
	// then take the next non-empty line as the second block's top row.
	i := startIdx + 8
	if i < 0 {
		i = 0
	}
	secondTop := ""
	for ; i < len(lines); i++ {
		if lines[i] != "" {
			secondTop = lines[i]
			break
		}
	}
	if secondTop == "" {
		t.Fatalf("could not detect second block top row; output:\n%s", out)
	}

	wantWorldTop := " __          __                 _       _"
	if !equalIgnoringEdgeSpaces(secondTop, wantWorldTop) {
		t.Fatalf("second block top row mismatch (ignoring edge spaces).\nwant: %q\ngot : %q\nfull output:\n%s", wantWorldTop, secondTop, out)
	}
}

func TestConvert_EmptyInput(t *testing.T) {
	out, err := asciiart.Convert("", "standard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("expected empty output for empty input, got: %q", out)
	}
}
