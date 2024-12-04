package ukcli

import (
	"context"
	"errors"
	"reflect"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukdec"
	"github.com/oligarch316/ukase/ukcore/ukexec"
	"github.com/oligarch316/ukase/ukcore/ukinit"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

var directiveNoop DirectiveFunc = func(State) error { return nil }

// =============================================================================
// Context
// =============================================================================

var _ Context = inputContext{}

type Context interface {
	context.Context

	// TODO: Document
	LoadEntry(target ...string) (ukexec.Entry, bool)

	// TODO: Document
	LoadSpec(t reflect.Type, opts ...ukspec.Option) (ukspec.Parameters, error)

	// TODO: Document
	Decode(input ukcore.Input, v any, opts ...ukdec.Option) error

	// TODO: Document
	Initialize(v any, opts ...ukspec.Option) error
}

type inputContext struct {
	context.Context
	state State
}

func newInputContext(ctx context.Context, state State) inputContext {
	return inputContext{Context: ctx, state: state}
}

func (ic inputContext) LoadEntry(target ...string) (ukexec.Entry, bool) {
	return ic.state.loadEntry(target)
}

func (ic inputContext) LoadSpec(t reflect.Type, opts ...ukspec.Option) (ukspec.Parameters, error) {
	return ic.state.LoadSpec(t, opts...)
}

func (ic inputContext) Decode(input ukcore.Input, v any, opts ...ukdec.Option) error {
	return ic.state.runDecode(input, v, opts...)
}

func (ic inputContext) Initialize(v any, opts ...ukspec.Option) error {
	return ic.state.runInit(v, opts...)
}

// =============================================================================
// Rule
// =============================================================================

func NewRule[Params any](rule func(*Params)) Directive {
	if rule == nil {
		return directiveNoop
	}

	initRule := ukinit.NewRule(rule)
	directive := func(s State) error { s.AddRule(initRule); return nil }
	return DirectiveFunc(directive)
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

	directive := func(s State) error { return s.AddInfo(i.Value, target) }
	return DirectiveFunc(directive)
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

	directive := func(s State) error { return e.add(s, target) }
	return DirectiveFunc(directive)
}

func (e Exec[Params]) add(state State, target []string) error {
	t := reflect.TypeFor[Params]()

	spec, err := state.LoadSpec(t)
	if err != nil {
		return err
	}

	exec := func(ctx context.Context, input ukcore.Input) error {
		inputCtx := newInputContext(ctx, state)
		return e(inputCtx, input)
	}

	return state.AddExec(exec, spec, target)
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
	directive := func(s State) error { return c.add(s, target) }
	return DirectiveFunc(directive)
}

func (c Command[Params]) add(state State, target []string) error {
	errExec := c.Exec.Bind(target...).UkaseDirective(state)
	errInfo := c.Info.Bind(target...).UkaseDirective(state)
	return errors.Join(errExec, errInfo)
}
