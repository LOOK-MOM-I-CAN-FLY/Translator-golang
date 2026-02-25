package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Simple Go Interpreter")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  translator                    - Start REPL")
		fmt.Println("  translator FILE               - Run file")
		fmt.Println("  translator -c CODE            - Run code")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  translator -c \"var x int = 5; fmt.Println(x);\"")
		os.Exit(1)
	}

	interp := NewInterpreter()

	if os.Args[1] == "-c" && len(os.Args) > 2 {
		code := os.Args[2]
		if err := interp.Run(code); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		filename := os.Args[1]
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		if err := interp.Run(string(content)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
