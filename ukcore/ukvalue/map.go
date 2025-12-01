package ukvalue

import (
	"reflect"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Map
// =============================================================================

// TODO:
// Formalize serialization of MapKey from Source<reflect.Type> + Index<[]int>

type MapKey struct{ Package, Name, Index string }

type Map[T any] map[MapKey]T

func (m Map[T]) Load(spec ukcore.Spec, index []int) (T, error) {
	if m == nil {
		return *new(T), ErrNotSpecified
	}

	for pivot := range len(index) {
		head, tail := index[:pivot], index[pivot:]

		source, ok := m.loadSource(spec, head)
		if !ok {
			continue
		}

		val, ok := m.loadValue(source, tail)
		if !ok {
			continue
		}

		return val, nil
	}

	return *new(T), ErrNotSpecified
}

func (m Map[T]) loadSource(spec ukcore.Spec, index []int) (reflect.Type, bool) {
	if len(index) == 0 {
		return spec.Source(), true
	}

	if field, ok := spec.Field(index...); ok {
		return field.Source.Type, true
	}

	return nil, false
}

func (m Map[T]) loadValue(source reflect.Type, index []int) (T, bool) {
	key := MapKey{
		Package: source.PkgPath(),
		Name:    source.Name(),
		Index:   ukspec.FormatIndex(index),
	}

	val, ok := m[key]
	return val, ok
}
