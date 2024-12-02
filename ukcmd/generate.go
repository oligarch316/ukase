package ukcmd

import (
	"errors"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Command
// =============================================================================

type GenerateParams struct {
	// TODO
}

func Generate(info any, opts ...GenerateOption) ukcli.Command[GenerateParams] {
	_ = newGenerateConfig(opts)

	exec := func(_ ukcli.Context, _ ukcore.Input) error {
		return errors.New("[TODO ukcmd.Generate] not yet implemented")
	}

	return ukcli.Command[GenerateParams]{
		Exec: exec,
		Info: ukcli.NewInfo(info),
	}
}

// =============================================================================
// Config
// =============================================================================

type GenerateOption interface{ UkaseApplyGenerate(*GenerateConfig) }

type GenerateConfig struct {
	// TODO
}

func newGenerateConfig(opts []GenerateOption) GenerateConfig {
	config := generateCfgDefault
	for _, opt := range opts {
		opt.UkaseApplyGenerate(&config)
	}
	return config
}

// =============================================================================
// Defaults
// =============================================================================

var generateCfgDefault = GenerateConfig{
	// TODO
}
