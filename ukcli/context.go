package ukcli

import (
	"cmp"
	"context"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukcontext"
	"github.com/oligarch316/ukase/ukcore/ukinput"
	"github.com/oligarch316/ukase/ukcore/ukspec"
	"github.com/oligarch316/ukase/ukcore/uktodo"
	"github.com/oligarch316/ukase/ukcore/uktree"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// Context
// =============================================================================

type Context ukcore.Context[Entry]

type Entry struct {
	Exec   Exec
	Info   string
	Spec   ukcore.Spec
	Values Values
}

func NewContext(ctx context.Context, state State) Context {
	derive := func(local StateEntry) (entry Entry, err error) {
		entry.Exec = deriveExec(state, local)
		entry.Info = deriveInfo(state, local)
		entry.Values = deriveValues(state, local)
		entry.Spec, err = deriveSpec(state, local)
		return
	}

	tree := uktree.Derive(state.Local, derive)
	return ukcontext.New(ctx, tree)
}

func deriveInfo(state State, local StateEntry) string {
	return cmp.Or(local.Info, state.Global.Info)
}

func deriveSpec(state State, local StateEntry) (ukcore.Spec, error) {
	t := cmp.Or(local.Exec.execT, state.Global.Exec.execT)
	return ukspec.New(t)
}

// -----------------------------------------------------------------------------
// Exec
// -----------------------------------------------------------------------------

type Exec interface {
	Execute(Context, ukcore.Input) error
}

func deriveExec(state State, local StateEntry) Exec {
	return cmp.Or(local.Exec.exec, state.Global.Exec.exec)
}

// -----------------------------------------------------------------------------
// Values
// -----------------------------------------------------------------------------

type Values struct {
	EnvNames  ukcore.SpecValue[[]string]
	FlagNames ukcore.SpecValue[[]string]
	MetaNames ukcore.SpecValue[[]string]
	ArgRange  ukcore.SpecValue[uktodo.ArgRange]
	FlagKind  ukcore.SpecValue[ukinput.FlagKind]
	Initial   ukcore.SpecValue[[]string]
	Info      ukcore.SpecValue[string]
}

func deriveValues(state State, local StateEntry) Values {
	return Values{
		EnvNames:  ukvalue.Or(local.Overrides.EnvNames.vMap, state.Global.Overrides.EnvNames.vMap, state.Values.EnvNames),
		FlagNames: ukvalue.Or(local.Overrides.FlagNames.vMap, state.Global.Overrides.FlagNames.vMap, state.Values.FlagNames),
		MetaNames: ukvalue.Or(local.Overrides.MetaNames.vMap, state.Global.Overrides.MetaNames.vMap, state.Values.MetaNames),
		ArgRange:  ukvalue.Or(local.Overrides.ArgRange.vMap, state.Global.Overrides.ArgRange.vMap, state.Values.ArgRange),
		FlagKind:  ukvalue.Or(local.Overrides.FlagKind.vMap, state.Global.Overrides.FlagKind.vMap, state.Values.FlagKind),
		Initial:   ukvalue.Or(local.Overrides.Initial.vMap, state.Global.Overrides.Initial.vMap, state.Values.Initial),
		Info:      ukvalue.Or(local.Overrides.Info.vMap, state.Global.Overrides.Info.vMap, state.Values.Info),
	}
}
