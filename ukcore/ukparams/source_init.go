package ukparams

import (
	"errors"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

var _ Source = SourceInit{}

// =============================================================================
// Source Init
// =============================================================================

type SourceInit struct {
	Initial ukcore.SpecValue[[]string]
}

func (si SourceInit) Decode(sink Sink, _ ukcore.Input) error {
	for field := range sink.Spec.Fields() {
		if err := si.decodeField(sink, field.Index); err != nil {
			return err
		}
	}

	return nil
}

func (si SourceInit) decodeField(sink Sink, index []int) error {
	switch srcs, err := si.Initial.Load(sink.Spec, index); {
	case errors.Is(err, ukvalue.ErrNotSpecified):
		return nil
	case err != nil:
		return err
	default:
		return sink.DecodeField(index, srcs...)
	}
}
