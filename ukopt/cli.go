package ukopt

import (
	"github.com/oligarch316/ukase"
	"github.com/oligarch316/ukase/ukcli"
)

var (
	_ ukcli.Option = CLI(nil)
	_ ukase.Option = CLI(nil)
)

// =============================================================================
// General
// =============================================================================

type CLI func(*ukcli.Config)

func (o CLI) UkaseApplyCLI(c *ukcli.Config) { o(c) }
func (o CLI) UkaseApply(c *ukase.Config)    { c.CLI = append(c.CLI, o) }

// =============================================================================
// Specific
// =============================================================================

func CLIUnknownError(err error) CLI  { return func(c *ukcli.Config) { c.UnknownError = err } }
func CLIUnknownInfo(info string) CLI { return func(c *ukcli.Config) { c.UnknownInfo = info } }
