package asciiart

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GenerateASCII takes a string and a banner name, and returns the ASCII art string.
func GenerateASCII(input, banner string) (string, error) {
	bannerFile := "banners/" + banner + ".txt"
	file, err := os.Open(bannerFile)
	if err != nil {
		return "", fmt.Errorf("failed to open banner file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	asciiMap := make(map[rune][]string)

	var currentRune rune = 32 // starting from ASCII space
	for scanner.Scan() {
		var lines []string
		for i := 0; i < 8; i++ {
			if !scanner.Scan() {
				return "", fmt.Errorf("unexpected EOF while reading ASCII lines")
			}
			lines = append(lines, scanner.Text())
		}
		asciiMap[currentRune] = lines
		currentRune++
	}

	// Split the input string into lines (in case of newlines)
	lines := strings.Split(input, "\\n")
	var result strings.Builder

	for _, line := range lines {
		for i := 0; i < 8; i++ {
			for _, char := range line {
				if art, ok := asciiMap[char]; ok {
					result.WriteString(art[i])
				} else {
					result.WriteString("        ") // fallback for unknown characters
				}
			}
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}
