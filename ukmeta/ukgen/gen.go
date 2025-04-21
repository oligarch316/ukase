package ukgen

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/oligarch316/ukase/internal/ilog"
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

func NewExec[Params any](handler func(context.Context, Params) ([]Option, error)) ukcli.Exec[Params] {
	return func(ctx ukcli.Context, input ukcore.Input) error {
		var params Params

		if err := ctx.Initialize(&params); err != nil {
			return err
		}

		if err := ctx.Decode(input, &params); err != nil {
			return err
		}

		opts, err := handler(ctx, params)
		if err != nil {
			return err
		}

		return Generate(ctx, opts...)
	}
}

// =============================================================================
//
// =============================================================================

func Generate(ctx ukcli.Context, opts ...Option) error {
	config := newConfig(opts)

	source, err := config.Extractor.Extract(ctx)
	if err != nil {
		return err
	}

	sink, err := config.Transformer.Transform(source)
	if err != nil {
		return err
	}

	return config.Renderer.Render(config.Creator, sink)
}

type Option interface{ UkaseApplyGen(*Config) }

type Config struct {
	// TODO: Document
	Log *slog.Logger

	// TODO: Document
	Creator Creator

	// TODO: Document
	Extractor Extractor

	// TODO: Document
	Transformer Transformer

	// TODO: Document
	Renderer Renderer
}

func newConfig(opts []Option) Config {
	config := cfgDefault
	for _, opt := range opts {
		opt.UkaseApplyGen(&config)
	}
	return config
}

var cfgDefault = Config{
	Log:         ilog.Discard,
	Creator:     Create(cfgCreate),
	Extractor:   Extract(cfgExtract),
	Transformer: Transform(cfgTransform),
	Renderer:    Render(cfgRender),
}

func cfgCreate(string) (io.WriteCloser, error) { return os.Stdout, nil }
func cfgExtract(ukcli.Context) (Source, error) { return nil, errors.New("no extractor specified") }
func cfgTransform(Source) (Sink, error)        { return nil, errors.New("no transformer specified") }
func cfgRender(Creator, Sink) error            { return errors.New("no renderer specified") }

// =============================================================================
//
// =============================================================================

type Source []ukspec.Parameters
type Sink []Parameters

// =============================================================================
//
// =============================================================================

type Parameters struct {
	Type   TypeNamed
	Fields []Field
}

type Field struct {
	Name string
	Type Type
	Tags map[string]string

	// TODO?
	// SourceIndex int
}
