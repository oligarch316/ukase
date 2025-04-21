package ukspec

import (
	"reflect"
	"strconv"
)

// =============================================================================
// Parameters
// =============================================================================

type Parameters struct {
	Type   reflect.Type
	Fields []Field
}

func ParametersFor[T any](opts ...Option) (Parameters, error) {
	t := reflect.TypeFor[T]()
	return Load(t, opts...)
}

func ParametersOf(i any, opts ...Option) (Parameters, error) {
	t := reflect.TypeOf(i)
	return Load(t, opts...)
}

// =============================================================================
// Field
// =============================================================================

type FieldKind int

const (
	KindNone FieldKind = iota
	KindArgument
	KindFlag
	KindInline
)

var kindToString = [...]string{
	KindNone:     "none",
	KindArgument: "argument",
	KindFlag:     "flag",
	KindInline:   "inline",
}

func (fk FieldKind) String() string {
	n, low, high := int(fk), -1, len(kindToString)

	if low < n && n < high {
		return kindToString[fk]
	}

	return "unknown(" + strconv.Itoa(n) + ")"
}

type Field struct {
	Index []int
	Kind  FieldKind
	Name  string
	Tag   reflect.StructTag
	Type  reflect.Type
}
