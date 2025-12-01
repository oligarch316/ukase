package ukparams

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/uktodo"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// TODO:
// This is all super messy, clean it up, especially the Unknown...Error stuff

var _ Source = SourceMeta{}
var _ Source = SourceArgs{}
var _ Source = SourceFlags{}

// =============================================================================
// Source Input
// =============================================================================

var ErrUnknownInput = errors.New("unknown input")

// -----------------------------------------------------------------------------
// Meta
// -----------------------------------------------------------------------------

type SourceMeta struct {
	MetaNames ukcore.SpecValue[[]string]
}

func (sm SourceMeta) Decode(sink Sink, input ukcore.Input) error {
	const nameProgram = "program"
	const nameTarget = "target"

	metaSet, err := sm.collect(sink)
	if err != nil {
		return err
	}

	if index, exists := metaSet[nameProgram]; exists {
		if err := sink.DecodeField(index, input.Program); err != nil {
			return err
		}
	}

	if index, exists := metaSet[nameTarget]; exists {
		if err := sink.DecodeField(index, input.Target...); err != nil {
			return err
		}
	}

	return nil
}

func (sm SourceMeta) collect(sink Sink) (map[string][]int, error) {
	metaSet := make(map[string][]int)

	for field := range sink.Spec.Fields() {
		names, err := sm.MetaNames.Load(sink.Spec, field.Index)
		switch {
		case errors.Is(err, ukvalue.ErrNotSpecified):
			continue
		case err != nil:
			return nil, err
		}

		for _, name := range names {
			if _, exists := metaSet[name]; exists {
				return nil, fmt.Errorf("[TODO SourceMeta.collect] meta name conflict on '%s'", name)
			}

			metaSet[name] = field.Index
		}
	}

	return metaSet, nil
}

// -----------------------------------------------------------------------------
// Flags
// -----------------------------------------------------------------------------

type UnknownFlagsError []ukcore.InputFlag

func (UnknownFlagsError) Is(target error) bool {
	return target == ErrUnknownInput
}

func (ufe UnknownFlagsError) Error() string {
	var names []string
	for _, flag := range ufe {
		names = append(names, flag.Name)
	}

	return fmt.Sprintf("unknown flag names: %s", strings.Join(names, " "))
}

type SourceFlags struct {
	FlagNames ukcore.SpecValue[[]string]
}

func (sf SourceFlags) Decode(sink Sink, input ukcore.Input) error {
	flagSet, err := sf.collect(sink)
	if err != nil {
		return err
	}

	var unknown []ukcore.InputFlag

	for _, flag := range input.Flags {
		index, exists := flagSet[flag.Name]
		if !exists {
			unknown = append(unknown, flag)
			continue
		}

		if err := sink.DecodeField(index, flag.Value); err != nil {
			return err
		}
	}

	if len(unknown) > 0 {
		return UnknownFlagsError(unknown)
	}

	return nil
}

func (sf SourceFlags) collect(sink Sink) (map[string][]int, error) {
	flagSet := make(map[string][]int)

	for field := range sink.Spec.Fields() {
		names, err := sf.FlagNames.Load(sink.Spec, field.Index)
		switch {
		case errors.Is(err, ukvalue.ErrNotSpecified):
			continue
		case err != nil:
			return nil, err
		}

		for _, name := range names {
			if _, exists := flagSet[name]; exists {
				return nil, fmt.Errorf("[TODO SourceFlags.collect] flag name conflict on '%s'", name)
			}

			flagSet[name] = field.Index
		}
	}

	return flagSet, nil
}

// -----------------------------------------------------------------------------
// Args
// -----------------------------------------------------------------------------

type UnknownArgsError []ukcore.InputArgument

func (UnknownArgsError) Is(target error) bool {
	return target == ErrUnknownInput
}

func (uae UnknownArgsError) Error() string {
	var positions []string
	for _, arg := range uae {
		position := strconv.Itoa(arg.Position)
		positions = append(positions, position)
	}

	return fmt.Sprintf("unknown arguments at positions: %s", strings.Join(positions, " "))
}

type SourceArgs struct {
	ArgRange ukcore.SpecValue[uktodo.ArgRange]
}

type argItem struct {
	argRange   uktodo.ArgRange
	fieldIndex []int
}

func (sa SourceArgs) Decode(sink Sink, input ukcore.Input) error {
	argSet, err := sa.collect(sink)
	if err != nil {
		return err
	}

	var unknown []ukcore.InputArgument

	for _, arg := range input.Arguments {
		index, exists := sa.lookup(argSet, arg.Position)
		if !exists {
			unknown = append(unknown, arg)
			continue
		}

		if err := sink.DecodeField(index, arg.Value); err != nil {
			return err
		}
	}

	if len(unknown) > 0 {
		return UnknownArgsError(unknown)
	}

	return nil
}

func (sa SourceArgs) collect(sink Sink) ([]argItem, error) {
	var argSet []argItem

	// ----- Load
	for field := range sink.Spec.Fields() {
		argRange, err := sa.ArgRange.Load(sink.Spec, field.Index)
		switch {
		case errors.Is(err, ukvalue.ErrNotSpecified):
			continue
		case err != nil:
			return nil, err
		}

		item := argItem{argRange: argRange, fieldIndex: field.Index}
		argSet = append(argSet, item)
	}

	if len(argSet) < 2 {
		return argSet, nil
	}

	// ----- Sort
	compare := func(a, b argItem) int { return cmp.Compare(a.argRange.Low, b.argRange.Low) }
	slices.SortFunc(argSet, compare)

	// ----- Validate
	prev := argSet[0]
	for _, cur := range argSet[1:] {
		// TODO: Clean up/formalize this

		if prev.argRange.High < 0 {
			return nil, fmt.Errorf("[TODO SourceArgs.collect] incompatible arg ranges")
		}

		if cur.argRange.Low <= prev.argRange.High {
			return nil, fmt.Errorf("[TODO SourceArgs.collect] incompatible arg ranges")
		}
	}

	return argSet, nil
}

func (sa SourceArgs) lookup(argSet []argItem, pos int) ([]int, bool) {
	match := func(item argItem) bool { return item.argRange.Contains(pos) }

	if setIdx := slices.IndexFunc(argSet, match); setIdx != -1 {
		return argSet[setIdx].fieldIndex, true
	}

	return nil, false
}
