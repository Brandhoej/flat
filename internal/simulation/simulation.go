package simulation

import (
	"errors"

	"github.com/brandhoej/flat/pkg/automaton"
	"github.com/brandhoej/flat/pkg/event"
)

var (
	ErrEventAssertionFailed = errors.New("an event assertion failed")
	ErrMissingUseCaseEvent  = errors.New("the event from the previously performed use case has not yet been received")
	ErrNoEnabledEdges       = errors.New("there are no enabled edges to traverse")
	ErrUnknownRole          = errors.New("the automaton location does not match an identity role")
	ErrUnknownUseCase       = errors.New("the automaton input action does not match a use case ID")
	ErrEventTestFailed      = errors.New("an event test failed")
)

const (
	// Actors:
	unknownActor = "unknown"
	guestActor   = "guest"
	userActor    = "user"
	// Sign in:
	signInUseCaseID = event.Index(0)
	signInUseCase   = "sign in"
	signInEvent     = "signed in"
	// Sign up:
	signUpUseCaseID      = event.Index(1)
	signUpUseCase        = "sign up"
	signUpEvent          = "signed up"
	isSignedUpIdentifier = "is signed up"
	// Accept sign up:
	acceptSignUpUseCaseID = event.Index(2)
	acceptSignUpUseCase   = "accept sign up"
	acceptSignUpEvent     = "accepted sign up"
	// Sign out:
	signOutUseCaseID = event.Index(3)
	signOutUseCase   = "sign out"
	signOutEvent     = "signed out"
)

type Simulation struct {
	stepCounter    uint
	symbols        automaton.SymbolTable
	inputToUseCase map[string]event.Index
	simulator      automaton.IOSimulator
	tester         Tester
}

func Start(bus event.Bus) *Simulation {
	symbols := automaton.NewSymbolTable()
	environment := automaton.NewEnvironment()

	isSignedUp := automaton.VariableExpression{
		Symbol: automaton.Symbol(symbols.Store(isSignedUpIdentifier)),
	}
	environment.Store(isSignedUp.Symbol, automaton.CreateBoolean(false))

	canSignUp := automaton.UnaryExpression{
		Operator:   automaton.LogicalNegation,
		Expression: isSignedUp,
	}
	setSignedUpTrue := automaton.AssignmentStatement{
		Variable: isSignedUp,
		Expression: automaton.ConstantExpression{
			Value: automaton.CreateBoolean(true),
		},
	}

	model := automaton.CreateBuilder(
		unknownActor,
		symbols,
	).AddUseCase(
		unknownActor,
		signInUseCase,
		automaton.WithDestination(guestActor),
		automaton.ExpectEvent(signInEvent),
	).AddUseCase(
		guestActor,
		signUpUseCase,
		automaton.ExpectEvent(signUpEvent),
		automaton.WithGuard(canSignUp),
		automaton.WithUpdate(setSignedUpTrue),
	).AddUseCase(
		guestActor,
		acceptSignUpUseCase,
		automaton.ExpectEvent(acceptSignUpEvent),
	).AddUseCase(
		guestActor,
		signInUseCase,
		automaton.ExpectEvent(signInEvent),
	).AddUseCase(
		guestActor,
		signOutUseCase,
		automaton.ExpectEvent(signOutEvent),
	).Build()

	simulator := automaton.CreateIOSimulator(model, environment)

	simulation := &Simulation{
		stepCounter: 0,
		symbols:     symbols,
		inputToUseCase: map[string]event.Index{
			signInUseCase:       signInUseCaseID,
			signUpUseCase:       signUpUseCaseID,
			acceptSignUpUseCase: acceptSignUpUseCaseID,
			signOutUseCase:      signOutUseCaseID,
		},
		simulator: simulator,
		tester:    NewHTTPTester(),
	}

	bus.Subscribe(simulation.tester)

	return simulation
}

func (simulation *Simulation) useCaseIDForEdge(edge automaton.IOEdge) (event.Index, bool) {
	// The builder gives actions their ID as the symbol.
	inputSymbol := int(edge.Input)
	inputName, exists := simulation.symbols.Lookup(inputSymbol)
	if !exists {
		return event.Index(0), false
	}

	// In order to follow the Open Closed Principle we use a map and not a switch like structure.
	useCaseID, exists := simulation.inputToUseCase[inputName]
	if !exists {
		return event.Index(0), false
	}

	return useCaseID, true
}

func (simulation *Simulation) Step() (bool, error) {
	initialStep := simulation.stepCounter == 0

	// We cannot advance by a step if we have not yet received the expected event.
	receivedEvent, eventError := simulation.tester.ReceivedUseCaseEvent()
	if !initialStep && !receivedEvent {
		return false, ErrMissingUseCaseEvent
	}

	if eventError != nil {
		return false, errors.Join(ErrEventTestFailed, eventError)
	}

	// We use the general IO automaton simulator to get the next step in the simulation.
	edge, done := simulation.simulator.Step()
	if done {
		return false, ErrNoEnabledEdges
	}

	simulation.stepCounter += 1

	/* The simulation is responsible for handling the automaton simulator.
	 *   But also maps automaton specific information to events. */
	useCaseID, exists := simulation.useCaseIDForEdge(edge)
	if !exists {
		return true, ErrUnknownUseCase
	}

	// Last we perform a positive test which is the test we expect.
	if err := simulation.tester.Test(useCaseID); err != nil {
		return true, err
	}

	return true, nil
}
