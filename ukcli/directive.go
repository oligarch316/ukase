package ukcli

import (
	"context"
	"errors"
	"reflect"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukinit"
)

// =============================================================================
// Directive
// =============================================================================

var directiveNoop directiveFunc = func(State) error { return nil }

type Directive interface{ UkaseRegister(State) error }

func NewDirective(directive func(State) error) Directive { return directiveFunc(directive) }
func ErrDirective(err error) Directive                   { return directiveErr{err} }

type directiveFunc func(State) error
type directiveErr struct{ error }

func (df directiveFunc) UkaseRegister(state State) error { return df(state) }
func (de directiveErr) UkaseRegister(State) error        { return de.error }

// =============================================================================
// Rule
// =============================================================================

type Rule[Params any] func(*Params)

func NewRule[Params any](rule func(*Params)) Rule[Params] {
	return Rule[Params](rule)
}

func (r Rule[Params]) UkaseRegister(state State) error {
	if r != nil {
		state.AddRule(ukinit.NewRule(r))
	}

	return nil
}

// =============================================================================
// Info
// =============================================================================

type Info struct{ Value any }

func NewInfo(info any) Info {
	return Info{Value: info}
}

func (i Info) Bind(target ...string) Directive {
	if i.Value == nil {
		return directiveNoop
	}

	dir := func(state State) error { return state.AddInfo(i.Value, target...) }
	return directiveFunc(dir)
}

// =============================================================================
// Exec
// =============================================================================

type Exec[Params any] func(Context, ukcore.Input) error

func NewExec[Params any](handler func(context.Context, Params) error) Exec[Params] {
	if handler == nil {
		return nil
	}

	return func(ctx Context, input ukcore.Input) error {
		var params Params

		if err := ctx.Initialize(&params); err != nil {
			return err
		}

		if err := ctx.Decode(input, &params); err != nil {
			return err
		}

		return handler(ctx, params)
	}
}

func (e Exec[Params]) Bind(target ...string) Directive {
	if e == nil {
		return directiveNoop
	}

	dir := func(state State) error { return e.register(state, target) }
	return directiveFunc(dir)
}

func (e Exec[Params]) register(state State, target []string) error {
	t := reflect.TypeFor[Params]()

	spec, err := state.LoadSpec(t)
	if err != nil {
		return err
	}

	exec := func(ctx context.Context, input ukcore.Input) error {
		inputCtx := newInputContext(ctx, state)
		return e(inputCtx, input)
	}

	return state.AddExec(exec, spec, target...)
}

// =============================================================================
// Command
// =============================================================================

type Command[Params any] struct {
	Exec Exec[Params]
	Info Info
}

func NewCommand[Params any](handler func(context.Context, Params) error, info any) Command[Params] {
	return Command[Params]{
		Exec: NewExec(handler),
		Info: NewInfo(info),
	}
}

func (c Command[Params]) Bind(target ...string) Directive {
	dir := func(state State) error { return c.register(state, target) }
	return NewDirective(dir)
}

func (c Command[Params]) register(state State, target []string) error {
	errExec := c.Exec.Bind(target...).UkaseRegister(state)
	errInfo := c.Exec.Bind(target...).UkaseRegister(state)
	return errors.Join(errExec, errInfo)
}
