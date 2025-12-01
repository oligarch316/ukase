package ukvalue

import (
	"fmt"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Function
// =============================================================================

type Func[T any] func(ukcore.Spec, []int) (T, error)

func (f Func[T]) Load(spec ukcore.Spec, index []int) (T, error) {
	return f(spec, index)
}

func From[T any](f func(ukcore.Spec, []int) (T, error)) Func[T] {
	return Func[T](f)
}

func FromConst[T any](v T) Func[T] {
	f := func(ukcore.Spec, []int) (T, error) { return v, nil }
	return Func[T](f)
}

func FromError[T any](err error) Func[T] {
	f := func(ukcore.Spec, []int) (T, error) { return *new(T), err }
	return Func[T](f)
}

func FromErrorf[T any](format string, a ...any) Func[T] {
	f := func(ukcore.Spec, []int) (T, error) { return *new(T), fmt.Errorf(format, a...) }
	return Func[T](f)
}
