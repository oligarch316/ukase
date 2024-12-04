package ukmeta

import (
	"context"
	"errors"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukexec"
)

// =============================================================================
// Context
// =============================================================================

type Context interface {
	ukcli.Context

	// TODO: Document
	MetaReference() []string
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

func (i Info) Auto(suffix ...string) ukcli.Proxy {
	auto := func(reference []string) ukcli.Directive {
		target := append(reference, suffix...)
		return i.Refer(reference...).Bind(target...)
	}

	return Auto(auto)
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
		inputCtx := inputContext{Context: ctx, reference: reference}
		return e(inputCtx, input)
	}

	return exec
}

func (e Exec[Params]) Auto(suffix ...string) ukcli.Proxy {
	auto := func(reference []string) ukcli.Directive {
		target := append(reference, suffix...)
		return e.Refer(reference...).Bind(target...)
	}

	return Auto(auto)
}

// =============================================================================
// Command
// =============================================================================

type Command[Params any] struct {
	Exec Exec[Params]
	Info Info
}

func NewCommand[Params any](handler func(context.Context, ukexec.Entry, Params) error, info any) Command[Params] {
	return Command[Params]{
		Exec: NewExec(handler),
		Info: NewInfo(info),
	}
}

func (c Command[Params]) Refer(reference ...string) ukcli.Command[Params] {
	return ukcli.Command[Params]{
		Exec: c.Exec.Refer(reference...),
		Info: c.Info.Refer(reference...),
	}
}

func (c Command[Params]) Auto(suffix ...string) ukcli.Proxy {
	auto := func(reference []string) ukcli.Directive {
		target := append(reference, suffix...)
		return c.Refer(reference...).Bind(target...)
	}

	return Auto(auto)
}
