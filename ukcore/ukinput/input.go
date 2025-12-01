package ukinput

import (
	"slices"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
//
// =============================================================================

type Context ukcore.Context[Entry]

type Entry map[string]FlagKind

type Value struct {
	FlagNames ukcore.SpecValue[[]string]
	FlagKind  ukcore.SpecValue[FlagKind]
}

func (v Value) Load(spec ukcore.Spec, index []int) (Entry, error) {
	valKind := ukvalue.OrErrorf(v.FlagKind, "[TODO Value.Load] missing flag kind")

	names, err := v.FlagNames.Load(spec, index)
	if err != nil {
		return nil, err
	}

	kind, err := valKind.Load(spec, index)
	if err != nil {
		return nil, err
	}

	entry := make(Entry)
	for _, name := range names {
		entry[name] = kind
	}

	return entry, nil
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

var (
	KindBasic   = NewFlagKind("")
	KindBoolean = NewFlagKind("true", "True", "TRUE", "false", "False", "FALSE")
)

type FlagKind struct {
	assumed string
	allowed []string
}

func NewFlagKind(assumed string, allowed ...string) FlagKind {
	slices.Sort(allowed)

	return FlagKind{
		assumed: assumed,
		allowed: append([]string{assumed}, allowed...),
	}
}

func (fk FlagKind) Allowed(s string) bool { return fk.Required() || slices.Contains(fk.allowed, s) }
func (fk FlagKind) Assumed() string       { return fk.assumed }
func (fk FlagKind) Equal(k FlagKind) bool { return slices.Equal(fk.allowed, k.allowed) }
func (fk FlagKind) Required() bool        { return fk.assumed == "" }
