package main

type functionType int

const (
	NONE functionType = iota
	FUNCTION
	INITIALIZER
	METHOD
)

type classType int

const (
	NO_CLASS classType = iota
	YES_CLASS
	SUB_CLASS
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []any
	currentFunction functionType
	currentClass    classType
}

func newResolver(interpreter *Interpreter) Resolver {
	return Resolver{interpreter: interpreter, currentFunction: NONE, currentClass: NO_CLASS}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) resolve(statements any) any {
	switch t := statements.(type) {
	case []Stmt:
		for _, stmt := range t {
			r.resolve(stmt)
		}
	case Stmt:
		t.accept(r)
	case Expr:
		t.accept(r)

	}

	return nil
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		_, ok := r.scopes[i].(map[string]bool)[name.lexeme]
		if ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
		}
	}
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) {
	// if (scopes.isEmpty()) return;
	if len(r.scopes) == 0 {
		return
	}
	var scope = r.scopes[len(r.scopes)-1].(map[string]bool)
	// 	if (scope.containsKey(name.lexeme)) {
	_, ok := scope[name.lexeme]
	if ok {
		TokenError(name, "Already variable with this name in this scope.")
	}
	scope[name.lexeme] = false
	return
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1].(map[string]bool)[name.lexeme] = true
}

func (r *Resolver) resolveFunction(function *Function, ftype functionType) {
	var enclosingFunction = r.currentFunction
	r.currentFunction = ftype
	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(function.body)
	r.endScope()
	r.currentFunction = enclosingFunction
}
func (r *Resolver) visitClassStmt(stmt *Class) any {
	var enclosingClass = r.currentClass
	r.currentClass = YES_CLASS
	r.declare(stmt.name)
	r.define(stmt.name)
	if stmt.superclass != nil && stmt.name.lexeme == stmt.superclass.name.lexeme {
		TokenError(stmt.superclass.name, "A class can't inherit from itself.")
	}
	if stmt.superclass != nil {
		r.currentClass = SUB_CLASS
		r.resolve(stmt.superclass)
	}
	if stmt.superclass != nil {
		r.beginScope()
		r.scopes[len(r.scopes)-1].(map[string]bool)["super"] = true
	}
	r.beginScope()
	r.scopes[len(r.scopes)-1].(map[string]bool)["this"] = true
	for _, method := range stmt.methods {
		var declaration = METHOD
		if method.name.lexeme == "init" {
			declaration = INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}
	r.endScope()
	if stmt.superclass != nil {
		r.endScope()
	}
	r.currentClass = enclosingClass
	return nil
}

func (r *Resolver) visitBlockStmt(stmt *Block) any {
	r.beginScope()
	r.resolve(stmt.statements)
	r.endScope()
	return nil
}

func (r *Resolver) visitVaStmt(stmt *Va) any {
	r.declare(stmt.name)
	if stmt.initializer != nil {
		r.resolve(stmt.initializer)
	}
	r.define(stmt.name)
	return nil
}

func (r *Resolver) visitVariableExpr(expr *Variable) any {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1].(map[string]bool)
		val, ok := scope[expr.name.lexeme]
		if ok && val == false {
			TokenError(expr.name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.name)
	return nil
}

func (r *Resolver) visitAssignExpr(expr *Assign) any {
	r.resolve(expr.value)
	r.resolveLocal(expr, expr.name)
	return nil
}
func (r *Resolver) visitExpressionStmt(stmt *Expression) any {
	r.resolve(stmt.expression)
	return nil
}

func (r *Resolver) visitFunctionStmt(stmt *Function) any {
	r.declare(stmt.name)
	r.define(stmt.name)
	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) visitIfStmt(stmt *If) any {
	r.resolve(stmt.condition)
	r.resolve(stmt.thenBranch)
	if stmt.elseBranch != nil {
		r.resolve(stmt.elseBranch)
	}
	return nil
}
func (r *Resolver) visitPrintStmt(stmt *Print) any {
	r.resolve(stmt.expression)
	return nil
}

func (r *Resolver) visitReturnStmt(stmt *Return) any {
	if stmt.value != nil {
		if r.currentFunction == INITIALIZER {
			TokenError(stmt.keyword, "Can't return a value from an initializer.")
		}
		r.resolve(stmt.value)
	}
	return nil
}
func (r *Resolver) visitWhileStmt(stmt *While) any {
	r.resolve(stmt.condition)
	r.resolve(stmt.body)
	return nil
}

func (r *Resolver) visitBinaryExpr(expr *Binary) any {
	r.resolve(expr.left)
	r.resolve(expr.right)
	return nil
}

func (r *Resolver) visitCallExpr(expr *Call) any {
	r.resolve(expr.callee)
	for _, argument := range expr.arguments {
		r.resolve(argument)
	}
	return nil
}

func (r *Resolver) visitGetExpr(expr *Get) any {
	r.resolve(expr.object)
	return nil
}

func (r *Resolver) visitSetExpr(expr *Set) any {
	r.resolve(expr.value)
	r.resolve(expr.object)
	return nil
}

func (r *Resolver) visitThisExpr(expr *This) any {
	if r.currentClass == NO_CLASS {
		TokenError(expr.keyword, "Can't use 'this' outside of a class.")
		return nil
	}
	r.resolveLocal(expr, expr.keyword)
	return nil
}

func (r *Resolver) visitGroupingExpr(expr *Grouping) any {
	r.resolve(expr.expression)
	return nil
}

func (r *Resolver) visitLiteralExpr(expr *Literal) any {
	return nil
}

func (r *Resolver) visitLogicalExpr(expr *Logical) any {
	r.resolve(expr.left)
	r.resolve(expr.right)
	return nil
}

func (r *Resolver) visitUnaryExpr(expr *Unary) any {
	r.resolve(expr.right)
	return nil
}

func (r *Resolver) visitSuperExpr(expr *Super) any {
	if r.currentClass == NO_CLASS {
		TokenError(expr.keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != SUB_CLASS {
		TokenError(expr.keyword, "Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.keyword)
	return nil
}
