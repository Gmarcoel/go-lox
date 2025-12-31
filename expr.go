package main

type Expr interface {
accept(exprVisitor) any
}

type exprVisitor interface {
visitAssignExpr(expr *Assign) any
visitBinaryExpr(expr *Binary) any
visitCallExpr(expr *Call) any
visitGetExpr(expr *Get) any
visitGroupingExpr(expr *Grouping) any
visitLiteralExpr(expr *Literal) any
visitLogicalExpr(expr *Logical) any
visitSetExpr(expr *Set) any
visitSuperExpr(expr *Super) any
visitThisExpr(expr *This) any
visitUnaryExpr(expr *Unary) any
visitVariableExpr(expr *Variable) any
 }

type Assign struct {
name Token
value Expr
}

func (assign_ *Assign) accept(visitor exprVisitor) any {
return visitor.visitAssignExpr(assign_)
}

func newAssign(name Token, value Expr, ) *Assign {
	return &Assign{
name: name,
value: value,
 }
 }
type Binary struct {
left Expr
operator Token
right Expr
}

func (binary_ *Binary) accept(visitor exprVisitor) any {
return visitor.visitBinaryExpr(binary_)
}

func newBinary(left Expr, operator Token, right Expr, ) *Binary {
	return &Binary{
left: left,
operator: operator,
right: right,
 }
 }
type Call struct {
callee Expr
paren Token
arguments []any
}

func (call_ *Call) accept(visitor exprVisitor) any {
return visitor.visitCallExpr(call_)
}

func newCall(callee Expr, paren Token, arguments []any, ) *Call {
	return &Call{
callee: callee,
paren: paren,
arguments: arguments,
 }
 }
type Get struct {
object Expr
name Token
}

func (get_ *Get) accept(visitor exprVisitor) any {
return visitor.visitGetExpr(get_)
}

func newGet(object Expr, name Token, ) *Get {
	return &Get{
object: object,
name: name,
 }
 }
type Grouping struct {
expression Expr
}

func (grouping_ *Grouping) accept(visitor exprVisitor) any {
return visitor.visitGroupingExpr(grouping_)
}

func newGrouping(expression Expr, ) *Grouping {
	return &Grouping{
expression: expression,
 }
 }
type Literal struct {
value any
}

func (literal_ *Literal) accept(visitor exprVisitor) any {
return visitor.visitLiteralExpr(literal_)
}

func newLiteral(value any, ) *Literal {
	return &Literal{
value: value,
 }
 }
type Logical struct {
left Expr
operator Token
right Expr
}

func (logical_ *Logical) accept(visitor exprVisitor) any {
return visitor.visitLogicalExpr(logical_)
}

func newLogical(left Expr, operator Token, right Expr, ) *Logical {
	return &Logical{
left: left,
operator: operator,
right: right,
 }
 }
type Set struct {
object Expr
name Token
value Expr
}

func (set_ *Set) accept(visitor exprVisitor) any {
return visitor.visitSetExpr(set_)
}

func newSet(object Expr, name Token, value Expr, ) *Set {
	return &Set{
object: object,
name: name,
value: value,
 }
 }
type Super struct {
keyword Token
method Token
}

func (super_ *Super) accept(visitor exprVisitor) any {
return visitor.visitSuperExpr(super_)
}

func newSuper(keyword Token, method Token, ) *Super {
	return &Super{
keyword: keyword,
method: method,
 }
 }
type This struct {
keyword Token
}

func (this_ *This) accept(visitor exprVisitor) any {
return visitor.visitThisExpr(this_)
}

func newThis(keyword Token, ) *This {
	return &This{
keyword: keyword,
 }
 }
type Unary struct {
operator Token
right Expr
}

func (unary_ *Unary) accept(visitor exprVisitor) any {
return visitor.visitUnaryExpr(unary_)
}

func newUnary(operator Token, right Expr, ) *Unary {
	return &Unary{
operator: operator,
right: right,
 }
 }
type Variable struct {
name Token
}

func (variable_ *Variable) accept(visitor exprVisitor) any {
return visitor.visitVariableExpr(variable_)
}

func newVariable(name Token, ) *Variable {
	return &Variable{
name: name,
 }
 }

