package ukparams

import "github.com/oligarch316/ukase/ukcore"

var _ Source = SourceFunc(nil)
var _ Source = SourceList{}

// =============================================================================
// Source
// =============================================================================

type Source interface {
	Decode(Sink, ukcore.Input) error
}

type SourceFunc func(ukcore.Input, Sink) error

func (sf SourceFunc) Decode(sink Sink, input ukcore.Input) error {
	return sf(input, sink)
}

type SourceList []Source

func (sl SourceList) Decode(sink Sink, input ukcore.Input) error {
	for _, source := range sl {
		if err := source.Decode(sink, input); err != nil {
			return err
		}
	}

	return nil
}
