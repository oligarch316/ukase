package ukdirect

import (
	"context"
	"errors"
	"reflect"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukspec"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// Operation
// =============================================================================

type Operation interface {
	UkaseOperation(*ukcli.StateEntry) error
}

type OperationFunc func(*ukcli.StateEntry) error

func (of OperationFunc) UkaseOperation(s *ukcli.StateEntry) error { return of(s) }

// -----------------------------------------------------------------------------
// Command
// -----------------------------------------------------------------------------

type Command[Params any] struct {
	Info Info
	Exec Exec[Params]
}

func (c Command[Params]) UkaseOperation(entry *ukcli.StateEntry) error {
	infoErr := c.Info.UkaseOperation(entry)
	execErr := c.Exec.UkaseOperation(entry)
	return errors.Join(infoErr, execErr)
}

// -----------------------------------------------------------------------------
// Info
// -----------------------------------------------------------------------------

type Info string

func (i Info) UkaseOperation(entry *ukcli.StateEntry) error {
	entry.Info = string(i)
	return nil
}

// -----------------------------------------------------------------------------
// Exec
// -----------------------------------------------------------------------------

type Exec[Params any] func(ukcli.Context, ukcore.Input) error

func (e Exec[Params]) Execute(ctx ukcli.Context, input ukcore.Input) error {
	return e(ctx, input)
}

func (e Exec[Params]) UkaseOperation(entry *ukcli.StateEntry) error {
	t := reflect.TypeFor[Params]()
	entry.Exec.Update(t, e)
	return nil
}

// -----------------------------------------------------------------------------
// Handler
// -----------------------------------------------------------------------------

func Handler[Params any](dec ukcli.Decoder, f func(context.Context, Params) error) Exec[Params] {
	return func(ctx ukcli.Context, input ukcore.Input) error {
		var params Params

		if err := dec.Decode(ctx, input, &params); err != nil {
			return err
		}

		return f(ctx, params)
	}
}

// -----------------------------------------------------------------------------
// Default
// -----------------------------------------------------------------------------

type Default struct{ Params any }

func (d Default) UkaseOperation(se *ukcli.StateEntry) error {
	spec, err := ukspec.Of(d.Params)
	if err != nil {
		return err
	}

	paramsType := spec.Source()
	paramsVal := reflect.ValueOf(d.Params)

	keyPkg, keyName := paramsType.PkgPath(), paramsType.Name()

	for field := range spec.Fields() {
		fieldVal, err := paramsVal.FieldByIndexErr(field.Index)
		if err != nil {
			return err
		}

		if fieldVal.IsZero() {
			continue
		}

		keyIndex := ukspec.FormatIndex(field.Index)
		key := ukvalue.MapKey{Package: keyPkg, Name: keyName, Index: keyIndex}
		val := fieldVal.Interface()

		se.Overrides.Initial.Update(key, val)
	}

	return nil
}
