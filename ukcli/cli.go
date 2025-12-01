package ukcli

import "github.com/oligarch316/ukase/ukcore"

// =============================================================================
// CLI
// =============================================================================

type Directive interface {
	UkaseDirective(*State) error
}

type Decoder interface {
	Decode(Context, ukcore.Input, any) error
}

type Parser interface {
	Format(string) string
	Parse(Context, []string) (ukcore.Input, error)
}
