package automaton

type ioTransitionsConfiguration struct {
	destination string
	guard       Expression
	update      Statement
	output      string
}

func (config *ioTransitionsConfiguration) configure(options ...ioTransitionsOption) {
	for idx := range options {
		options[idx](config)
	}
}

type ioTransitionsOption func(*ioTransitionsConfiguration)

func WithDestination(destination string) ioTransitionsOption {
	return func(configuration *ioTransitionsConfiguration) {
		configuration.destination = destination
	}
}

func WithGuard(guard Expression) ioTransitionsOption {
	return func(configuration *ioTransitionsConfiguration) {
		configuration.guard = guard
	}
}

func WithUpdate(update Statement) ioTransitionsOption {
	return func(configuration *ioTransitionsConfiguration) {
		configuration.update = update
	}
}

func ExpectEvent(event string) ioTransitionsOption {
	return func(configuration *ioTransitionsConfiguration) {
		configuration.output = event
	}
}

type Builder struct {
	symbols         SymbolTable
	initialLocation Location
	actors          map[string]Location
	transitions     SparseTransitions[IOTransition]
}

func CreateBuilder(initial string, symbols SymbolTable) *Builder {
	builder := &Builder{
		symbols: symbols,
		actors:  map[string]Location{},
		transitions: SparseTransitions[IOTransition]{
			Map: make(map[Location]map[Action]IOTransition),
		},
	}
	initialLocation := builder.addActor(initial)
	builder.initialLocation = initialLocation
	return builder
}

func (builder *Builder) addActor(name string) Location {
	if location, exists := builder.actors[name]; exists {
		return location
	}

	location := Location(builder.symbols.Store(name))
	builder.actors[name] = location
	builder.transitions.Map[location] = make(map[Action]IOTransition)
	return location
}

func (builder *Builder) AddActor(name string) *Builder {
	builder.addActor(name)
	return builder
}

func (builder *Builder) AddUseCase(
	actor string, useCase string, options ...ioTransitionsOption,
) *Builder {
	var configuration ioTransitionsConfiguration
	configuration.destination = actor
	configuration.configure(options...)

	transition := IOTransition{}

	source, exists := builder.actors[actor]
	if !exists {
		source = builder.addActor(actor)
	}

	transition.Destination = source
	if configuration.destination == "" {
		transition.Destination = builder.addActor(configuration.destination)
	}

	if configuration.output != "" {
		transition.Output = Action(builder.symbols.Store(configuration.output))
	}

	transition.Guard = configuration.guard
	transition.Update = configuration.update

	input := Action(builder.symbols.Store(useCase))

	builder.transitions.Map[source][input] = transition

	return builder
}

func (builder *Builder) Build() IOAutomaton {
	var locations []Location
	var finals []Location

	for _, location := range builder.actors {
		locations = append(locations, location)
		finals = append(finals, location)
	}

	var inputs []Action
	var outputs []Action

	for _, transitions := range builder.transitions.Map {
		for input, transition := range transitions {
			inputs = append(inputs, input)
			outputs = append(outputs, transition.Output)
		}
	}

	return CreateIO(
		locations, inputs,
		outputs,
		builder.transitions,
		builder.initialLocation,
		finals,
	)
}
