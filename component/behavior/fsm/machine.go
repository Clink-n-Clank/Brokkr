package fsm

// State identifier, must be unique
type State interface {
	~byte | ~int | string
}

// StateTransitionHandlerFunc that handling transitions between states and returns new
type StateTransitionHandlerFunc[S State, A any] func(A) (S, A)

// Machine it lets change a behavior when the state changes
type Machine[S State, A any] struct {
	stateStart    S
	stateEnding   map[S]bool
	stateHandlers map[S]StateTransitionHandlerFunc[S, A]
}

// AddState and associated with handler function for transition
func (m *Machine[S, A]) AddState(state S, handler StateTransitionHandlerFunc[S, A]) {
	if m.stateHandlers == nil {
		m.stateHandlers = make(map[S]StateTransitionHandlerFunc[S, A], 0)
	}

	m.stateHandlers[state] = handler
}

// AddEndState for transition applies for handler and terminates when the end state is reached
func (m *Machine[S, A]) AddEndState(state S) {
	if m.stateEnding == nil {
		m.stateEnding = make(map[S]bool, 0)
	}

	m.stateEnding[state] = true
}

// MakeTransition with given arguments
func (m *Machine[S, A]) MakeTransition(handlerArgs A) {
	stateHandler, isStateExist := m.stateHandlers[m.stateStart]
	if !isStateExist {
		return
	}

	for {
		if stateHandler == nil {
			break
		}

		nextState, nextArgs := stateHandler(handlerArgs)
		_, finished := m.stateEnding[nextState]

		if finished {
			break
		} else {
			stateHandler, isStateExist = m.stateHandlers[nextState]
			handlerArgs = nextArgs
		}
	}
}
