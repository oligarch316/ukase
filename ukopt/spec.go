package ukopt

import (
	"github.com/oligarch316/ukase"
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore/ukdec"
	"github.com/oligarch316/ukase/ukcore/ukinit"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// General
// =============================================================================

var (
	_ ukspec.Option = Spec(nil)
	_ ukdec.Option  = Spec(nil)
	_ ukinit.Option = Spec(nil)
	_ ukcli.Option  = Spec(nil)
	_ ukase.Option  = Spec(nil)
)

type Spec func(*ukspec.Config)

func (o Spec) UkaseApplySpec(c *ukspec.Config) { o(c) }
func (o Spec) UkaseApplyDec(c *ukdec.Config)   { c.Spec = append(c.Spec, o) }
func (o Spec) UkaseApplyInit(c *ukinit.Config) { c.Spec = append(c.Spec, o) }
func (o Spec) UkaseApplyApp(c *ukase.Config)   { c.CLI = append(c.CLI, o) }

func (o Spec) UkaseApplyCLI(c *ukcli.Config) {
	c.Spec = append(c.Spec, o)
	c.Decode = append(c.Decode, o)
	c.Init = append(c.Init, o)
}

// =============================================================================
// Specific
// =============================================================================

// TODO
