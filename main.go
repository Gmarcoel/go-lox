package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	DEFAULT_COLOR = "\033[0m"
	RED_COLOR     = "\033[0;31m"
	YELLOW_COLOR  = "\033[0;33m"
)

type TokenType int

const (
	// NOT FOUND
	NOT_FOUND TokenType = iota

	// Single-character tokens.
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}



var errorFlag = false

func runFile(filepath string) {
	f, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("error reading from file")
		return
	}

	// fmt.Println(f)
	run(string(f))

}

func runPrompt() {
	var reader = bufio.NewReader(os.Stdin)
	for {
		var line string
		fmt.Printf(YELLOW_COLOR + "> " + DEFAULT_COLOR)
		line, _ = reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		}
		run(line)

		errorFlag = false
	}
}

func run(source string) {
	var interpreter = newInterpreter()
	var scanner = newScanner(source)
	var tokens = scanner.scanTokens()
	var parser = newParser(tokens)
	var statements = parser.parse()
	var resolver = newResolver(interpreter)
	resolver.resolve(statements)
	if errorFlag {
		return
	}
	if statements != nil {
		interpreter.interpret(statements)
	}
}

func report(line int, where string, message string) {
	fmt.Println("error report in line ", line, "in", where, "with message ", message)
	errorFlag = true
}

func error(line int, message string) {
	report(line, "", message)
}

func TokenError(token Token, message string) {
	if token.tokenType == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func RuntimeError(token Token, message string) {
	report(token.line, " somewhere ", message)
}

func main() {

	if len(os.Args) == 1 {
		runPrompt()
		return
	}

	var filepath = os.Args[1]
	runFile(filepath)

}


