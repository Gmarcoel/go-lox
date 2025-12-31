// class LoxInstance {
package main

type LoxInstance struct {
	klass  *LoxClass
	fields map[string]any
}

func newLoxInstance(klass *LoxClass) *LoxInstance {
	return &LoxInstance{klass: klass}
}

func (li *LoxInstance) String() string {
	return li.klass.name + " instance"
}

func (li *LoxInstance) get(name Token) any {
	// if (fields.containsKey(name.lexeme)) {
	value, ok := li.fields[name.lexeme]
	if ok {
		return value
	}
	var method = li.klass.findMethod(name.lexeme)
	if method != nil {
		return method.bind(li)
	}

	RuntimeError(name, "Undefined property '"+name.lexeme+"'.")
	return nil
}

func (li *LoxInstance) set(name Token, value any) any {
	li.fields[name.lexeme] = value
	return nil
}

