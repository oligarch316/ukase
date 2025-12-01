package ukopt

import (
	"github.com/oligarch316/ukase"
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcli/ukparse"
	"github.com/oligarch316/ukase/ukcore/ukinput"
)

var (
	_ ukparse.Option = Parse(nil)
	_ ukase.Option   = Parse(nil)
)

// =============================================================================
// General
// =============================================================================

type Parse func(*ukparse.Config)

func (o Parse) UkaseApplyParse(c *ukparse.Config) { o(c) }
func (o Parse) UkaseApply(c *ukase.Config)        { c.Parse = append(c.Parse, o) }

// =============================================================================
// Specific
// =============================================================================

var ParseModeStrict = ParseMode(ukparse.ModeStrict(parseValue))
var ParseModeLoose = ParseMode(ukparse.ModeLoose(parseValue))

func ParseMode(mode ukparse.Mode) Parse {
	return func(c *ukparse.Config) { c.Mode = mode }
}

func parseValue(vals ukcli.Values) ukinput.Value {
	return ukinput.Value{FlagNames: vals.FlagNames, FlagKind: vals.FlagKind}
}

var ParseStyleGo = ParseStyle(ukinput.StyleGo{})
var ParseStyleUkase = ParseStyle(ukinput.StyleUkase{Delimiter: "--"})

func ParseStyle(style ukinput.Style) Parse {
	return func(c *ukparse.Config) { c.Style = style }
}
