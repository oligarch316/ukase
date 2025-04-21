package ukgen

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"golang.org/x/tools/go/packages"
)

// =============================================================================
// Type
// =============================================================================

type Type interface{ ukgenType() }

func TypeFor[T any](opts ...func(*packages.Config)) (Type, error) {
	t := reflect.TypeFor[T]()
	return LoadType(t, opts...)
}

func TypeOf(i any, opts ...func(*packages.Config)) (Type, error) {
	t := reflect.TypeOf(i)
	return LoadType(t, opts...)
}

func LoadType(t reflect.Type, opts ...func(*packages.Config)) (Type, error) {
	kind, name, pkgPath := t.Kind(), t.Name(), t.PkgPath()

	// No package path ⇒ Native type
	if pkgPath == "" {
		// TODO: Deal with `any`
		// reflect.TypeFor[string]()
		// > kind == "string", name == "string", path == ""
		//
		// reflect.TypeFor[any]()
		// > kind == "inteface", name == "", path == ""

		return TypeNative(name), nil
	}

	// Basic kind ⇒ Direct named type
	if _, ok := typeBasicSet[kind]; ok {
		pkg, err := LoadPackage(pkgPath, opts...)
		namedType := TypeNamed{Name: name, Package: pkg}
		return namedType, err

		// // TODO: Load package using opts...
		// return TypeNamed{Name: name, Package: pkgPath}, nil
	}

	// Indirect kind ⇒ Recurse on elements
	switch kind {
	case reflect.Map:
		keyType, keyErr := LoadType(t.Key(), opts...)
		valType, valErr := LoadType(t.Elem(), opts...)
		return TypeMap{Key: keyType, Value: valType}, errors.Join(keyErr, valErr)
	case reflect.Pointer:
		elemType, err := LoadType(t.Elem(), opts...)
		return TypePointer{Element: elemType}, err
	case reflect.Slice:
		elemType, err := LoadType(t.Elem(), opts...)
		return TypeSlice{Element: elemType}, err
	}

	// Unsupported kind
	return nil, fmt.Errorf("[TODO NewType] unsupported kind: %s", kind)
}

var typeBasicSet = map[reflect.Kind]struct{}{
	reflect.Bool:       {},
	reflect.Int:        {},
	reflect.Int8:       {},
	reflect.Int16:      {},
	reflect.Int32:      {},
	reflect.Int64:      {},
	reflect.Uint:       {},
	reflect.Uint8:      {},
	reflect.Uint16:     {},
	reflect.Uint32:     {},
	reflect.Uint64:     {},
	reflect.Float32:    {},
	reflect.Float64:    {},
	reflect.Complex64:  {},
	reflect.Complex128: {},
	reflect.String:     {},
	reflect.Struct:     {},
	reflect.Interface:  {},
}

type TypeNamed struct {
	Name    string
	Package Package
}

type TypeNative string
type TypeMap struct{ Key, Value Type }
type TypePointer struct{ Element Type }
type TypeSlice struct{ Element Type }

func (TypeNative) ukgenType()  {}
func (TypeNamed) ukgenType()   {}
func (TypeMap) ukgenType()     {}
func (TypePointer) ukgenType() {}
func (TypeSlice) ukgenType()   {}

// =============================================================================
// Dependencies
// =============================================================================

type depID any
type depIDLocal struct{}
type depIDBlank struct{}
type depIDNamed string

type Dependencies struct {
	pathToID  map[string]depID
	nameToIdx map[string]int
}

func NewDependencies(pkgPath string) *Dependencies {
	deps := &Dependencies{
		pathToID:  make(map[string]depID),
		nameToIdx: make(map[string]int),
	}

	deps.pathToID[pkgPath] = depIDLocal{}
	return deps
}

func (d *Dependencies) Imports() map[string]string {
	imports := make(map[string]string)

	for pkgPath, id := range d.pathToID {
		switch idVal := id.(type) {
		case depIDBlank:
			imports[pkgPath] = "_"
		case depIDNamed:
			imports[pkgPath] = string(idVal)
		}
	}

	return imports
}

// -----------------------------------------------------------------------------
// Dependencies› Add
// -----------------------------------------------------------------------------

func (d *Dependencies) Add(name, pkgPath string) {
	if _, exists := d.pathToID[pkgPath]; exists {
		return
	}

	if name == "_" {
		d.pathToID[pkgPath] = depIDBlank{}
		return
	}

	idx := d.nameToIdx[name]
	pkgName := name + strconv.Itoa(idx)

	d.nameToIdx[name] = idx + 1
	d.pathToID[pkgPath] = depIDNamed(pkgName)
}

func (d *Dependencies) AddType(t Type) {
	switch tVal := t.(type) {
	case TypeNamed:
		d.Add(tVal.Package.Name, tVal.Package.Path)
	case TypeMap:
		d.AddType(tVal.Key)
		d.AddType(tVal.Value)
	case TypePointer:
		d.AddType(tVal.Element)
	case TypeSlice:
		d.AddType(tVal.Element)
	}
}

// -----------------------------------------------------------------------------
// Dependencies› Render
// -----------------------------------------------------------------------------

func (d *Dependencies) RenderType(t Type) (string, error) {
	switch tVal := t.(type) {
	case TypeNative:
		return string(tVal), nil
	case TypeNamed:
		return d.renderTypeNamed(tVal)
	case TypeMap:
		return d.renderTypeMap(tVal)
	case TypePointer:
		return d.renderTypePointer(tVal)
	case TypeSlice:
		return d.renderTypeSlice(tVal)
	}

	return "", fmt.Errorf("[TODO Dependencies.RenderType] invalid type '%v'", t)
}

func (d *Dependencies) renderTypeNamed(t TypeNamed) (string, error) {
	id, ok := d.pathToID[t.Package.Path]
	if !ok {
		return "", fmt.Errorf("[TODO Dependencies.renderTypeNamed] package path '%s' not found", t.Package.Path)
	}

	switch idVal := id.(type) {
	case depIDLocal:
		return t.Name, nil
	case depIDNamed:
		return fmt.Sprintf("%s.%s", idVal, t.Name), nil
	}

	return "", fmt.Errorf("[TODO Dependencies.renderTypeNamed] invalid id '%v'", id)
}

func (d *Dependencies) renderTypeMap(t TypeMap) (string, error) {
	key, keyErr := d.RenderType(t.Key)
	val, valErr := d.RenderType(t.Value)

	str := fmt.Sprintf("map[%s]%s", key, val)
	err := errors.Join(keyErr, valErr)
	return str, err
}

func (d *Dependencies) renderTypePointer(t TypePointer) (string, error) {
	elem, err := d.RenderType(t.Element)
	return fmt.Sprintf("*%s", elem), err
}

func (d *Dependencies) renderTypeSlice(t TypeSlice) (string, error) {
	elem, err := d.RenderType(t.Element)
	return fmt.Sprintf("[]%s", elem), err
}
