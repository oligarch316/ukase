package ukhelp

import (
	_ "embed"

	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/oligarch316/ukase/ukmeta"
)

// =============================================================================
// Info
// =============================================================================

// TODO: Document
type Info struct{ Long, Short string }

// =============================================================================
// Encoder
// =============================================================================

// TODO: Document
type Encoder func(Data) error

// TODO: Document
func (e Encoder) Render(_ ukmeta.Context, data Data) error { return e(data) }

// =============================================================================
// Template
// =============================================================================

//go:embed text.tmpl
var templateText string

var templateFuncs = template.FuncMap{
	"hasCommand":     hasCommand,
	"hasSubcommands": hasSubcommands,
	"hasArguments":   hasArguments,
	"hasFlags":       hasFlags,

	"infoCommand":    infoCommand,
	"infoSubcommand": infoSubcommand,
	"infoArgument":   infoArgument,
	"infoFlag":       infoFlag,

	"labelCommand":    labelCommand,
	"labelSubcommand": labelSubcommand,
	"labelArgument":   labelArgument,
	"labelFlag":       labelFlag,

	"measureSubcommands": measureSubcommands,
	"measureArguments":   measureArguments,
	"measureFlags":       measureFlags,
}

func renderTemplate(w io.Writer, data Data) error {
	t, err := template.New("ukhelp").Funcs(templateFuncs).Parse(templateText)
	if err != nil {
		return err
	}

	return t.Execute(w, data)
}

// -----------------------------------------------------------------------------
// Has
// -----------------------------------------------------------------------------

func hasCommand(d Data) bool     { return d.Command.Exec }
func hasSubcommands(d Data) bool { return len(d.Subcommands) != 0 }
func hasArguments(d Data) bool   { return len(d.Arguments) != 0 }
func hasFlags(d Data) bool       { return len(d.Flags) != 0 }

// -----------------------------------------------------------------------------
// Info
// -----------------------------------------------------------------------------

func infoCommand(d DataCommand) (Info, error)       { return info(d.Info) }
func infoSubcommand(d DataSubcommand) (Info, error) { return info(d.Info) }
func infoArgument(d DataArgument) (Info, error)     { return info(d.Info) }
func infoFlag(d DataFlag) (Info, error)             { return info(d.Info) }

func info(info any) (Info, error) {
	type custom interface{ UkaseHelpInfo() Info }

	switch v := info.(type) {
	case nil:
		return Info{}, nil
	case string:
		return Info{Long: v, Short: v}, nil
	case Info:
		return v, nil
	case custom:
		return v.UkaseHelpInfo(), nil
	default:
		return Info{}, fmt.Errorf("[TODO ukhelp.renderInfo] invalid info type %T", info)
	}
}

// -----------------------------------------------------------------------------
// Label
// -----------------------------------------------------------------------------

func labelCommand(d DataCommand) string {
	segs := append([]string{d.Program}, d.Target...)
	return strings.Join(segs, " ")
}

func labelSubcommand(d DataSubcommand) string { return d.Name }

func labelArgument(d DataArgument) string {
	// TODO: More human friendly display than half open range???
	return strings.Replace(d.Position.String(), ":", "...", 1)
}

func labelFlag(d DataFlag) string {
	var items []string

	for _, name := range d.Names {
		switch len(name) {
		case 0:
		case 1:
			items = append(items, "-"+name)
		default:
			items = append(items, "--"+name)
		}
	}

	return strings.Join(items, ", ")
}

// -----------------------------------------------------------------------------
// Measure
// -----------------------------------------------------------------------------

func measureSubcommands(d Data) int { return measure(d.Subcommands, labelSubcommand) }
func measureArguments(d Data) int   { return measure(d.Arguments, labelArgument) }
func measureFlags(d Data) int       { return measure(d.Flags, labelFlag) }

func measure[L ~[]T, T any](list L, label func(T) string) (max int) {
	for _, item := range list {
		if candidate := len(label(item)); candidate > max {
			max = candidate
		}
	}
	return
}
