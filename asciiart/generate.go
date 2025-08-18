package asciiart

import (
	"bufio"
	"bytes"
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

// preloadBanner loads and caches the embedded banner lines.
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

	// read embedded file
	data, err := bannerFS.ReadFile("banners/" + banner + ".txt")
	if err != nil {
		return nil, fmt.Errorf("failed to load embedded banner %q: %w", banner, err)
	}

	// ensure trailing LF so the scanner doesn’t drop the last row
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	// scan into logical lines; strip any stray \r (CRLF hardening)
	var lines []string
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasSuffix(line, "\r") {
			line = strings.TrimSuffix(line, "\r")
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	// canonical banners: 856 total lines (95 chars × 9 + final blank)
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
	// Each glyph uses 8 displayed rows; file stores 9 per glyph (8 + 1 separator)
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
func Convert(multilineText, banner string) (string, error) {
	// normalise Windows inputs
	cleaned := strings.ReplaceAll(multilineText, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "")

	parts := strings.Split(cleaned, "\n")
	var out strings.Builder

	for i, line := range parts {
		if line == "" {
			// explicit empty line in input -> emit one blank line
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
