package automaton

import (
	"math/rand"
)

type State struct {
	Location    Location
	Environment Environment
}

type IOSimulator struct {
	State     State
	Automaton IOAutomaton
}

func CreateIOSimulator(
	automaton IOAutomaton, environment Environment,
) IOSimulator {
	return IOSimulator{
		State: State{
			Location:    automaton.initialLocation,
			Environment: environment,
		},
		Automaton: automaton,
	}
}

func (simulator *IOSimulator) edges() []IOEdge {
	var edges []IOEdge

	for _, input := range simulator.Automaton.inputs {
		transition := simulator.Automaton.Transition(
			simulator.State.Location, input,
		)

		if enabled, _ := transition.Enabled(simulator.State.Environment); enabled {
			edge := IOEdge{
				IOTransition: transition,
				Source:       simulator.State.Location,
				Input:        input,
			}
			edges = append(edges, edge)
		}
	}

	return edges
}

func (simulator *IOSimulator) transition(transition IOTransition) {
	interpreter := Interpreter{Environment: simulator.State.Environment}

	// Execute the update and update the state environment.
	if transition.Update != nil {
		interpreter.Execute(transition.Update)
		simulator.State.Environment = interpreter.Environment
	}

	simulator.State.Location = transition.Destination
}

func (simulator *IOSimulator) Step() (IOEdge, bool) {
	var edge IOEdge

	// Step 1: Get all enabled edges.
	edges := simulator.edges()
	if len(edges) == 0 {
		return edge, true
	}

	// Step 2: Choose a random edge.
	edge = edges[rand.Intn(len(edges))]

	// Step 3: Transition and update state.
	simulator.transition(edge.IOTransition)

	return edge, false
}
