package main

// class LoxClass {
type LoxClass struct {
	name       string
	methods    map[string]*LoxFunction
	superclass *LoxClass
}

func newLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name: name, methods: methods, superclass: superclass}
}

func (lc *LoxClass) String() string {
	return lc.name
}

func (lc *LoxClass) call(interpreter *Interpreter, arguments []any) any {
	var instance *LoxInstance = newLoxInstance(lc)
	var initializer *LoxFunction = lc.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}
	return instance
}

func (lc *LoxClass) arity() int {
	var initializer *LoxFunction = lc.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

func (lc *LoxClass) findMethod(name string) *LoxFunction {
	_, ok := lc.methods[name]
	if ok {
		return lc.methods[name]
	}
	if lc.superclass != nil {
		return lc.superclass.findMethod(name)
	}
	return nil
}
