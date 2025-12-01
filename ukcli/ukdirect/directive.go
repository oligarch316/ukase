package ukdirect

import (
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore/uktree"
)

// =============================================================================
// Directive
// =============================================================================

type Func func(*ukcli.State) error

func (f Func) UkaseDirective(s *ukcli.State) error { return f(s) }

func Global(op Operation) Func {
	return func(s *ukcli.State) error {
		return op.UkaseOperation(&s.Global)
	}
}

func Local(op Operation, target ...string) Func {
	return func(s *ukcli.State) error {
		return uktree.Update(s.Local, op.UkaseOperation, target...)
	}
}
