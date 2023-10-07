package automaton

import (
	"errors"
)

type Action int

type Location int

var ErrIOTransitionGuardDidNotEvaluateToABoolean = errors.New("the interpretation result of the guard was not a boolean")

type IOTransition struct {
	Destination Location
	Output      Action
	Guard       Expression
	Update      Statement
}

type IOEdge struct {
	IOTransition
	Source Location
	Input  Action
}

func (transition IOTransition) GetOutput() (Action, bool) {
	return transition.Output, transition.Output > 0
}

func (transition IOTransition) Enabled(environment Environment) (bool, error) {
	if transition.Guard == nil {
		return true, nil
	}

	interpreter := Interpreter{Environment: environment}
	value, err := interpreter.Evaluate(transition.Guard)
	if err != nil {
		return false, err
	}

	if !value.IsBoolean() {
		return false, ErrIOTransitionGuardDidNotEvaluateToABoolean
	}

	return value.AsBoolean()
}

type Transitions[T any] interface {
	Transition(source Location, input Action) T
}

type SparseTransitions[T any] struct {
	Map map[Location]map[Action]T
}

func (sparse SparseTransitions[IOTransition]) Transition(source Location, input Action) IOTransition {
	return sparse.Map[source][input]
}

type IOAutomaton struct {
	locations       []Location
	inputs          []Action
	outputs         []Action
	transitions     Transitions[IOTransition]
	initialLocation Location
	finals          []Location
}

func CreateIO(
	locations []Location,
	inputs []Action,
	outputs []Action,
	transitions Transitions[IOTransition],
	initial Location,
	finals []Location,
) IOAutomaton {
	return IOAutomaton{
		locations,
		inputs,
		outputs,
		transitions,
		initial,
		finals,
	}
}

func (automaton IOAutomaton) InitialLocation() Location {
	return automaton.initialLocation
}

func (automaton IOAutomaton) Transition(source Location, input Action) IOTransition {
	return automaton.transitions.Transition(source, input)
}

func (automaton IOAutomaton) IsFinal(location Location) bool {
	for idx := range automaton.finals {
		if automaton.finals[idx] == location {
			return true
		}
	}
	return false
}
