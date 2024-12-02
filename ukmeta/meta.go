package ukmeta

import (
	"context"
	"errors"
	"slices"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukexec"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Parameter
// =============================================================================

type Parameter[T any] func(index []int) (T, error)

func ParameterConst[T any](val T, err error) Parameter[T] {
	return func([]int) (T, error) { return val, err }
}

// =============================================================================
// Context
// =============================================================================

type Context interface {
	ukcli.Context

	// TODO: Document
	MetaReference() []string
}

func NewContext(ctx ukcli.Context, reference []string) Context {
	return inputContext{Context: ctx, reference: reference}
}

type inputContext struct {
	ukcli.Context
	reference []string
}

func (ic inputContext) MetaReference() []string { return ic.reference }

// =============================================================================
// Info
// =============================================================================

type Info func(reference []string) any

func NewInfo(info any) Info { return func([]string) any { return info } }

func (i Info) Refer(reference ...string) ukcli.Info {
	return ukcli.NewInfo(i(reference))
}

func (i Info) Auto(suffix ...string) func(ukcli.State) ukcli.State {
	bind := func(reference []string) ukcli.Directive {
		target := append(reference, suffix...)
		return i.Refer(reference...).Bind(target...)
	}

	return Auto(bind)
}

// =============================================================================
// Exec
// =============================================================================

type Exec[Params any] func(Context, ukcore.Input) error

func NewExec[Params any](handler func(context.Context, ukexec.Entry, Params) error) Exec[Params] {
	return func(ctx Context, input ukcore.Input) error {
		var params Params

		if err := ctx.Initialize(&params); err != nil {
			return err
		}

		if err := ctx.Decode(input, &params); err != nil {
			return err
		}

		reference := ctx.MetaReference()
		entry, ok := ctx.LoadEntry(reference...)
		if !ok {
			return errors.New("TODO")
		}

		return handler(ctx, entry, params)
	}
}

func (e Exec[Params]) Refer(reference ...string) ukcli.Exec[Params] {
	exec := func(ctx ukcli.Context, input ukcore.Input) error {
		inputCtx := NewContext(ctx, reference)
		return e(inputCtx, input)
	}

	return exec
}

func (e Exec[Params]) Auto(suffix ...string) func(ukcli.State) ukcli.State {
	bind := func(reference []string) ukcli.Directive {
		target := append(reference, suffix...)
		return e.Refer(reference...).Bind(target...)
	}

	return Auto(bind)
}

// =============================================================================
// Command
// =============================================================================

type Command[Params any] struct {
	Exec Exec[Params]
	Info Info
}

func (c Command[Params]) Refer(reference ...string) ukcli.Command[Params] {
	return ukcli.Command[Params]{
		Exec: c.Exec.Refer(reference...),
		Info: c.Info.Refer(reference...),
	}
}

func (c Command[Params]) Auto(suffix ...string) func(ukcli.State) ukcli.State {
	bind := func(reference []string) ukcli.Directive {
		target := append(reference, suffix...)
		return c.Refer(reference...).Bind(target...)
	}

	return Auto(bind)
}

// =============================================================================
// Auto
// =============================================================================

func Auto(bind func(reference []string) ukcli.Directive) func(ukcli.State) ukcli.State {
	return func(state ukcli.State) ukcli.State {
		return &autoStateX{State: state, bind: bind}
	}
}

type autoMemo map[string]autoMemo

type autoStateX struct {
	ukcli.State

	bind func(reference []string) ukcli.Directive
	memo autoMemo
}

func (as *autoStateX) AddExec(exec ukcore.Exec, spec ukspec.Parameters, target ...string) error {
	for _, ref := range as.sift(target) {
		directive := as.bind(ref)
		directive.UkaseRegister(as.State)
	}

	return as.State.AddExec(exec, spec, target...)
}

// Ensure each sub-path of the given target is marked as visited.
// Return a list of those that have not previously been visited.
func (as *autoStateX) sift(target []string) [][]string {
	var paths [][]string

	if as.memo == nil {
		// Root (empty) target not yet visited

		// ⇒ Initialize memo to mark as visited
		as.memo = make(autoMemo)

		// ⇒ Add empty path (nil) to the result list
		paths = [][]string{nil}
	}

	for cur, i := as.memo, 0; i < len(target); i++ {
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
