package ukcmd

import (
	"errors"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Command
// =============================================================================

type CompleteParams struct {
	// TODO
}

func Complete(info any, opts ...CompleteOption) ukcli.Command[CompleteParams] {
	_ = newCompleteConfig(opts)

	exec := func(_ ukcli.Context, _ ukcore.Input) error {
		return errors.New("[TODO ukcmd.Complete] not yet implemented")
	}

	return ukcli.Command[CompleteParams]{
		Exec: exec,
		Info: ukcli.NewInfo(info),
	}
}

// =============================================================================
// Config
// =============================================================================

type CompleteOption interface{ UkaseApplyComplete(*CompleteConfig) }

type CompleteConfig struct {
	// TODO
}

func newCompleteConfig(opts []CompleteOption) CompleteConfig {
	config := completeCfgDefault
	for _, opt := range opts {
		opt.UkaseApplyComplete(&config)
	}
	return config
}

// =============================================================================
// Defaults
// =============================================================================

var completeCfgDefault = CompleteConfig{
	// TODO
}
