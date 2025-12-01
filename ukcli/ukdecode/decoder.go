package ukdecode

import (
	"slices"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukparams"
	"github.com/oligarch316/ukase/ukcore/uktree"
)

// =============================================================================
// Decoder
// =============================================================================

type Decoder struct{ config Config }

func New(opts ...Option) Decoder {
	config := newConfig(opts)
	return Decoder{config: config}
}

func (d Decoder) Decode(ctx ukcli.Context, input ukcore.Input, v any) error {
	entry, err := uktree.Read(ctx, input.Target...)
	if err != nil {
		return err
	}

	sink, err := ukparams.NewSink(v)
	if err != nil {
		return err
	}

	sources := ukparams.SourceList{}
	for _, item := range d.config {
		source := item(entry.Values)
		sources = append(sources, source)
	}

	return sources.Decode(sink, input)
}

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

type Option interface{ UkaseApplyDecode(*Config) }

type Config []Source

func newConfig(opts []Option) Config {
	config := slices.Clone(cfgDefault)
	for _, opt := range opts {
		opt.UkaseApplyDecode(&config)
	}
	return config
}

var cfgDefault = Config{
	SourceInit,
	SourceEnv,
	SourceInput,
}
