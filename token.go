package main
import (
	"fmt"
	"strconv"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func newToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t *Token) String() string {
	return " type:" + strconv.Itoa(int(t.tokenType)) + ", lexeme: " + t.lexeme + ", literal: " + fmt.Sprint(t.literal)
}
