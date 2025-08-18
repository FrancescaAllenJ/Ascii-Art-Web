package asciiart

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	bannerCache = make(map[string][]string)
	cacheLock   sync.RWMutex
)

// preloadBanner loads and caches banner lines to avoid reading from disk repeatedly.
func preloadBanner(banner string) ([]string, error) {
	cacheLock.RLock()
	if lines, ok := bannerCache[banner]; ok {
		cacheLock.RUnlock()
		return lines, nil
	}
	cacheLock.RUnlock()

	filePath := "banners/" + banner + ".txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read banner file (%s): %w", filePath, err)
	}

	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

	cacheLock.Lock()
	bannerCache[banner] = lines
	cacheLock.Unlock()

	return lines, nil
}

// GenerateASCII renders a single line of text using the specified banner file.
func GenerateASCII(text, banner string) (string, error) {
	lines, err := preloadBanner(banner)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	for i := 0; i < 8; i++ {
		for _, r := range text {
			if r < 32 || r > 126 {
				return "", fmt.Errorf("unsupported character: %q", r)
			}
			start := int(r-32) * 9
			if start+i >= len(lines) {
				return "", fmt.Errorf("index out of range for character: %q", r)
			}
			builder.WriteString(lines[start+i])
		}
		builder.WriteString("\n")
	}
	return builder.String(), nil
}

// Convert handles multiline input and renders ASCII for each line separately.
func Convert(multilineText, banner string) (string, error) {
	cleaned := strings.ReplaceAll(multilineText, "\r", "")
	lines := strings.Split(cleaned, "\n")

	var final strings.Builder
	for _, line := range lines {
		if line == "" {
			final.WriteString("\n")
			continue
		}
		ascii, err := GenerateASCII(line, banner)
		if err != nil {
			return "", err
		}
		final.WriteString(ascii)
	}
	return final.String(), nil
}
