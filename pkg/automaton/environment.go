package automaton

type Environment struct {
	bindings map[Symbol]Value
}

func NewEnvironment() Environment {
	return Environment{
		bindings: map[Symbol]Value{},
	}
}

func (environment Environment) Contains(symbol Symbol) bool {
	_, exists := environment.bindings[symbol]
	return exists
}

func (environment Environment) Store(symbol Symbol, binding Value) {
	environment.bindings[symbol] = binding
}

func (environment Environment) Lookup(symbol Symbol) (Value, bool) {
	value, exists := environment.bindings[symbol]
	return value, exists
}
