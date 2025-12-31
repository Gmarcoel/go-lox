package main

type LoxFunction struct {
	declaration   *Function
	closure       *Environment
	isInitializer bool
}

func (lf *LoxFunction) arity() int {
	return len(lf.declaration.params)
}

func (lf *LoxFunction) String() string {
	return "<fn " + lf.declaration.name.lexeme + ">"
}

func (lf *LoxFunction) call(i *Interpreter, arguments []any) any {
	var environment *Environment = newEnvironment(lf.closure)
	for i := 0; i < len(lf.declaration.params); i++ {
		var value = arguments[i]
		environment.define(lf.declaration.params[i].lexeme, value)
	}
	var ret_value = i.executeBlock(lf.declaration.body, environment)
	if ret_value != nil {
		if lf.isInitializer {
			return lf.closure.getAt(0, "this")
		}
		return ret_value
	}
	if lf.isInitializer {
		return lf.closure.getAt(0, "this")
	}
	return nil
}

func (lf *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	var environment *Environment = newEnvironment(lf.closure)
	environment.define("this", instance)
	return newLoxFunction(lf.declaration, environment, lf.isInitializer)
}

func newLoxFunction(declaration *Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{declaration: declaration, closure: closure, isInitializer: isInitializer}
}
