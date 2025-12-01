package ukcore

import (
	"context"
	"iter"
	"reflect"
)

// =============================================================================
// Context
// =============================================================================

type Context[T any] interface {
	context.Context
	Tree[T]
}

// =============================================================================
// Tree
// =============================================================================

type Tree[T any] interface {
	Load() (T, error)
	List() iter.Seq2[string, Tree[T]]
	Lookup(target ...string) (Tree[T], bool)
}

// =============================================================================
// Spec
// =============================================================================

type Spec interface {
	Source() reflect.Type
	Fields() iter.Seq[SpecField]
	Field(index ...int) (SpecField, bool)
}

type SpecField struct {
	Source reflect.StructField
	Index  []int
}

type SpecValue[T any] interface {
	Load(spec Spec, index []int) (T, error)
}

// =============================================================================
// Input
// =============================================================================

type Input struct {
	Program   string
	Target    []string
	Arguments []InputArgument
	Flags     []InputFlag
}

type InputArgument struct {
	Position int
	Value    string
}

type InputFlag struct {
	Name  string
	Value string
}
