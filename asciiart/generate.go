package asciiart

import (
	"embed"
	"fmt"
	"strings"
	"sync"
)

//go:embed banners/*.txt
var bannerFS embed.FS

var (
	bannerCache = make(map[string][]string)
	cacheLock   sync.RWMutex
)

var allowed = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
}

// preloadBanner loads & caches banner lines (expects 856 total, incl. final blank).
func preloadBanner(banner string) ([]string, error) {
	if !allowed[banner] {
		return nil, fmt.Errorf("unknown banner: %q", banner)
	}

	// cache hit
	cacheLock.RLock()
	if lines, ok := bannerCache[banner]; ok {
		cacheLock.RUnlock()
		return lines, nil
	}
	cacheLock.RUnlock()

	data, err := bannerFS.ReadFile("banners/" + banner + ".txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded banner %q: %w", banner, err)
	}

	// Normalise line endings: CRLF/CR -> LF
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// Ensure a trailing LF so Split will include a final empty element
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}

	// Split into lines; this yields an empty last element for the final blank line
	lines := strings.Split(s, "\n")

	// Trim any accidental trailing \r per line (belt & braces)
	for i := range lines {
		if strings.HasSuffix(lines[i], "\r") {
			lines[i] = strings.TrimSuffix(lines[i], "\r")
		}
	}

	// Expect 95*9 + 1 = 856
	if len(lines) != 856 {
		return nil, fmt.Errorf("banner %q malformed: got %d lines, want 856", banner, len(lines))
	}

	cacheLock.Lock()
	bannerCache[banner] = lines
	cacheLock.Unlock()

	return lines, nil
}

// GenerateASCII renders a single (one-line) string using the specified banner.
func GenerateASCII(text, banner string) (string, error) {
	lines, err := preloadBanner(banner)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	// Each glyph uses 8 visible rows; file stores 9 per glyph (8 + 1 separator)
	for row := 0; row < 8; row++ {
		for _, r := range text {
			if r < 32 || r > 126 {
				return "", fmt.Errorf("unsupported character: %q", r)
			}
			start := int(r-32) * 9
			idx := start + row
			if idx >= len(lines) {
				return "", fmt.Errorf("index out of range for character: %q", r)
			}
			b.WriteString(lines[idx])
		}
		b.WriteString("\n")
	}
	return b.String(), nil
}

// Convert renders multi-line input by joining per-line ASCII blocks,
// inserting a single blank separator line between blocks.
// Empty input returns an empty output (per the tests).
func Convert(multilineText, banner string) (string, error) {
	// Normalise Windows inputs
	cleaned := strings.ReplaceAll(multilineText, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")

	// Special-case: empty input -> empty output
	if cleaned == "" {
		return "", nil
	}

	parts := strings.Split(cleaned, "\n")
	var out strings.Builder

	for i, line := range parts {
		if line == "" {
			// explicit empty line in the input -> emit one blank line
			out.WriteString("\n")
			continue
		}
		ascii, err := GenerateASCII(line, banner)
		if err != nil {
			return "", err
		}
		out.WriteString(ascii)

		// add a single blank separator line between blocks
		if i < len(parts)-1 {
			out.WriteString("\n")
		}
	}
	return out.String(), nil
}
