package ukspec

import (
	"log/slog"
	"reflect"
)

// =============================================================================
// Config
// =============================================================================

type Option interface{ UkaseApplySpec(*Config) }

type Config struct {
	// TODO: Document
	Log *slog.Logger

	// TODO: Document
	FieldKind func(reflect.StructField) (FieldKind, error)
}

func newConfig(opts []Option) Config {
	config := cfgDefault
	for _, opt := range opts {
		opt.UkaseApplySpec(&config)
	}
	return config
}

var cfgDefault = Config{
	FieldKind: myTagMap.Blah,
}

// ------------
// TODO

var myTagMap = tagMapThing{
	KindArgument: "ukarg",
	KindFlag:     "ukflag",
	KindInline:   "ukinline",
}

type tagMapThing map[FieldKind]string

func (tmt tagMapThing) Blah(sf reflect.StructField) (FieldKind, error) {
	for kind, tagKey := range tmt {
		if _, ok := sf.Tag.Lookup(tagKey); ok {
			return kind, nil
		}
	}

	return KindNone, nil
}
