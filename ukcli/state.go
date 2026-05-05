package ukcli

import (
	"errors"
	"reflect"

	"github.com/oligarch316/ukase/ukcore/ukinput"
	"github.com/oligarch316/ukase/ukcore/uktodo"
	"github.com/oligarch316/ukase/ukcore/uktree"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// State
// =============================================================================

type State struct {
	Global StateEntry
	Local  *uktree.Node[StateEntry]
	Values Values
}

type StateEntry struct {
	Exec      StateExec
	Info      string
	Overrides StateOverrides
}

func NewState(directives ...Directive) (State, error) {
	local := uktree.New[StateEntry]()
	state := State{Local: local}
	errs := []error{}

	for _, directive := range directives {
		err := directive.UkaseDirective(&state)
		errs = append(errs, err)
	}

	return state, errors.Join(errs...)
}

// -----------------------------------------------------------------------------
// Exec
// -----------------------------------------------------------------------------

type StateExec struct {
	exec  Exec
	execT reflect.Type
}

func (se *StateExec) Update(t reflect.Type, exec Exec) {
	se.exec, se.execT = exec, t
}

// -----------------------------------------------------------------------------
// Overrides
// -----------------------------------------------------------------------------

type StateOverrides struct {
	EnvNames  StateOverride[[]string]
	FlagNames StateOverride[[]string]
	MetaNames StateOverride[[]string]
	ArgRange  StateOverride[uktodo.ArgRange]
	FlagKind  StateOverride[ukinput.FlagKind]
	Initial   StateOverride[any]
	Info      StateOverride[string]
}

type StateOverride[T any] struct{ vMap ukvalue.Map[T] }

func (so *StateOverride[T]) Update(key ukvalue.MapKey, val T) {
	if so.vMap == nil {
		so.vMap = make(ukvalue.Map[T])
	}

	so.vMap[key] = val
}
