package main

// expression → equality ;
// equality → comparison ( ( "!=" | "==" ) comparison )* ;
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term → factor ( ( "-" | "+" ) factor )* ;
// factor → unary ( ( "/" | "*" ) unary )* ;
// unary → ( "!" | "-" ) unary
// | primary ;
// primary → NUMBER | STRING | "true" | "false" | "nil"
// | "(" expression ")" ;

type Parser struct {
	tokens  []Token
	current int
}

type ParseError struct {
}

func newParser(tokens []Token) *Parser {
	return &Parser{
		current: 0,
		tokens:  tokens,
	}
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return (p.peek().tokenType == t)
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}
func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}
	// throw error(peek(), message);
	p.error(p.peek(), message)
	return p.advance()
}

func (p *Parser) error(token Token, message string) ParseError {
	TokenError(token, message)
	return ParseError{}
}

// expression → equality ;
func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	var expr Expr = p.or()
	if p.match(EQUAL) {
		var equals Token = p.previous()
		var value Expr = p.assignment()
		expr_var, ok := expr.(*Variable)
		if ok {
			var name Token = expr_var.name
			return newAssign(name, value)
		} else {
			get, ok := expr.(*Get)
			if ok {
				return newSet(get.object, get.name, value)
			}
		}
		TokenError(equals, "Invalid assignment target.")
	}
	return expr
}

// equality → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	var expr Expr = p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		var operator Token = p.previous()
		var right Expr = p.comparison()
		expr = newBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) and() Expr {
	// Expr expr = equality();
	var expr Expr = p.equality()
	// while (match(AND)) {
	for p.match(AND) {
		// Token operator = previous();
		var operator = p.previous()
		// Expr right = equality();
		var right = p.equality()
		// expr = new Expr.Logical(expr, operator, right);
		expr = newLogical(expr, operator, right)
	}
	return expr
}

func (p *Parser) or() Expr {
	var expr = p.and()
	for p.match(OR) {
		var operator Token = p.previous()
		var right Expr = p.and()
		expr = newLogical(expr, operator, right)
	}
	return expr
}

// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	var expr Expr = p.term()
	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		var operator Token = p.previous()
		var right Expr = p.term()
		expr = newBinary(expr, operator, right)
	}
	return expr
}

// term → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	var expr Expr = p.factor()
	for p.match(MINUS, PLUS) {
		var operator Token = p.previous()
		var right Expr = p.factor()
		expr = newBinary(expr, operator, right)
	}
	return expr
}

// factor → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	var expr Expr = p.unary()
	for p.match(SLASH, STAR) {
		var operator Token = p.previous()
		var right Expr = p.unary()
		expr = newBinary(expr, operator, right)
	}
	return expr
}

// unary → ( "!" | "-" ) unary | primary
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		var operator = p.previous()
		var right = p.unary()
		return newUnary(operator, right)
	}
	return p.call()
}

// primary → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return newLiteral(false)
	}
	if p.match(TRUE) {
		return newLiteral(true)
	}
	if p.match(NIL) {
		return newLiteral(nil)
	}
	if p.match(NUMBER, STRING) {
		return newLiteral(p.previous().literal)
	}
	if p.match(SUPER) {
		var keyword Token = p.previous()
		p.consume(DOT, "Expect '.' after 'super'.")
		var method Token = p.consume(IDENTIFIER, "Expect superclass method name.")
		return newSuper(keyword, method)
	}
	if p.match(THIS) {
		return newThis(p.previous())
	}
	if p.match(IDENTIFIER) {
		return newVariable(p.previous())
	}
	if p.match(LEFT_PAREN) {
		var expr Expr = p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return newGrouping(expr)
	}
	if p.match(EOF) {
		return nil
	}

	TokenError(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) call() Expr {
	var expr Expr = p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			var name Token = p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = newGet(expr, name)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments = []any{}
	if !p.check(RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
		for p.match(COMMA) {
			if len(arguments) > 255 {
				p.error(p.peek(), "Can't have more than 255 argumets")
			}
			arguments = append(arguments, p.expression())
		}
	}
	var paren Token = p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return newCall(callee, paren, arguments)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}
		switch p.peek().tokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) parse() []Stmt {
	var statements = []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
	// parser.synchronize()
}

func (p *Parser) classDeclaration() Stmt {
	var name Token = p.consume(IDENTIFIER, "Expect class name.")
	// 	Expr.Variable superclass = null;
	var superclass *Variable
	// if (match(LESS)) {
	if p.match(LESS) {
		p.consume(IDENTIFIER, "Expect superclass name.")
		superclass = newVariable(p.previous())
	}
	p.consume(LEFT_BRACE, "Expect '{' before class body.")
	var methods []*Function = []*Function{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))

	}
	p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	return newClass(name, superclass, methods)
}

func (p *Parser) varDeclaration() Stmt {
	var name Token = p.consume(IDENTIFIER, "Expect variable name.")
	var initializer Expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return newVa(name, initializer)
}

func (p *Parser) function(kind string) *Function {
	var name Token = p.consume(IDENTIFIER, "Expected "+kind+" name.")
	p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name.")
	var parameters []Token = []Token{}
	if !p.check(RIGHT_PAREN) {
		parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
		for p.match(COMMA) {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body.")
	var body []Stmt = p.block()
	return newFunction(name, parameters, body)
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return newBlock(p.block())
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	var value Expr = p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return newPrint(value)
}

func (p *Parser) returnStatement() Stmt {
	var keyword Token = p.previous()
	var value Expr = nil
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value")
	return newReturn(keyword, value)
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	var condition Expr = p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	var body Stmt = p.statement()
	return newWhile(condition, body)
}

func (p *Parser) expressionStatement() Stmt {
	var expr Expr = p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return newExpression(expr)
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	var condition Expr = nil
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr = nil
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	var body = p.statement()

	if increment != nil {
		body = newBlock([]Stmt{body, newExpression(increment)})
	}

	if condition == nil {
		condition = newLiteral(true)
	}
	body = newWhile(condition, body)

	if initializer != nil {
		body = newBlock([]Stmt{initializer, body})
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	var condition Expr = p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	var thenBranch Stmt = p.statement()
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return newIf(condition, thenBranch, elseBranch)
}

func (p *Parser) block() []Stmt {
	var statements = []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}
