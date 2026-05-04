package ukhelp

import (
	"cmp"
	"io"
	"slices"
	"text/template"

	"github.com/oligarch316/ukase/ukcore/uktree"
)

// =============================================================================
// Encoder
// =============================================================================

// type Encoder interface {
// 	Encode(ctx Context, target ...string) error
// }

// -----------------------------------------------------------------------------
// Template
// -----------------------------------------------------------------------------

type TemplateEncoder struct {
	Template *template.Template
	Out      io.Writer
}

type Data struct {
	Command     EntryCommand
	Subcommands []DataItem
	Flags       []DataItem
	Arguments   []DataItem
}

type DataItem struct {
	Label string
	Info  string
}

func (te TemplateEncoder) Encode(ctx Context, target ...string) error {
	data, err := te.build(ctx, target)
	if err != nil {
		return err
	}

	return te.Template.Execute(te.Out, data)
}

func (te TemplateEncoder) build(ctx Context, target []string) (Data, error) {
	entry, err := uktree.Read(ctx, target...)
	if err != nil {
		return Data{}, err
	}

	children, err := uktree.ReadChildren(ctx, target...)
	if err != nil {
		return Data{}, err
	}

	data := Data{
		Command:     entry.Command,
		Subcommands: te.buildSubcommands(children),
		Flags:       te.buildFlags(entry),
		Arguments:   te.buildArguments(entry),
	}

	compare := func(a, b DataItem) int { return cmp.Compare(a.Label, b.Label) }
	slices.SortFunc(data.Subcommands, compare)
	slices.SortFunc(data.Flags, compare)
	slices.SortFunc(data.Arguments, compare)

	return data, nil
}

func (TemplateEncoder) buildSubcommands(children map[string]Entry) (items []DataItem) {
	for name, entry := range children {
		item := DataItem{Label: name, Info: entry.Command.Info}
		items = append(items, item)
	}

	return
}

func (TemplateEncoder) buildFlags(entry Entry) (items []DataItem) {
	for _, field := range entry.Fields {
		if field.FlagLabel != "" {
			item := DataItem{Label: field.FlagLabel, Info: field.Info}
			items = append(items, item)
		}
	}

	return
}

func (TemplateEncoder) buildArguments(entry Entry) (items []DataItem) {
	for _, field := range entry.Fields {
		if field.ArgumentLabel != "" {
			item := DataItem{Label: field.ArgumentLabel, Info: field.Info}
			items = append(items, item)
		}
	}

	return
}
