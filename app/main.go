package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) != 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	if pattern == `\d` {
		return isDigit(line)
	}

	var ok bool

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	// Uncomment this to pass the first stage

	ok = bytes.ContainsAny(line, pattern)

	return ok, nil
}

func isDigit(line []byte) (bool, error) {
	for _, b := range line {
		if !unicode.IsDigit(rune(b)) {
			return false, nil
		}
	}
	return true, nil
}