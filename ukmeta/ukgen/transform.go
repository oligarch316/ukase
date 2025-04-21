package ukgen

import (
	"reflect"

	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Transform
// =============================================================================

var _ Transformer = Transform(nil)
var _ Transformer = TransformParameters(nil)
var _ Transformer = TransformSet{}

// TODO: Document
type Transformer interface {
	Transform(Source) (Sink, error)
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type Transform func(Source) (Sink, error)

func (t Transform) Transform(s Source) (Sink, error) { return t(s) }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type TransformParameters func(ukspec.Parameters) (Parameters, error)

func (tp TransformParameters) Transform(source Source) (Sink, error) {
	var sink Sink

	for _, spec := range source {
		params, err := tp(spec)
		if err != nil {
			return nil, err
		}

		sink = append(sink, params)
	}

	return sink, nil
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type TransformType func(reflect.Type) (TypeNamed, error)

func (tt TransformType) Transform(source Source) (Sink, error) {
	tp := TransformParameters(tt.transformParameters)
	return tp.Transform(source)
}

func (tt TransformType) transformParameters(spec ukspec.Parameters) (Parameters, error) {
	// ----- TODO
	basicFieldType, err := TypeFor[string]()
	if err != nil {
		return Parameters{}, err
	}
	// -----

	paramsType, err := tt(spec.Type)
	if err != nil {
		return Parameters{}, err
	}

	params := Parameters{Type: paramsType}

	for _, argSpec := range spec.Arguments {
		if len(argSpec.FieldIndex) != 1 {
			continue
		}

		field := Field{Name: argSpec.FieldName, Type: basicFieldType}
		params.Fields = append(params.Fields, field)
	}

	for _, flagSpec := range spec.Flags {
		if len(flagSpec.FieldIndex) != 1 {
			continue
		}

		field := Field{Name: flagSpec.FieldName, Type: basicFieldType}
		params.Fields = append(params.Fields, field)
	}

	for _, inlineSpec := range spec.Inlines {
		if len(inlineSpec.FieldIndex) != 1 {
			continue
		}

		fieldType, err := tt(inlineSpec.FieldType)
		if err != nil {
			return params, err
		}

		field := Field{Name: inlineSpec.FieldName, Type: fieldType}
		params.Fields = append(params.Fields, field)
	}

	return params, nil
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type TransformPackage func(pkgPath string) (Package, error)

func (tp TransformPackage) Transform(source Source) (Sink, error) {
	tt := TransformType(tp.transformType)
	return tt.Transform(source)
}

func (tp TransformPackage) transformType(t reflect.Type) (tn TypeNamed, err error) {
	tn.Name = t.Name()
	tn.Package, err = tp(t.PkgPath())
	return
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type TransformSet struct {
	Parameters func(reflect.Type) (TypeNamed, error)
	Argument   func(reflect.Type) (Type, error)
	Flag       func(reflect.Type) (Type, error)
}

func (ts TransformSet) Transform(source Source) (Sink, error) {
	return TransformParameters(ts.transformParameters).Transform(source)
}

func (ts TransformSet) transformParameters(spec ukspec.Parameters) (Parameters, error) {
	paramsType, err := ts.Parameters(spec.Type)
	if err != nil {
		return Parameters{}, err
	}

	params := Parameters{Type: paramsType}

	for _, argSpec := range spec.Arguments {
		if len(argSpec.FieldIndex) != 1 {
			continue
		}

		fieldType, err := ts.Argument(argSpec.FieldType)
		if err != nil {
			return params, err
		}

		field := Field{Name: argSpec.FieldName, Type: fieldType}
		params.Fields = append(params.Fields, field)
	}

	for _, flagSpec := range spec.Flags {
		if len(flagSpec.FieldIndex) != 1 {
			continue
		}

		fieldType, err := ts.Flag(flagSpec.FieldType)
		if err != nil {
			return params, err
		}

		field := Field{Name: flagSpec.FieldName, Type: fieldType}
		params.Fields = append(params.Fields, field)
	}

	for _, inlineSpec := range spec.Inlines {
		if len(inlineSpec.FieldIndex) != 1 {
			continue
		}

		fieldType, err := ts.Parameters(inlineSpec.FieldType)
		if err != nil {
			return params, err
		}

		field := Field{Name: inlineSpec.FieldName, Type: fieldType}
		params.Fields = append(params.Fields, field)
	}

	return params, nil
}
