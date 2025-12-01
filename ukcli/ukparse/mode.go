package ukparse

import (
	"fmt"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukcontext"
	"github.com/oligarch316/ukase/ukcore/ukinput"
	"github.com/oligarch316/ukase/ukcore/uktree"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// Mode
// =============================================================================

type Mode interface {
	Load(ukcli.Context) ukinput.Context
}

type ModeFunc func(ukcli.Context) ukinput.Context

func (mf ModeFunc) Load(ctx ukcli.Context) ukinput.Context { return mf(ctx) }

// -----------------------------------------------------------------------------
// Strict
// -----------------------------------------------------------------------------

type ModeStrict func(ukcli.Values) ukinput.Value

func (ms ModeStrict) Load(ctx ukcli.Context) ukinput.Context {
	return ukcontext.Derive(ctx, ms.derive)
}

func (ms ModeStrict) derive(source ukcli.Entry) (ukinput.Entry, error) {
	val := ms(source.Values)

	sinks, err := ukvalue.Collect(source.Spec, val)
	if err != nil {
		return nil, err
	}

	return mergeEntries(sinks)
}

// -----------------------------------------------------------------------------
// Loose
// -----------------------------------------------------------------------------

type ModeLoose func(ukcli.Values) ukinput.Value

func (ml ModeLoose) Load(ctx ukcli.Context) ukinput.Context {
	inputCtx := ModeStrict(ml).Load(ctx)
	mergedTree, err := mergeChildren(inputCtx)
	if err != nil {
		deriveErr := func(ukinput.Entry) (ukinput.Entry, error) { return nil, err }
		return ukcontext.Derive(inputCtx, deriveErr)
	}

	return ukcontext.New(inputCtx, mergedTree)
}

// -----------------------------------------------------------------------------
// Merge
// -----------------------------------------------------------------------------

func mergeEntries(entries []ukinput.Entry) (ukinput.Entry, error) {
	merged := make(ukinput.Entry)

	for _, entry := range entries {
		for name, kind := range entry {
			if orig, exists := merged[name]; exists && !orig.Equal(kind) {
				return nil, fmt.Errorf("[TODO mergeEntries] flag kind conflict on '%s'", name)
			}

			merged[name] = kind
		}
	}

	return merged, nil
}

func mergeChildren(source ukcore.Tree[ukinput.Entry]) (*uktree.Node[ukinput.Entry], error) {
	sourceEntry, err := source.Load()
	if err != nil {
		return nil, err
	}

	sink := uktree.New[ukinput.Entry]()
	list := []ukinput.Entry{sourceEntry}

	for childName, childSource := range source.List() {
		childSink, err := mergeChildren(childSource)
		if err != nil {
			return nil, err
		}

		sink.Children[childName] = childSink
		list = append(list, childSink.Entry)
	}

	sink.Entry, err = mergeEntries(list)
	return sink, err
}
