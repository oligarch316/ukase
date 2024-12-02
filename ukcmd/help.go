package ukcmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oligarch316/ukase/ukmeta"
	"github.com/oligarch316/ukase/ukmeta/ukhelp"
)

// =============================================================================
// Command
// =============================================================================

type HelpParams struct {
	Format    string   `ukflag:"f format"`
	Reference []string `ukarg:":"`
}

func Help(info any, opts ...HelpOption) ukmeta.Command[HelpParams] {
	config := newHelpConfig(opts)

	pipeline := ukhelp.Pipeline[HelpParams]{
		Locate:  config.locate,
		Collect: config.collect,
		Render:  config.render,
	}

	return ukmeta.Command[HelpParams]{
		Exec: pipeline.Execute,
		Info: ukmeta.NewInfo(info),
	}
}

// =============================================================================
// Config
// =============================================================================

type HelpOption interface{ UkaseApplyHelp(*HelpConfig) }

type HelpConfig struct {
	Provider ukhelp.Provider
	Encoders map[string]ukhelp.Encoder
}

func newHelpConfig(opts []HelpOption) HelpConfig {
	config := helpCfgDefault
	for _, opt := range opts {
		opt.UkaseApplyHelp(&config)
	}
	return config
}

func (HelpConfig) locate(ctx ukmeta.Context, params HelpParams) (ukhelp.Entry, error) {
	entry, err := ukhelp.Locate(ctx, "TODO-program")
	if err != nil {
		return entry, err
	}

	child, ok := entry.Child(params.Reference...)
	if !ok {
		return entry, fmt.Errorf("[TODO HelpConfig.locate] entry.Child empty")
	}

	entry.Entry = child
	entry.Target = append(entry.Target, params.Reference...)
	return entry, nil
}

func (hc HelpConfig) collect(ctx ukmeta.Context, entry ukhelp.Entry, _ HelpParams) (ukhelp.Data, error) {
	return hc.Provider.Collect(ctx, entry)
}

func (hc HelpConfig) render(ctx ukmeta.Context, data ukhelp.Data, params HelpParams) error {
	enc, ok := hc.Encoders[params.Format]
	if !ok {
		return fmt.Errorf("[TODO HelpConfig.render] invalid format '%s'", params.Format)
	}

	return enc.Render(ctx, data)
}

// =============================================================================
// Defaults
// =============================================================================

var helpCfgDefault = HelpConfig{
	Provider: ukhelp.Provide,
	Encoders: map[string]ukhelp.Encoder{
		"json":        ukhelp.EncodeJSON(os.Stdout),
		"json-pretty": ukhelp.EncodeJSON(os.Stdout, helpCfgJSONPretty),
		"text":        ukhelp.EncodeText(os.Stdout),
	},
}

func helpCfgJSONPretty(enc *json.Encoder) { enc.SetIndent("", "  ") }
