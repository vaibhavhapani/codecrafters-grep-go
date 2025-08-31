package main

import (
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

	ok, err := handler(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func handler(line []byte, pattern string) (bool, error){
	input := string(line)

	for i := 0; i <= len(input); i++ {
		if(matchAt(input, pattern, i)) {
			return true, nil
		}
	}

	return false, nil
}

func matchAt(input string, pattern string, pos int) bool {
	if len(pattern) == 0 {
		return true
	}

	if pos >= len(input) {
		return false
	}

	element, elementLen := parsePatternElement(pattern)

	char, charLen := utf8.DecodeRuneInString(input[pos:])
	if char == utf8.RuneError {
		return false
	}

	if(!matchElement(char, element)) {
		return false
	}

	return matchAt(input, pattern[elementLen:], pos+charLen)
}

func parsePatternElement(pattern string) (string, int) {
	if len(pattern) == 0 {
		return "", 0
	}

	if pattern[0] == '\\' && len(pattern) > 1 {
		return pattern[:2], 2
	}

	if pattern[0] == '[' {
		for i := 1; i < len(pattern); i++ {
			if pattern[i] == ']' {
				return pattern[:i+1], i+1
			}
		}
		fmt.Printf("Error: There's no closing bracket, malformed pattern!!!")
		return pattern[:1], 1
	}

	return pattern[:1], 1
}


func matchElement(char rune, element string) bool {
	switch element {
	case `\d`:
		return unicode.IsDigit(char)
	case `\w`:
		return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_'
	default:
		if len(element) >= 3 && element[0] == '[' && element[len(element)-1] == ']' {
			return matchCharGroup(char, element)
		}
		if len(element) == 1 {
			return char == rune(element[0])
		}
	}
	return false
}

func matchCharGroup(char rune, pattern string) bool {
	if len(pattern) < 3 {
		return false
	}

	if pattern[1] == '^' {
		for _, c := range pattern[2:len(pattern)-1] {
			if char == c {
				return false
			}
		}
		return true
	}

	for _, c := range pattern[1:len(pattern)-1] {
		if char == c {
			return true
		}
	}

	return false
}