package ukcmd

import (
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukmeta"
	"github.com/oligarch316/ukase/ukmeta/ukgen"
)

// =============================================================================
// Seed
// =============================================================================

type Seed struct {
	ukmeta.Record[ukmeta.EventExec]
}

func (s *Seed) UkaseApplyComplete(config *CompleteConfig) {
	// TODO
}

func (s *Seed) UkaseApplyGen(config *ukgen.Config) {
	config.Extractor = ukgen.ExtractRecord(&s.Record)
}

// =============================================================================
// Pipeline
// =============================================================================

type Pipeline[Params, Source, Sink any] struct {
	Extract   func(ukcli.Context, Params) (Source, error)
	Transform func(ukcli.Context, Params, Source) (Sink, error)
	Render    func(ukcli.Context, Params, Sink) error
}

func (p Pipeline[Params, Source, Sink]) Execute(ctx ukcli.Context, input ukcore.Input) error {
	var params Params

	if err := ctx.Initialize(&params); err != nil {
		return err
	}

	if err := ctx.Decode(input, &params); err != nil {
		return err
	}

	source, err := p.Extract(ctx, params)
	if err != nil {
		return err
	}

	sink, err := p.Transform(ctx, params, source)
	if err != nil {
		return err
	}

	return p.Render(ctx, params, sink)
}
