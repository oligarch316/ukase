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
	return func(state *ukcli.State) error {
		return op.UkaseOperation(&state.Global)
	}
}

func Local(op Operation, target ...string) Func {
	return func(state *ukcli.State) error {
		return uktree.Update(state.Local, op.UkaseOperation, target...)
	}
}

func Auto(op Operation, name string) Func {
	return func(state *ukcli.State) error {
		queue := []*uktree.Node[ukcli.StateEntry]{state.Local}

		for node := queue[0]; len(queue) > 0; node, queue = queue[0], queue[1:] {
			if err := uktree.Update(node, op.UkaseOperation, name); err != nil {
				return err
			}

			for childName, childNode := range node.Children {
				if childName != name {
					queue = append(queue, childNode)
				}
			}
		}

		return nil
	}
}
