package main

import (
	"fmt"
	"reflect"
	"time"
)

type LoxCallable interface {
	call(interpreter *Interpreter, arguments []any) any
	arity() int
}

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
}

type clockNative struct{}

func (c *clockNative) arity() int32 {
	return 0
}

func (c *clockNative) call(i *Interpreter, arguments []any) float32 {
	return float32(time.Now().UnixNano()) / 1e9
}
func (c *clockNative) String() string {
	return "<native fn>"
}

func newInterpreter() *Interpreter {
	var globals = newEnvironment()
	var environment = globals

	globals.define("clock", clockNative{})

	return &Interpreter{globals, environment, map[Expr]int{}}
}

// visit statements
func (i *Interpreter) visitClassStmt(stmt *Class) any {
	var superclass any = nil
	if stmt.superclass != nil {
		superclass = i.evaluate(stmt.superclass)
		var _, ok = superclass.(*LoxClass)
		if !ok {

			RuntimeError(stmt.superclass.name, "Superclass must be a class.")
		}
	}
	i.environment.define(stmt.name.lexeme, nil)
	if stmt.superclass != nil {
		i.environment = newEnvironment(i.environment)
		i.environment.define("super", superclass)
	}
	var methods = map[string]*LoxFunction{}
	for _, method := range stmt.methods {
		var isInitializer bool = method.name.lexeme == "init"
		var function *LoxFunction = newLoxFunction(method, i.environment, isInitializer)
		methods[method.name.lexeme] = function
	}
	var klass any
	if superclass == nil {
		klass = newLoxClass(stmt.name.lexeme, nil, methods)
	} else {
		klass = newLoxClass(stmt.name.lexeme, superclass.(*LoxClass), methods)
		i.environment = i.environment.enclosing
	}
	i.environment.assign(stmt.name, klass)
	return nil
}
func (i *Interpreter) visitVaStmt(stmt *Va) any {
	var value any = nil
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}
	i.environment.define(stmt.name.lexeme, value)
	return nil
}

func (i *Interpreter) visitExpressionStmt(stmt *Expression) any {
	i.evaluate(stmt.expression)
	return nil
}
func (i *Interpreter) visitPrintStmt(stmt *Print) any {
	var value = i.evaluate(stmt.expression)
	fmt.Println(fmt.Sprint(value))
	return nil
}

func (i *Interpreter) visitBlockStmt(stmt *Block) any {
	return i.executeBlock(stmt.statements, newEnvironment(i.environment))
}

func (i *Interpreter) visitIfStmt(stmt *If) any {
	var ret_value any = nil
	if i.isTruthy(i.evaluate(stmt.condition)) {
		ret_value = i.execute(stmt.thenBranch)
	} else if stmt.elseBranch != nil {
		ret_value = i.execute(stmt.elseBranch)
	}
	return ret_value
}

func (i *Interpreter) visitWhileStmt(stmt *While) any {
	var ret_value any = nil
	for i.isTruthy(i.evaluate(stmt.condition)) {
		ret_value = i.execute(stmt.body)
		if ret_value != nil {
			break
		}
	}
	return ret_value
}

func (i *Interpreter) visitFunctionStmt(stmt *Function) any {
	var function = newLoxFunction(stmt, i.environment, false)
	i.environment.define(stmt.name.lexeme, function)
	return nil
}

func (i *Interpreter) visitReturnStmt(stmt *Return) any {
	var value any = nil
	if stmt.value != nil {
		value = i.evaluate(stmt.value)
	}
	return value
}

// visit expressions
func (i *Interpreter) visitAssignExpr(expr *Assign) any {
	var value any = i.evaluate(expr.value)
	distance, ok := i.locals[expr]
	if ok {
		i.environment.assignAt(distance, expr.name, value)
	} else {
		i.globals.assign(expr.name, value)
	}
	return value
}

func (i *Interpreter) visitVariableExpr(expr *Variable) any {
	return i.lookUpVariable(expr.name, expr)
}

func (i *Interpreter) visitLiteralExpr(expr *Literal) any {
	return expr.value
}

func (i *Interpreter) visitGroupingExpr(expr *Grouping) any {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitBinaryExpr(expr *Binary) any {
	var left = i.evaluate(expr.left)
	var right = i.evaluate(expr.right)
	var err any = nil
	switch expr.operator.tokenType {
	case GREATER:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) >= right.(float64)
	case LESS:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) <= right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	case PLUS:
		var left_num, okl = left.(float64)
		var right_num, okr = right.(float64)
		if okl && okr {
			return left_num + right_num
		}
		left_str, okl := left.(string)
		right_str, okr := right.(string)
		if okl && okr {
			return left_str + right_str
		}
		RuntimeError(expr.operator, "Operands "+fmt.Sprint(left)+" and "+fmt.Sprint(right)+" must be two numbers or two strings.")
	case MINUS:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) - right.(float64)
	case SLASH:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) / right.(float64)
	case STAR:
		err = i.checkNumberOperands(expr.operator, left, right)
		if err != nil {
			return nil
		}
		return left.(float64) * right.(float64)
	}
	return nil
}

func (i *Interpreter) visitUnaryExpr(expr *Unary) any {
	var right = i.evaluate(expr.right)
	var err any = nil
	switch expr.operator.tokenType {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		err = i.checkNumberOperand(expr.operator, right)
		if err != nil {
			return nil
		}
		return -right.(float64)
	}
	return nil
}

func (i *Interpreter) visitLogicalExpr(expr *Logical) any {
	var left = i.evaluate(expr.left)
	if expr.operator.tokenType == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.right)
}

func (i *Interpreter) visitSetExpr(expr *Set) any {
	var object = i.evaluate(expr.object)
	li_object, ok := object.(*LoxInstance)
	if !ok {
		RuntimeError(expr.name, "Only instances have fields.")
	}
	var value = i.evaluate(expr.value)
	li_object.set(expr.name, value)
	return value
}

func (i *Interpreter) visitSuperExpr(expr *Super) any {
	var distance int = i.locals[expr]
	var superclass *LoxClass = i.environment.getAt(distance, "super").(*LoxClass)
	var object *LoxInstance = i.environment.getAt(distance-1, "this").(*LoxInstance)
	var method *LoxFunction = superclass.findMethod(expr.method.lexeme)
	if method == nil {
		RuntimeError(expr.method, "Undefined property '"+expr.method.lexeme+"'.")
	}
	return method.bind(object)
}

func (i *Interpreter) visitThisExpr(expr *This) any {
	return i.lookUpVariable(expr.keyword, expr)
}

func (i *Interpreter) visitCallExpr(expr *Call) any {
	var callee = i.evaluate(expr.callee)
	var arguments = []any{}
	for _, a := range expr.arguments {
		arguments = append(arguments, i.evaluate(a.(Expr)))
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		RuntimeError(expr.paren, "Can only call functions and classes.")
		return nil
	}
	if len(arguments) != function.arity() {
		RuntimeError(expr.paren, fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments)))
		return nil
	}
	return function.call(i, arguments)
}

func (i *Interpreter) visitGetExpr(expr *Get) any {
	var object = i.evaluate(expr.object)
	li_object, ok := object.(*LoxInstance)
	if ok {
		return li_object.get(expr.name)
	}
	RuntimeError(expr.name, "Only instances have properties.")
	return nil
}

// methods
func (i *Interpreter) evaluate(expr Expr) any {
	return expr.accept(i)
}

func (i *Interpreter) isEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	var b, ok = object.(bool)
	if ok {
		return b
	}
	return true
}

func (i *Interpreter) checkNumberOperands(operator Token, left any, right any) any {
	var _, okl = left.(float64)
	var _, okr = right.(float64)
	if okl && okr {
		return nil
	}
	RuntimeError(operator, "Operands "+fmt.Sprint(left)+" and "+fmt.Sprint(right)+" must be numbers.")
	return "error"
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) any {
	var _, ok = operand.(float64)
	if ok {
		return nil
	}
	RuntimeError(operator, "Operand "+fmt.Sprint(operand)+" must be a number.")
	return "error"
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
	// fmt.Println(fmt.Sprint(value))
}

func (i *Interpreter) execute(stmt Stmt) any {
	return stmt.accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, env *Environment) any {
	var ret_value any = nil
	var previous = i.environment
	i.environment = env
	for _, stmt := range statements {
		ret_value = i.execute(stmt)
		if ret_value != nil {
			break
		}
	}

	i.environment = previous
	return ret_value
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) any {
	var distance, ok = i.locals[expr]
	if ok {
		return i.environment.getAt(distance, name.lexeme)
	} else {
		return i.globals.get(name)
	}
}
