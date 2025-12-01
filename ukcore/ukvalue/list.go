package ukvalue

import (
	"errors"
	"slices"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// List
// =============================================================================

type List[T any] []ukcore.SpecValue[T]

func (l List[T]) Load(spec ukcore.Spec, index []int) (T, error) {
	for items := l; len(items) > 0; items = items[1:] {
		if val, err := items[0].Load(spec, index); !errors.Is(err, ErrNotSpecified) {
			return val, err
		}
	}

	return *new(T), ErrNotSpecified
}

func Or[T any](vals ...ukcore.SpecValue[T]) List[T] {
	isNil := func(val ukcore.SpecValue[T]) bool { return val == nil }
	vals = slices.DeleteFunc(vals, isNil)
	return List[T](vals)
}

func OrConst[T any](val ukcore.SpecValue[T], v T) List[T] {
	return Or(val, FromConst(v))
}

func OrError[T any](val ukcore.SpecValue[T], err error) List[T] {
	return Or(val, FromError[T](err))
}

func OrErrorf[T any](val ukcore.SpecValue[T], format string, a ...any) List[T] {
	return Or(val, FromErrorf[T](format, a...))
}
