package ukspec

import (
	"fmt"
	"maps"
	"reflect"
)

// =============================================================================
//
// =============================================================================

func Load(t reflect.Type, opts ...Option) (Parameters, error) {
	seed, err := newLoadScope(t)
	if err != nil {
		return Parameters{}, err
	}

	config := newConfig(opts)
	params := Parameters{Type: seed.Type}

	for queue := []loadScope{seed}; len(queue) > 0; queue = queue[1:] {
		scope := queue[0]

		for i := range scope.Type.NumField() {
			fieldInfo := scope.Type.Field(i)
			if !fieldInfo.IsExported() {
				continue
			}

			fieldKind, err := config.FieldKind(fieldInfo)
			if err != nil {
				return params, err
			}

			if fieldKind == KindNone {
				continue
			}

			field := Field{
				Index: append(scope.Index, fieldInfo.Index...),
				Kind:  fieldKind,
				Name:  fieldInfo.Name,
				Tag:   fieldInfo.Tag,
				Type:  fieldInfo.Type,
			}

			params.Fields = append(params.Fields, field)

			if fieldKind != KindInline {
				continue
			}

			childScope, err := scope.child(field.Type, field.Index)
			if err != nil {
				return params, err
			}

			queue = append(queue, childScope)
		}
	}

	return params, nil
}

type loadScope struct {
	Type  reflect.Type
	Index []int
	seen  map[reflect.Type]struct{}
}

func newLoadScope(t reflect.Type) (loadScope, error) {
	empty := make(map[reflect.Type]struct{})
	scope := loadScope{seen: empty}
	return scope.child(t, nil)
}

func (ls loadScope) child(t reflect.Type, index []int) (loadScope, error) {
	childType, err := ls.normalize(t)
	if err != nil {
		return loadScope{}, err
	}

	childSeen, err := ls.check(childType)
	if err != nil {
		return loadScope{}, err
	}

	child := loadScope{Type: childType, Index: index, seen: childSeen}
	return child, nil
}

func (loadScope) normalize(t reflect.Type) (reflect.Type, error) {
	normalT := t

	for normalT.Kind() == reflect.Pointer {
		normalT = normalT.Elem()
	}

	if normalT.Kind() != reflect.Struct {
		return nil, fmt.Errorf("[TODO] not a struct or struct pointer: %s", t)
	}

	return normalT, nil
}

func (ls loadScope) check(t reflect.Type) (map[reflect.Type]struct{}, error) {
	if _, exists := ls.seen[t]; exists {
		return nil, fmt.Errorf("[TODO] cycle stuffz: %s", t)
	}

	clone := maps.Clone(ls.seen)
	clone[t] = struct{}{}
	return clone, nil
}
