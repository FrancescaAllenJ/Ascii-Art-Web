package main

import (
	"fmt"

	"learn.01founders.co/git/ftafrial/ascii-art-web/asciiart"
)

func main() {
	input := "Hello"
	banner := "standard"

	output, err := asciiart.GenerateASCII(input, banner)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(output)
}
