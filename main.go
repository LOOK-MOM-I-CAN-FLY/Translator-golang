package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"translator/parser"
)

type syntaxErrorListener struct {
	*antlr.DefaultErrorListener
	err error
}

func (l *syntaxErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line, column int, msg string, _ antlr.RecognitionException) {
	if l.err == nil {
		l.err = fmt.Errorf("syntax error at %d:%d: %s", line, column, msg)
	}
}

func runCode(code string, interp *Interpreter) error {
	input := antlr.NewInputStream(code)
	lexer := parser.NewSimpleLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewSimpleParser(stream)
	p.BuildParseTrees = true

	errListener := &syntaxErrorListener{}
	p.RemoveErrorListeners()
	p.AddErrorListener(errListener)

	tree := p.Program()
	if errListener.err != nil {
		return errListener.err
	}

	interp.Visit(tree)
	return interp.Err()
}

func runOnce(code string) error {
	interp := NewInterpreter()
	return runCode(code, interp)
}

func normalizeREPLInput(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}
	// add semicolon for single-line statements when missing
	if strings.HasSuffix(trimmed, ";") || strings.HasSuffix(trimmed, "}") {
		return trimmed
	}
	return trimmed + ";"
}

func repl() {
	fmt.Println("Simple Go Interpreter (REPL mode)")
	fmt.Println("Type 'exit' to quit, 'run FILE' to run a file")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	interp := NewInterpreter()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		if strings.HasPrefix(input, "run ") {
			filename := strings.TrimPrefix(input, "run ")
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				continue
			}
			if err := runCode(string(content), interp); err != nil {
				fmt.Printf("Error: %v\n", err)
				interp.err = nil
			}
			continue
		}

		code := normalizeREPLInput(input)
		if code == "" {
			continue
		}
		if err := runCode(code, interp); err != nil {
			fmt.Printf("Error: %v\n", err)
			interp.err = nil
		}
	}
}

func printUsage() {
	fmt.Println("Simple Go Interpreter")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  translator                    - Start REPL")
	fmt.Println("  translator FILE               - Run file")
	fmt.Println("  translator -c CODE            - Run code")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  translator -c \"var x int = 5; fmt.Println(x);\"")
}

func main() {
	if len(os.Args) < 2 {
		repl()
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return
	}

	if os.Args[1] == "-c" && len(os.Args) > 2 {
		if err := runOnce(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if os.Args[1] == "repl" {
		repl()
		return
	}

	if strings.HasPrefix(os.Args[1], "-") {
		printUsage()
		os.Exit(1)
	}

	filename := os.Args[1]
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	if err := runOnce(string(content)); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
