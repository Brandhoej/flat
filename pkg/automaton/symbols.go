package automaton

type Symbol int32

type SymbolTable struct {
	identifiers map[string]int
	symbols     map[int]string
}

func NewSymbolTable() SymbolTable {
	return SymbolTable{
		identifiers: make(map[string]int),
		symbols:     make(map[int]string),
	}
}

func (table *SymbolTable) Store(identifier string) int {
	// If the value already exists then we do nothing and return the existing symbol.
	// This gurantees that the returned symbol is always pointing at the same identifier.
	// If we did not have this rule and stored the same identifier twice then could result
	// in cases where lookup for the symbol returned a different identifier.
	if value, exists := table.identifiers[identifier]; exists {
		return value
	}

	// We start at one so the zero'th value can be reserved as a replacement for nil.
	// That way we can reduce the size of our edges and transitions and not make omittable outputs a pointer.
	symbol := len(table.identifiers) + 1
	table.identifiers[identifier] = symbol
	table.symbols[symbol] = identifier
	return symbol
}

func (table *SymbolTable) Lookup(symbol int) (string, bool) {
	value, exists := table.symbols[symbol]
	return value, exists
}
