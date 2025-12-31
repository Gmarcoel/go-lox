package main
import (
	"unicode"
	"strconv"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func newScanner(source string) *Scanner {
	return &Scanner{source: source, tokens: []Token{}, start: 0, current: 0, line: 1}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len([]rune(s.source))
}

func (s *Scanner) scanTokens() []Token {

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, *newToken(EOF, "", "", s.line))
	return s.tokens
}

func (s *Scanner) advance() rune {
	s.current++
	return []rune(s.source)[s.current-1]
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\n'
	}
	return []rune(s.source)[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len([]rune(s.source)) {
		return '\n'
	}
	return []rune(s.source)[s.current+1]
}

func (s *Scanner) addToken(token TokenType, literal ...any) {

	var text string = string([]rune(s.source)[s.start:s.current])
	if len(literal) == 0 {
		literal = append(literal, nil)
	}
	s.tokens = append(s.tokens, *newToken(token, text, literal[0], s.line))
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if []rune(s.source)[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
			s.advance()
		}
		if s.isAtEnd() {
			error(s.line, "unterminated string")
			return
		}

		// closing
		s.advance()

	}
	s.advance()
	// trim surrounding quotes
	var value = string([]rune(s.source)[s.start+1 : s.current-1])
	s.addToken(STRING, value)
}

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}
	// Look for a fractional part.
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}
	var num_string = string([]rune(s.source)[s.start:s.current])
	var number, err = strconv.ParseFloat(num_string, 64)
	if err != nil {
		error(s.line, "float parse error scanning")
	}

	s.addToken(NUMBER, number)

	// addToken(NUMBER,
	// Double.parseDouble(source.substring(start, current)));
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) {
		s.advance()
	}
	var text = string([]rune(s.source)[s.start:s.current])
	var tokenType TokenType = keywords[text]
	if tokenType == NOT_FOUND {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType, text)
}

func (s *Scanner) scanToken() {
	var c rune = s.advance()
	switch c {
	case '(':
		{
			s.addToken(LEFT_PAREN)
			break
		}
	case ')':
		{
			s.addToken(RIGHT_PAREN)
			break
		}
	case '{':
		{
			s.addToken(LEFT_BRACE)
			break
		}
	case '}':
		{
			s.addToken(RIGHT_BRACE)
			break
		}
	case ',':
		{
			s.addToken(COMMA)
			break
		}
	case '.':
		{
			s.addToken(DOT)
			break
		}
	case '-':
		{
			s.addToken(MINUS)
			break
		}
	case '+':
		{
			s.addToken(PLUS)
			break
		}
	case ';':
		{
			s.addToken(SEMICOLON)
			break
		}
	case '*':
		{
			s.addToken(STAR)
			break
		}
	case '!':
		{
			if s.match('=') {
				s.addToken(BANG_EQUAL)
			} else {
				s.addToken(BANG)
			}
			break
		}
	case '=':
		{
			if s.match('=') {
				s.addToken(EQUAL_EQUAL)
			} else {
				s.addToken(EQUAL)
			}
			break
		}
	case '<':
		{
			if s.match('=') {
				s.addToken(LESS_EQUAL)
			} else {
				s.addToken(LESS)
			}
			break
		}
	case '>':
		{
			if s.match('=') {
				s.addToken(GREATER_EQUAL)
			} else {
				s.addToken(GREATER)
			}
			break
		}
	case '/':
		{
			if s.match('/') {
				// A comment goes until the end of the line.
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else if s.peek() == '*' {
				for {
					if s.peek() == '*' && s.peekNext() == '/' {
						s.advance()
						s.advance()
						break
					}
					s.advance()
				}
			} else {
				s.addToken(SLASH)
			}
			break
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		{
			s.line += 1
			break
		}
	case '"':
		{
			s.string()
			break
		}
	default:
		{
			if unicode.IsDigit(c) {
				s.number()
			} else if unicode.IsLetter(c) {
				s.identifier()
			} else {
				error(s.line, "unexpected character")
			}
			break
		}
	}
}
