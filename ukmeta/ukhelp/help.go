package ukhelp

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukexec"
	"github.com/oligarch316/ukase/ukcore/ukspec"
	"github.com/oligarch316/ukase/ukmeta"
)

// =============================================================================
// Pipeline
// =============================================================================

// TODO: Document
type Pipeline[Params any] struct {
	Locate  func(ukmeta.Context, Params) (Entry, error)
	Collect func(ukmeta.Context, Entry, Params) (Data, error)
	Render  func(ukmeta.Context, Data, Params) error
}

// TODO: Document
func (p Pipeline[Params]) Execute(ctx ukmeta.Context, input ukcore.Input) error {
	var params Params

	if err := ctx.Initialize(&params); err != nil {
		return err
	}

	if err := ctx.Decode(input, &params); err != nil {
		return err
	}

	entry, err := p.Locate(ctx, params)
	if err != nil {
		return err
	}

	data, err := p.Collect(ctx, entry, params)
	if err != nil {
		return err
	}

	return p.Render(ctx, data, params)
}

// -----------------------------------------------------------------------------
// Locate
// -----------------------------------------------------------------------------

// TODO: Document
type Entry struct {
	ukexec.Entry
	Program string
	Target  []string
}

// TODO: Document
func Locate(ctx ukmeta.Context, program string) (Entry, error) {
	reference := ctx.MetaReference()
	entry, ok := ctx.LoadEntry(reference...)
	if !ok {
		return Entry{}, fmt.Errorf("[TODO ukhelp.Locate] no entry found for '%s'", reference)
	}

	return Entry{Entry: entry, Program: program, Target: reference}, nil
}

// -----------------------------------------------------------------------------
// Collect
// -----------------------------------------------------------------------------

// TODO: Document
type Data struct {
	Command     DataCommand
	Subcommands []DataSubcommand
	Arguments   []DataArgument
	Flags       []DataFlag
}

// TODO: Document
type DataCommand struct {
	Info    any
	Exec    bool
	Program string
	Target  []string
}

// TODO: Document
type DataSubcommand struct {
	Info any
	Name string
}

// TODO: Document
type DataArgument struct {
	Info     any
	Default  any
	Position ukspec.ArgumentPosition
}

// TODO: Document
type DataFlag struct {
	Info    any
	Default any
	Names   ukspec.FlagNames
}

// TODO: Document
func Collect(ctx ukmeta.Context, entry Entry) (Data, error) {
	return Provider(Provide).Collect(ctx, entry)
}

// TODO: Document
func Provide(ctx ukcli.Context, spec ukspec.Parameters) (Parameters, error) {
	def := ParameterDefault(ctx, spec.Type)
	info := ParameterInfo(ctx, spec.Type)
	return Parameters{Default: def, Info: info}, nil
}

// -----------------------------------------------------------------------------
// Render
// -----------------------------------------------------------------------------

// TODO: Document
func Render(ctx ukmeta.Context, data Data, w io.Writer) error {
	return EncodeText(w).Render(ctx, data)
}

// TODO: Document
func EncodeJSON(w io.Writer, opts ...func(*json.Encoder)) Encoder {
	enc := json.NewEncoder(w)
	for _, opt := range opts {
		opt(enc)
	}

	return func(data Data) error { return enc.Encode(data) }
}

// TODO: Document
func EncodeText(w io.Writer) Encoder {
	return func(data Data) error { return renderTemplate(w, data) }
}
