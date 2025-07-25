package asciiart

import (
	"fmt"
	"os"
	"strings"
)

func GenerateASCII(text, banner string) (string, error) {
	filePath := "banners/" + banner + ".txt"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read banner file (%s): %w", filePath, err)
	}

	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

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
