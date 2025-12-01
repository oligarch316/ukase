package ukinput

import (
	"errors"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Unix Style
// =============================================================================

type StyleUnix struct{}

// -----------------------------------------------------------------------------
// Format
// -----------------------------------------------------------------------------

func (StyleUnix) Format(flagName string) string {
	// TODO: Need to check rune length instead
	switch len(flagName) {
	case 0:
		return flagName
	case 1:
		return "-" + flagName
	default:
		return "--" + flagName
	}
}

// -----------------------------------------------------------------------------
// Parse
// -----------------------------------------------------------------------------

func (sg StyleUnix) Parse(ctx Context, args []string) (ukcore.Input, error) {
	return ukcore.Input{}, errors.New("[TODO StyleUnix.Parse] not yet implemented")
}
