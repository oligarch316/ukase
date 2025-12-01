package ukparams

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Sink
// =============================================================================

type Sink struct {
	Data  any
	Value reflect.Value
	Spec  ukcore.Spec
}

func NewSink(v any) (Sink, error) {
	val := reflect.ValueOf(v)
	switch {
	case val.Kind() != reflect.Pointer:
		return Sink{}, errors.New("value is not a pointer")
	case val.IsNil():
		return Sink{}, errors.New("value is a nil pointer")
	case val.Elem().Kind() != reflect.Struct:
		return Sink{}, errors.New("value does not point to a struct")
	}

	val = val.Elem()
	spec, err := ukspec.New(val.Type())

	return Sink{Data: v, Value: val, Spec: spec}, err
}

// -----------------------------------------------------------------------------
// Load
// -----------------------------------------------------------------------------

func (s Sink) LoadField(index []int) (val reflect.Value, err error) {
	defer s.recoverField(&err)
	return s.loadField(s.Value, index), err
}

func (Sink) loadField(val reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if val.Kind() == reflect.Pointer {
			if val.IsZero() {
				newVal := reflect.New(val.Type().Elem())
				val.Set(newVal)
			}

			val = val.Elem()
		}

		val = val.Field(i)
	}

	return val
}

func (Sink) recoverField(err *error) {
	if recover() != nil {
		*err = errors.New("[TODO Sink.recoverField] caught us a panic")
	}
}

// -----------------------------------------------------------------------------
// Decode
// -----------------------------------------------------------------------------

func (s Sink) DecodeField(index []int, srcs ...string) error {
	dst, err := s.LoadField(index)
	if err != nil {
		return err
	}

	for _, src := range srcs {
		if err := decodeField(dst, src); err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Decode Field
// › Dispatch entrypoint
// › Indirect and container decode logic recurses back here
// -----------------------------------------------------------------------------

func decodeField(dst reflect.Value, src string) error {
	if complete, err := decodeFieldIndirect(dst, src); complete {
		return err
	}

	if complete, err := decodeFieldCustom(dst, src); complete {
		return err
	}

	if complete, err := decodeFieldContainer(dst, src); complete {
		return err
	}

	if complete, err := decodeFieldDirect(dst, src); complete {
		return err
	}

	// TODO:
	// We should error on unsupported field kinds during spec creation
	// When/if that happens, this can become an internal error
	return fmt.Errorf("[TODO decodeField] unsupported destination kind '%s'", dst.Kind())
}

// -----------------------------------------------------------------------------
// Decode Indirect Field
// › Handles interfaces and pointers
// -----------------------------------------------------------------------------

func decodeFieldIndirect(dst reflect.Value, src string) (bool, error) {
	switch dst.Kind() {
	case reflect.Interface:
		return true, decodeFieldInterface(dst, src)
	case reflect.Pointer:
		return true, decodeFieldPointer(dst, src)
	default:
		return false, nil
	}
}

func decodeFieldInterface(dst reflect.Value, src string) error {
	// Interface already contains a concrete type+value
	// ⇒ Copy that type+value to attain "settability"
	// ⇒ Decode into this copy and set `dst` on success
	if !dst.IsZero() {
		elemOld := dst.Elem()

		elemNew := reflect.New(elemOld.Type()).Elem()
		elemNew.Set(elemOld)

		if err := decodeField(elemNew, src); err != nil {
			return err
		}

		dst.Set(elemNew)
		return nil
	}

	// Interface contains no concrete type+value
	// ⇒ Choose the simple string type as a sane default
	// ⇒ Confirm assignability and set `dst` on success
	srcVal := reflect.ValueOf(src)
	if srcVal.Type().AssignableTo(dst.Type()) {
		dst.Set(srcVal)
		return nil
	}

	// Interface won't accept a simple string
	// ⇒ A reference value is required in this case, so fail
	return errors.New("[TODO decodeFieldInterface] interface destination neither contains a non-zero value nor is string-assignable")
}

func decodeFieldPointer(dst reflect.Value, src string) error {
	if !dst.IsZero() {
		// Why bother with this?
		// ⇒ TODO: Document
		return decodeField(dst.Elem(), src)
	}

	elemType := dst.Type().Elem()

	val := reflect.New(elemType)
	if err := decodeField(val.Elem(), src); err != nil {
		return err
	}

	dst.Set(val)
	return nil
}

// -----------------------------------------------------------------------------
// Decode Custom Field
// › Handles encoding.TextUnmarshaler implementations
// -----------------------------------------------------------------------------

var typeTextUnmarshaler = reflect.TypeFor[encoding.TextUnmarshaler]()

func decodeFieldCustom(dst reflect.Value, src string) (bool, error) {
	unmarshaler, ok := loadTextUnmarshaler(dst)
	if !ok {
		return false, nil
	}

	if err := unmarshaler.UnmarshalText([]byte(src)); err != nil {
		return true, err
	}

	return true, nil
}

func loadTextUnmarshaler(val reflect.Value) (encoding.TextUnmarshaler, bool) {
	typ := val.Type()

	if typ.Implements(typeTextUnmarshaler) {
		unmarshaler, ok := val.Interface().(encoding.TextUnmarshaler)
		return unmarshaler, ok
	}

	if reflect.PointerTo(typ).Implements(typeTextUnmarshaler) {
		if !val.CanAddr() {
			// This should never be the case as `val` is expected to be a field
			// of an addressable (provided via pointer) struct.
			// Still, safety first
			return nil, false
		}

		unmarshaler, ok := val.Addr().Interface().(encoding.TextUnmarshaler)
		return unmarshaler, ok
	}

	return nil, false
}

// -----------------------------------------------------------------------------
// Decode Container Field
// › Handles slices, arrays and structs
// -----------------------------------------------------------------------------

func decodeFieldContainer(dst reflect.Value, src string) (bool, error) {
	switch dst.Kind() {
	case reflect.Slice:
		return true, decodeFieldSlice(dst, src)
	case reflect.Array:
		return true, decodeFieldArray(dst, src)
	case reflect.Struct:
		return true, decodeFieldStruct(dst, src)
	default:
		return false, nil
	}
}

func decodeFieldSlice(dst reflect.Value, src string) error {
	elemType := dst.Type().Elem()
	elemVal := reflect.New(elemType).Elem()

	if err := decodeField(elemVal, src); err != nil {
		return err
	}

	val := reflect.Append(dst, elemVal)
	dst.Set(val)
	return nil
}

func decodeFieldArray(dst reflect.Value, src string) error {
	_, _ = dst, src
	return errors.New("[TODO decodeFieldArray] not yet implemented")
}

func decodeFieldStruct(dst reflect.Value, src string) error {
	_, _ = dst, src
	return errors.New("[TODO decodeFieldStruct] not yet implemented")
}

// -----------------------------------------------------------------------------
// Decode Direct Field
// › Handles "basic" types (bool, numeric, string)
// -----------------------------------------------------------------------------

var basicDecoders = map[reflect.Kind]func(reflect.Value, string) error{
	reflect.Bool:       decodeFieldBool,
	reflect.Int:        decodeFieldInt,
	reflect.Int8:       decodeFieldInt,
	reflect.Int16:      decodeFieldInt,
	reflect.Int32:      decodeFieldInt,
	reflect.Int64:      decodeFieldInt,
	reflect.Uint:       decodeFieldUint,
	reflect.Uint8:      decodeFieldUint,
	reflect.Uint16:     decodeFieldUint,
	reflect.Uint32:     decodeFieldUint,
	reflect.Uint64:     decodeFieldUint,
	reflect.Float32:    decodeFieldFloat,
	reflect.Float64:    decodeFieldFloat,
	reflect.Complex64:  decodeFieldComplex,
	reflect.Complex128: decodeFieldComplex,
	reflect.String:     decodeFieldString,
}

func decodeFieldDirect(dst reflect.Value, src string) (bool, error) {
	kind := dst.Kind()

	if decodeFunc, ok := basicDecoders[kind]; ok {
		// TODO:
		// The "user" errors coming from here are not "pretty" enough for actual end-users
		// example > strconv.ParseBool: parsing "ipsum": invalid syntax
		return true, decodeFunc(dst, src)
	}

	return false, nil
}

func decodeFieldBool(dst reflect.Value, src string) error {
	boolVal, err := strconv.ParseBool(src)
	if err != nil {
		return err
	}

	dst.SetBool(boolVal)
	return nil
}

func decodeFieldInt(dst reflect.Value, src string) error {
	intVal, err := strconv.ParseInt(src, 10, dst.Type().Bits())
	if err != nil {
		return err
	}

	dst.SetInt(intVal)
	return nil
}

func decodeFieldUint(dst reflect.Value, src string) error {
	uintVal, err := strconv.ParseUint(src, 10, dst.Type().Bits())
	if err != nil {
		return err
	}

	dst.SetUint(uintVal)
	return nil
}

func decodeFieldFloat(dst reflect.Value, src string) error {
	floatVal, err := strconv.ParseFloat(src, dst.Type().Bits())
	if err != nil {
		return err
	}

	dst.SetFloat(floatVal)
	return nil
}

func decodeFieldComplex(dst reflect.Value, src string) error {
	complexVal, err := strconv.ParseComplex(src, dst.Type().Bits())
	if err != nil {
		return err
	}

	dst.SetComplex(complexVal)
	return nil
}

func decodeFieldString(dst reflect.Value, src string) error {
	dst.SetString(src)
	return nil
}
