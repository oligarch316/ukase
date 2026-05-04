package ukcommand

import (
	_ "embed"

	"io"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukcontext"
	"github.com/oligarch316/ukase/ukcore/ukhelp"
	"github.com/oligarch316/ukase/ukcore/uktodo"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// Helper
// =============================================================================

type HelpParams struct {
	Program string   `ukmeta:"program"`
	Target  []string `ukmeta:"target"`
	Query   []string `ukarg:":" ukinfo:"subcommand query"`
}

type Helper struct {
	decoder ukcli.Decoder
	config  HelpConfig
}

func NewHelper(dec ukcli.Decoder, opts ...HelpOption) Helper {
	config := newHelpConfig(opts)
	return Helper{decoder: dec, config: config}
}

func (h Helper) UkaseOperation(entry *ukcli.StateEntry) error {
	execT := reflect.TypeFor[HelpParams]()
	entry.Exec.Update(execT, h)
	return nil
}

func (h Helper) Execute(ctx ukcli.Context, input ukcore.Input) error {
	var params HelpParams

	if err := h.decoder.Decode(ctx, input, &params); err != nil {
		return err
	}

	helpCtx := h.buildContext(ctx, params)
	helpTgt := h.buildTarget(params)

	helpEnc, err := h.buildEncoder()
	if err != nil {
		return err
	}

	return helpEnc.Encode(helpCtx, helpTgt...)
}

func (h Helper) buildContext(ctx ukcli.Context, params HelpParams) ukhelp.Context {
	derive := func(target []string, source ukcli.Entry) (sink ukhelp.Entry, err error) {
		usageSegs := []string{params.Program}
		usageSegs = append(usageSegs, target...)

		sink.Command.Info = source.Info
		sink.Command.Usage = strings.Join(usageSegs, " ")
		sink.Command.Exec = source.Exec != nil // TODO

		// TODO: Is here the right place to parser.Format(...) flag names if at all?
		sinkVals := h.config.Format(source.Values)
		sink.Fields, err = ukvalue.Collect(source.Spec, sinkVals)

		return
	}

	return ukcontext.DeriveTarget(ctx, derive)
}

func (Helper) buildTarget(params HelpParams) []string {
	l := max(0, len(params.Target)-1)
	return append(params.Target[:l], params.Query...)
}

func (h Helper) buildEncoder() (ukhelp.TemplateEncoder, error) {
	t, err := template.New("ukhelp").Funcs(h.config.Funcs).Parse(h.config.Text)
	enc := ukhelp.TemplateEncoder{Template: t, Out: h.config.Out}
	return enc, err
}

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

type HelpOption interface{ UkaseApplyHelp(*HelpConfig) }

type HelpConfig struct {
	Text   string
	Funcs  template.FuncMap
	Out    io.Writer
	Format func(ukcli.Values) ukhelp.Values
}

func newHelpConfig(opts []HelpOption) HelpConfig {
	config := helpCfgDefault
	for _, opt := range opts {
		opt.UkaseApplyHelp(&config)
	}
	return config
}

var helpCfgDefault = HelpConfig{
	Text: string(helpText),
	Funcs: template.FuncMap{
		"hasInfo":        helpText.tHasInfo,
		"hasExec":        helpText.tHasExec,
		"hasSubcommands": helpText.tHasSubcommands,
		"hasFlags":       helpText.tHasFlags,
		"hasArguments":   helpText.tHasArguments,
		"hasUsage":       helpText.tHasUsage,
		"maxLabel":       helpText.tMaxLabel,
	},
	Out:    os.Stdout,
	Format: helpCfgFormat,
}

func helpCfgFormat(vals ukcli.Values) ukhelp.Values {
	deriveFlag := func(names []string) (string, error) { return strings.Join(names, ", "), nil }
	deriveArg := func(rng uktodo.ArgRange) (string, error) { return uktodo.FormatArgRange(rng), nil }

	return ukhelp.Values{
		Info:          vals.Info,
		FlagLabel:     ukvalue.Derive(vals.FlagNames, deriveFlag),
		ArgumentLabel: ukvalue.Derive(vals.ArgRange, deriveArg),
	}
}

// -----------------------------------------------------------------------------
// Template
// -----------------------------------------------------------------------------

//go:embed help.tmpl
var helpText helpTemplate

type helpTemplate string

func (helpTemplate) tHasInfo(d ukhelp.Data) bool        { return d.Command.Info != "" }
func (helpTemplate) tHasExec(d ukhelp.Data) bool        { return d.Command.Exec }
func (helpTemplate) tHasSubcommands(d ukhelp.Data) bool { return len(d.Subcommands) > 0 }
func (helpTemplate) tHasFlags(d ukhelp.Data) bool       { return len(d.Flags) > 0 }
func (helpTemplate) tHasArguments(d ukhelp.Data) bool   { return len(d.Arguments) > 0 }
func (helpTemplate) tHasUsage(d ukhelp.Data) bool       { return d.Command.Exec || len(d.Subcommands) > 0 }

func (helpTemplate) tMaxLabel(items []ukhelp.DataItem) (x int) {
	for _, item := range items {
		x = max(x, len(item.Label))
	}
	return
}
