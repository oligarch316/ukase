package ukparams

import (
	"errors"
	"os"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

var _ Source = SourceEnv{}

// =============================================================================
// Source Env
// =============================================================================

type SourceEnv struct {
	EnvNames ukcore.SpecValue[[]string]
}

// TODO:
// Make `SourceEnv` a `SpecValue` itself?
// Have it do the os.LookupEnv within load
// Then we can just ukvaluex.Collect on it...? (not actually so simple)
//
// Look into other sources doing similar, etc...

func (se SourceEnv) Decode(sink Sink, _ ukcore.Input) error {
	for field := range sink.Spec.Fields() {
		if err := se.decodeField(sink, field.Index); err != nil {
			return err
		}
	}

	return nil
}

func (se SourceEnv) decodeField(sink Sink, index []int) error {
	keys, err := se.EnvNames.Load(sink.Spec, index)
	switch {
	case errors.Is(err, ukvalue.ErrNotSpecified):
		return nil
	case err != nil:
		return err
	}

	var srcs []string
	for _, key := range keys {
		if src, ok := os.LookupEnv(key); ok {
			srcs = append(srcs, src)
		}
	}

	return sink.DecodeField(index, srcs...)
}
