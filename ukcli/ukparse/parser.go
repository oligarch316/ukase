package ukparse

import (
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukinput"
)

// =============================================================================
// Parser
// =============================================================================

type Parser struct{ config Config }

func New(opts ...Option) Parser {
	config := newConfig(opts)
	return Parser{config: config}
}

func (p Parser) Format(flagName string) string {
	return p.config.Style.Format(flagName)
}

func (p Parser) Parse(ctx ukcli.Context, args []string) (ukcore.Input, error) {
	inputCtx := p.config.Mode.Load(ctx)
	return p.config.Style.Parse(inputCtx, args)
}

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

type Option interface{ UkaseApplyParse(*Config) }

type Config struct {
	Mode  Mode
	Style ukinput.Style
}

func newConfig(opts []Option) Config {
	config := cfgDefault
	for _, opt := range opts {
		opt.UkaseApplyParse(&config)
	}
	return config
}

var cfgDefault = Config{
	Mode:  ModeStrict(cfgModeValue),
	Style: ukinput.StyleUkase{Delimiter: "--"},
}

func cfgModeValue(vals ukcli.Values) ukinput.Value {
	return ukinput.Value{FlagNames: vals.FlagNames, FlagKind: vals.FlagKind}
}
