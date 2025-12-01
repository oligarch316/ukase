package ukspec

import (
	"fmt"
	"iter"
	"maps"
	"reflect"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Params
// =============================================================================

type Params struct {
	source reflect.Type
	fields map[string]ukcore.SpecField
}

func (p Params) Source() reflect.Type {
	return p.source
}

func (p Params) Fields() iter.Seq[ukcore.SpecField] {
	return maps.Values(p.fields)
}

func (p Params) Field(index ...int) (ukcore.SpecField, bool) {
	key := FormatIndex(index)
	val, ok := p.fields[key]
	return val, ok
}

// -----------------------------------------------------------------------------
// Load
// -----------------------------------------------------------------------------

func load(t reflect.Type) (Params, error) {
	type Scope struct {
		Type  reflect.Type
		Index []int
		Seen  map[reflect.Type]struct{}
	}

	paramsSource, valid := normalize(t)
	if !valid {
		return Params{}, fmt.Errorf("[TODO load] not a struct or struct pointer: %s", t)
	}

	params := Params{
		source: paramsSource,
		fields: make(map[string]ukcore.SpecField),
	}

	seedScope := Scope{
		Type: paramsSource,
		Seen: map[reflect.Type]struct{}{paramsSource: {}},
	}

	for queue := []Scope{seedScope}; len(queue) > 0; queue = queue[1:] {
		scope := queue[0]

		for i := range scope.Type.NumField() {
			fieldSource := scope.Type.Field(i)
			if ignored(fieldSource) {
				continue
			}

			fieldIndex := append(scope.Index, i)
			fieldKey := FormatIndex(fieldIndex)
			params.fields[fieldKey] = ukcore.SpecField{Source: fieldSource, Index: fieldIndex}

			nextSource, validStruct := normalize(fieldSource.Type)
			if !validStruct {
				continue
			}

			if _, seen := scope.Seen[nextSource]; seen {
				return Params{}, fmt.Errorf("[TODO load] cycle detected: %s", nextSource)
			}

			nextSeen := maps.Clone(scope.Seen)
			nextSeen[nextSource] = struct{}{}

			nextScope := Scope{Type: nextSource, Index: fieldIndex, Seen: nextSeen}
			queue = append(queue, nextScope)
		}
	}

	return params, nil
}

func normalize(t reflect.Type) (reflect.Type, bool) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t, t.Kind() == reflect.Struct
}

func ignored(sf reflect.StructField) bool {
	const ignoreTagKey = "ukase"
	const ignoreTagVal = "-"

	return !sf.IsExported() || sf.Tag.Get(ignoreTagKey) == ignoreTagVal
}
