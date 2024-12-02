package ukhelp

import (
	"cmp"
	"reflect"
	"slices"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore/ukspec"
	"github.com/oligarch316/ukase/ukmeta"
)

// =============================================================================
// Parameters
// =============================================================================

// TODO: Document
type Parameters struct{ Default, Info ukmeta.Parameter[any] }

// TODO: Document
func ParameterDefault(ctx ukcli.Context, t reflect.Type) ukmeta.Parameter[any] {
	ptrVal := reflect.New(t)

	if err := ctx.Initialize(ptrVal.Interface()); err != nil {
		return ukmeta.ParameterConst[any](nil, err)
	}

	val := ptrVal.Elem()

	return func(index []int) (any, error) {
		fieldVal, err := val.FieldByIndexErr(index)
		if err != nil {
			return nil, err
		}

		return fieldVal.Interface(), nil
	}
}

// TODO: Document
func ParameterInfo(ctx ukcli.Context, t reflect.Type) ukmeta.Parameter[any] {
	// TODO:
	// Implement pulling from struct tag `ukinfo:"..."` as a no-generate option
	return ukmeta.ParameterConst[any](nil, nil)
}

// =============================================================================
// Provider
// =============================================================================

// TODO: Document
type Provider func(ukcli.Context, ukspec.Parameters) (Parameters, error)

// TODO: Document
func (p Provider) Collect(ctx ukmeta.Context, entry Entry) (Data, error) {
	var data Data

	data.Command = collectCommand(entry)
	data.Subcommands = collectSubcommands(entry)

	spec, ok := entry.Spec()
	if !ok {
		return data, nil
	}

	params, err := p(ctx, spec)
	if err != nil {
		return data, err
	}

	if data.Arguments, err = collectArguments(params, spec.Arguments); err != nil {
		return data, err
	}

	if data.Flags, err = collectFlags(params, spec.Flags); err != nil {
		return data, err
	}

	return data, nil
}

// -----------------------------------------------------------------------------
// Command
// -----------------------------------------------------------------------------

func collectCommand(entry Entry) DataCommand {
	data := DataCommand{Program: entry.Program, Target: entry.Target}
	data.Info, _ = entry.Info()
	_, data.Exec = entry.Exec()
	return data
}

// -----------------------------------------------------------------------------
// Subcommand
// -----------------------------------------------------------------------------

func collectSubcommands(entry Entry) []DataSubcommand {
	var list []DataSubcommand

	for name, child := range entry.Children() {
		info, _ := child.Info()
		item := DataSubcommand{Info: info, Name: name}
		list = append(list, item)
	}

	// Sort by lexicographic order of subcommand name
	compare := func(a, b DataSubcommand) int { return cmp.Compare(a.Name, b.Name) }
	slices.SortFunc(list, compare)

	return list
}

// -----------------------------------------------------------------------------
// Argument
// -----------------------------------------------------------------------------

func collectArguments(params Parameters, specs []ukspec.Argument) ([]DataArgument, error) {
	var list []DataArgument

	for _, spec := range specs {
		item, err := collectArgument(params, spec)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	// Sort by position
	compare := func(a, b DataArgument) int {
		switch {
		case a.Position.Low == nil && b.Position.Low == nil:
			return 0
		case a.Position.Low == nil:
			return -1
		case b.Position.Low == nil:
			return +1
		default:
			return cmp.Compare(*a.Position.Low, *b.Position.Low)
		}
	}

	slices.SortFunc(list, compare)

	return list, nil
}

func collectArgument(params Parameters, spec ukspec.Argument) (DataArgument, error) {
	info, err := params.Info(spec.FieldIndex)
	if err != nil {
		return DataArgument{}, err
	}

	def, err := params.Default(spec.FieldIndex)
	if err != nil {
		return DataArgument{}, err
	}

	data := DataArgument{Info: info, Default: def, Position: spec.Position}
	return data, nil
}

// -----------------------------------------------------------------------------
// Flag
// -----------------------------------------------------------------------------

func collectFlags(params Parameters, specs []ukspec.Flag) ([]DataFlag, error) {
	var list []DataFlag

	for _, spec := range specs {
		item, err := collectFlag(params, spec)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	// Sort by lexicographic order of 1st flag name
	compare := func(a, b DataFlag) int { return cmp.Compare(a.Names[0], b.Names[0]) }
	slices.SortFunc(list, compare)

	return list, nil
}

func collectFlag(params Parameters, spec ukspec.Flag) (DataFlag, error) {
	info, err := params.Info(spec.FieldIndex)
	if err != nil {
		return DataFlag{}, err
	}

	def, err := params.Default(spec.FieldIndex)
	if err != nil {
		return DataFlag{}, err
	}

	data := DataFlag{Info: info, Default: def, Names: spec.Names}
	return data, nil
}
