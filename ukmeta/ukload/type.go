package ukload

import (
	"errors"
	"reflect"
)

// =============================================================================
// Type
// =============================================================================

type Type interface{ ukloadType() }

func TypeFor[T any](opts ...Option) (Type, error) {
	t := reflect.TypeFor[T]()
	return NewType(t, opts...)
}

func TypeOf(i any, opts ...Option) (Type, error) {
	t := reflect.TypeOf(i)
	return NewType(t, opts...)
}

func NewType(t reflect.Type, opts ...Option) (Type, error) {
	kind, name, pkgPath := t.Kind(), t.Name(), t.PkgPath()
	_, _, _ = kind, name, pkgPath

	// TODO
	return nil, errors.New("[TODO NewType] not yet implemented")
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type TypeNamed struct {
	Name    string
	Package Package
}

type TypeNative string
type TypeMap struct{ Key, Value Type }
type TypePointer struct{ Element Type }
type TypeSlice struct{ Element Type }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

var _ Type = TypeNamed{}
var _ Type = TypeNative("")
var _ Type = TypeMap{}
var _ Type = TypePointer{}
var _ Type = TypeSlice{}

func (TypeNative) ukloadType()  {}
func (TypeNamed) ukloadType()   {}
func (TypeMap) ukloadType()     {}
func (TypePointer) ukloadType() {}
func (TypeSlice) ukloadType()   {}
