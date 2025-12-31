package main

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func newEnvironment(environment ...*Environment) *Environment {
	var values = map[string]any{}
	if len(environment) > 0 {
		return &Environment{values, environment[0]}
	}
	return &Environment{values, nil}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}
func (e *Environment) get(name Token) any {
	var value, ok = e.values[name.lexeme]
	if ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}
	RuntimeError(name, "Undefined variable '"+name.lexeme+"'.")
	return nil
}

func (e *Environment) assign(name Token, value any) {
	_, ok := e.values[name.lexeme]
	if ok {
		e.values[name.lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return
	}
	RuntimeError(name, "Undefined variable '"+name.lexeme+"'.")
}

func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	var environment *Environment = e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) assignAt(distance int, name Token, value any) {
	e.ancestor(distance).values[name.lexeme] = value
}
