package ukopt

import (
	"github.com/oligarch316/ukase"
	"github.com/oligarch316/ukase/ukcli/ukdecode"
)

var (
	_ ukdecode.Option = Decode(nil)
	_ ukase.Option    = Decode(nil)
)

// =============================================================================
// General
// =============================================================================

type Decode func(*ukdecode.Config)

func (o Decode) UkaseApplyDecode(c *ukdecode.Config) { o(c) }
func (o Decode) UkaseApply(c *ukase.Config)          { c.Decode = append(c.Decode, o) }

// =============================================================================
// Specific
// =============================================================================

func DecodeSources(sources ...ukdecode.Source) Decode {
	return func(c *ukdecode.Config) { *c = sources }
}
