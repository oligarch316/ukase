package ukexec

import (
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// TODO: Document
type Entry interface {
	// TODO: Document
	Exec() (ukcore.Exec, bool)

	// TODO: Document
	Info() (any, bool)

	// TODO: Document
	Spec() (ukspec.Parameters, bool)

	// TODO: Document
	Child(target ...string) (Entry, bool)

	// TODO: Document
	Children() map[string]Entry
}

type entry struct {
	exec ukcore.Exec
	info any
	spec ukspec.Parameters

	children map[string]*entry
	flags    map[string]ukspec.Flag
}

func newEntry() *entry {
	return &entry{
		children: make(map[string]*entry),
		flags:    make(map[string]ukspec.Flag),
	}
}

func (e *entry) Exec() (ukcore.Exec, bool)       { return e.exec, e.exec != nil }
func (e *entry) Info() (any, bool)               { return e.info, e.info != nil }
func (e *entry) Spec() (ukspec.Parameters, bool) { return e.spec, e.exec != nil }

func (e *entry) Child(target ...string) (Entry, bool) {
	child, ok := e, true

	for _, name := range target {
		if child, ok = child.children[name]; !ok {
			break
		}
	}

	return child, ok
}

func (e *entry) Children() map[string]Entry {
	children := make(map[string]Entry, len(e.children))

	for name, child := range e.children {
		children[name] = child
	}

	return children
}
