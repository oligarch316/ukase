package ukvalue

import (
	"errors"
	"fmt"
	"iter"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Value
// =============================================================================

// TODO: Do we need 2 sentinel errors here?
// > One for "Def not set" and another for "I personally don't know"
// > Otherwise, how can an end user override a value to "unset" something
// > Ex: "Remove this field as a flag"

var ErrNotSpecified = errors.New("not specified")

// -----------------------------------------------------------------------------
// Reify
// -----------------------------------------------------------------------------

func All[T any](spec ukcore.Spec, val ukcore.SpecValue[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for field := range spec.Fields() {
			v, err := val.Load(spec, field.Index)
			if errors.Is(err, ErrNotSpecified) {
				continue
			}

			if !yield(v, err) {
				return
			}
		}
	}
}

func Collect[T any](spec ukcore.Spec, val ukcore.SpecValue[T]) ([]T, error) {
	var vs []T
	var errs []error

	for field := range spec.Fields() {
		v, err := val.Load(spec, field.Index)
		if errors.Is(err, ErrNotSpecified) {
			continue
		}

		vs, errs = append(vs, v), append(errs, err)
	}

	return vs, errors.Join(errs...)
}

// -----------------------------------------------------------------------------
// Derive
// -----------------------------------------------------------------------------

func Derive[T, U any](source ukcore.SpecValue[T], derive func(T) (U, error)) ukcore.SpecValue[U] {
	f := func(spec ukcore.Spec, index []int) (U, error) {
		sourceVal, err := source.Load(spec, index)
		if err != nil {
			return *new(U), err
		}

		return derive(sourceVal)
	}

	return From(f)
}

func DeriveField[T any](derive func(ukcore.SpecField) (T, error)) ukcore.SpecValue[T] {
	f := func(spec ukcore.Spec, index []int) (T, error) {
		if field, ok := spec.Field(index...); ok {
			return derive(field)
		}

		return *new(T), fmt.Errorf("[TODO FromField] invalid index '%v' for spec type '%s'", index, spec.Source())

	}

	return From(f)
}

func DeriveTag[T any](key string, derive func(string) (T, error)) ukcore.SpecValue[T] {
	f := func(field ukcore.SpecField) (T, error) {
		if tagVal, ok := field.Source.Tag.Lookup(key); ok {
			return derive(tagVal)
		}

		return *new(T), ErrNotSpecified
	}

	return DeriveField(f)
}
