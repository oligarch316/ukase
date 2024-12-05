package ukmeta

import (
	"errors"
	"slices"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukinit"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Proxy
// =============================================================================

func Proxy[E Event](hook func(E) error) ukcli.Proxy {
	return Hook[E](hook)
}

func ProxyState[E Event](state ukcli.State, hook func(E) error) ukcli.State {
	return Hook[E](hook).UkaseProxy(state)
}

// -----------------------------------------------------------------------------
// Event
// -----------------------------------------------------------------------------

var (
	_ Event = EventExec{}
	_ Event = EventInfo{}
	_ Event = EventRule{}
)

type Event interface{ ukmetaEvent() }

type EventExec struct {
	Exec   ukcore.Exec
	Spec   ukspec.Parameters
	Target []string
}

type EventInfo struct {
	Info   any
	Target []string
}

type EventRule struct {
	Rule ukinit.Rule
}

func (EventExec) ukmetaEvent() {}
func (EventInfo) ukmetaEvent() {}
func (EventRule) ukmetaEvent() {}

// -----------------------------------------------------------------------------
// Hook
// -----------------------------------------------------------------------------

var _ ukcli.State = hookState{}

type Hook[E Event] func(E) error

func (h Hook[E]) UkaseProxy(s ukcli.State) ukcli.State {
	return hookState{State: s, hook: h.hook}
}

func (h Hook[E]) hook(event Event) error {
	if e, ok := event.(E); ok {
		return h(e)
	}

	return nil
}

type hookState struct {
	ukcli.State
	hook func(Event) error
}

func (hs hookState) AddExec(exec ukcore.Exec, spec ukspec.Parameters, target []string) error {
	event := EventExec{Exec: exec, Spec: spec, Target: target}
	hookErr := hs.hook(event)
	baseErr := hs.State.AddExec(exec, spec, target)
	return errors.Join(hookErr, baseErr)
}

func (hs hookState) AddInfo(info any, target []string) error {
	event := EventInfo{Info: info, Target: target}
	hookErr := hs.hook(event)
	baseErr := hs.State.AddInfo(info, target)
	return errors.Join(hookErr, baseErr)
}

func (hs hookState) AddRule(rule ukinit.Rule) error {
	event := EventRule{Rule: rule}
	hookErr := hs.hook(event)
	baseErr := hs.State.AddRule(rule)
	return errors.Join(hookErr, baseErr)
}

// =============================================================================
// Record
// =============================================================================

type Record[E Event] []E

func (r *Record[E]) UkaseProxy(state ukcli.State) ukcli.State {
	hook := func(e E) error { r.add(e); return nil }
	return ProxyState(state, hook)
}

func (r *Record[E]) add(event E) { *r = append(*r, event) }

// =============================================================================
// Auto
// =============================================================================

type Auto func(reference []string) ukcli.Directive

func (a Auto) UkaseProxy(state ukcli.State) ukcli.State {
	memo := make(autoMemo)
	hook := func(e EventExec) error { return a.add(state, memo, e.Target) }
	return ProxyState(state, hook)
}

func (a Auto) add(state ukcli.State, memo autoMemo, target []string) error {
	for _, reference := range memo.sift(target) {
		directive := a(reference)
		if err := directive.UkaseDirective(state); err != nil {
			return err
		}
	}

	return nil
}

type autoMemo map[string]autoMemo

// Ensure each sub-path of the given target is marked as visited.
// Return a list of those that have not previously been visited.
func (am autoMemo) sift(target []string) [][]string {
	const rootName = ""

	var paths [][]string

	if _, seen := am[rootName]; !seen {
		// Root (empty) target not yet visited

		// ⇒ Initialize memo to mark as visited
		am[rootName] = make(autoMemo)

		// ⇒ Add empty path (nil) to the result list
		paths = [][]string{nil}
	}

	for cur, i := am, 0; i < len(target); i++ {
		name := target[i]

		next, seen := cur[name]
		if !seen {
			// Path [0...1] not yet visited

			// ⇒ Initialize and add to memo to mark as visited
			next = make(autoMemo)
			cur[name] = next

			// ⇒ Add path (shallow copy) to the result list
			path := slices.Clone(target[:i+1])
			paths = append(paths, path)
		}

		cur = next
	}

	return paths
}
