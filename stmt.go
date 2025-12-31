package main

type Stmt interface {
accept(stmtVisitor) any
}

type stmtVisitor interface {
visitBlockStmt(stmt *Block) any
visitClassStmt(stmt *Class) any
visitExpressionStmt(stmt *Expression) any
visitFunctionStmt(stmt *Function) any
visitIfStmt(stmt *If) any
visitPrintStmt(stmt *Print) any
visitReturnStmt(stmt *Return) any
visitVaStmt(stmt *Va) any
visitWhileStmt(stmt *While) any
 }

type Block struct {
statements []Stmt
}

func (block_ *Block) accept(visitor stmtVisitor) any {
return visitor.visitBlockStmt(block_)
}

func newBlock(statements []Stmt, ) *Block {
	return &Block{
statements: statements,
 }
 }
type Class struct {
name Token
superclass *Variable
methods []*Function
}

func (class_ *Class) accept(visitor stmtVisitor) any {
return visitor.visitClassStmt(class_)
}

func newClass(name Token, superclass *Variable, methods []*Function, ) *Class {
	return &Class{
name: name,
superclass: superclass,
methods: methods,
 }
 }
type Expression struct {
expression Expr
}

func (expression_ *Expression) accept(visitor stmtVisitor) any {
return visitor.visitExpressionStmt(expression_)
}

func newExpression(expression Expr, ) *Expression {
	return &Expression{
expression: expression,
 }
 }
type Function struct {
name Token
params []Token
body []Stmt
}

func (function_ *Function) accept(visitor stmtVisitor) any {
return visitor.visitFunctionStmt(function_)
}

func newFunction(name Token, params []Token, body []Stmt, ) *Function {
	return &Function{
name: name,
params: params,
body: body,
 }
 }
type If struct {
condition Expr
thenBranch Stmt
elseBranch Stmt
}

func (if_ *If) accept(visitor stmtVisitor) any {
return visitor.visitIfStmt(if_)
}

func newIf(condition Expr, thenBranch Stmt, elseBranch Stmt, ) *If {
	return &If{
condition: condition,
thenBranch: thenBranch,
elseBranch: elseBranch,
 }
 }
type Print struct {
expression Expr
}

func (print_ *Print) accept(visitor stmtVisitor) any {
return visitor.visitPrintStmt(print_)
}

func newPrint(expression Expr, ) *Print {
	return &Print{
expression: expression,
 }
 }
type Return struct {
keyword Token
value Expr
}

func (return_ *Return) accept(visitor stmtVisitor) any {
return visitor.visitReturnStmt(return_)
}

func newReturn(keyword Token, value Expr, ) *Return {
	return &Return{
keyword: keyword,
value: value,
 }
 }
type Va struct {
name Token
initializer Expr
}

func (va_ *Va) accept(visitor stmtVisitor) any {
return visitor.visitVaStmt(va_)
}

func newVa(name Token, initializer Expr, ) *Va {
	return &Va{
name: name,
initializer: initializer,
 }
 }
type While struct {
condition Expr
body Stmt
}

func (while_ *While) accept(visitor stmtVisitor) any {
return visitor.visitWhileStmt(while_)
}

func newWhile(condition Expr, body Stmt, ) *While {
	return &While{
condition: condition,
body: body,
 }
 }

