package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"
)

func main() {
	banner := flag.String("banner", "standard", "standard|shadow|thinkertoy")
	text := flag.String("text", "", "text to render; use \\n for newlines")
	flag.Parse()

	if *text == "" {
		fmt.Fprintln(os.Stderr, "usage: ascii-check -banner standard -text \"Hello\\nThere\"")
		os.Exit(2)
	}
	// allow literal "\n" in the flag to become a real newline
	in := strings.ReplaceAll(*text, `\n`, "\n")

	out, err := asciiart.Convert(in, *banner)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(out)
}
