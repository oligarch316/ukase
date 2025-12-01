package ukcli

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukinput"
	"github.com/oligarch316/ukase/ukcore/uktodo"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// Runtime
// =============================================================================

type Runtime struct{ directives []Directive }

func New(opts ...Option) *Runtime {
	config := newConfig(opts)
	directives := []Directive{config}
	return &Runtime{directives: directives}
}

func (r *Runtime) Add(directives ...Directive) {
	r.directives = append(r.directives, directives...)
}

func (r *Runtime) Build(ctx context.Context) (Context, error) {
	state, err := NewState(r.directives...)
	return NewContext(ctx, state), err
}

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

type Option interface{ UkaseApplyCLI(*Config) }

type Config struct {
	BaseValues   Values
	UnknownError error
	UnknownInfo  string
}

func (c Config) UkaseDirective(state *State) error {
	state.Values = c.BaseValues
	state.Global.Info = c.UnknownInfo
	state.Global.Exec.exec = execError{c.UnknownError}
	state.Global.Exec.execT = reflect.TypeFor[struct{}]()
	return nil
}

func newConfig(opts []Option) Config {
	config := cfgDefault
	for _, opt := range opts {
		opt.UkaseApplyCLI(&config)
	}
	return config
}

var cfgDefault = Config{
	BaseValues:   cfgBaseValues,
	UnknownError: errors.New("unknown target"),
	UnknownInfo:  "",
}

var cfgBaseValues = Values{
	EnvNames:  ukvalue.DeriveTag("ukenv", cfgTagFields),
	FlagNames: ukvalue.DeriveTag("ukflag", cfgTagFields),
	MetaNames: ukvalue.DeriveTag("ukmeta", cfgTagFields),
	ArgRange:  ukvalue.DeriveTag("ukarg", uktodo.ParseArgRange),
	FlagKind:  ukvalue.DeriveField(cfgFlagKind),
	Initial:   ukvalue.DeriveTag("ukinit", cfgTagFields),
	Info:      ukvalue.DeriveTag("ukinfo", cfgTagString),
}

func cfgTagString(v string) (string, error) {
	return strings.TrimSpace(v), nil
}

func cfgTagFields(v string) ([]string, error) {
	v = strings.TrimSpace(v)
	return strings.Fields(v), nil
}

func cfgFlagKind(field ukcore.SpecField) (ukinput.FlagKind, error) {
	fieldType := field.Source.Type

	for fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	if fieldType.Kind() == reflect.Bool {
		return ukinput.KindBoolean, nil
	}

	return ukinput.KindBasic, nil
}

type execError struct{ error }

func (ee execError) Execute(Context, ukcore.Input) error { return ee.error }
