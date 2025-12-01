package uktree

import (
	"iter"

	"github.com/oligarch316/ukase/ukcore"
)

type derived[T, U any] struct {
	source ukcore.Tree[T]
	derive func(T) (U, error)
}

func Derive[T, U any](source ukcore.Tree[T], derive func(T) (U, error)) ukcore.Tree[U] {
	return &derived[T, U]{source: source, derive: derive}
}

func (d *derived[T, U]) Load() (U, error) {
	entry, err := d.source.Load()
	if err != nil {
		return *new(U), err
	}

	return d.derive(entry)
}

func (d *derived[T, U]) List() iter.Seq2[string, ukcore.Tree[U]] {
	return func(yield func(string, ukcore.Tree[U]) bool) {
		for name, source := range d.source.List() {
			derived := &derived[T, U]{source: source, derive: d.derive}
			if !yield(name, derived) {
				return
			}
		}
	}
}

func (d *derived[T, U]) Lookup(target ...string) (ukcore.Tree[U], bool) {
	if source, exists := d.source.Lookup(target...); exists {
		derived := &derived[T, U]{source: source, derive: d.derive}
		return derived, true
	}

	return nil, false
}

type derivedTarget[T, U any] struct {
	source ukcore.Tree[T]
	target []string
	derive func([]string, T) (U, error)
}

func DeriveTarget[T, U any](source ukcore.Tree[T], derive func([]string, T) (U, error)) ukcore.Tree[U] {
	return &derivedTarget[T, U]{source: source, derive: derive}
}

func (d *derivedTarget[T, U]) Load() (U, error) {
	entry, err := d.source.Load()
	if err != nil {
		return *new(U), err
	}

	return d.derive(d.target, entry)
}

func (d *derivedTarget[T, U]) List() iter.Seq2[string, ukcore.Tree[U]] {
	return func(yield func(string, ukcore.Tree[U]) bool) {
		for name, source := range d.source.List() {
			target := append(d.target, name)
			derived := &derivedTarget[T, U]{source: source, target: target, derive: d.derive}
			if !yield(name, derived) {
				return
			}
		}
	}
}

func (d *derivedTarget[T, U]) Lookup(target ...string) (ukcore.Tree[U], bool) {
	if source, exists := d.source.Lookup(target...); exists {
		target := append(d.target, target...)
		derived := &derivedTarget[T, U]{source: source, target: target, derive: d.derive}
		return derived, true
	}

	return nil, false
}
