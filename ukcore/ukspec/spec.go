package ukspec

import (
	"reflect"
	"strconv"
	"strings"
)

// =============================================================================
// Spec
// =============================================================================

func For[T any]() (Params, error) {
	t := reflect.TypeFor[T]()
	return load(t)
}

func Of(i any) (Params, error) {
	t := reflect.TypeOf(i)
	return load(t)
}

func New(t reflect.Type) (Params, error) {
	return load(t)
}

// -----------------------------------------------------------------------------
// Index Conversion
// -----------------------------------------------------------------------------

const indexSep = "."

func FormatIndex(index []int) string {
	segs := make([]string, len(index))
	for i, item := range index {
		segs[i] = strconv.Itoa(item)
	}
	return strings.Join(segs, indexSep)
}

func ParseIndex(s string) ([]int, error) {
	segs := strings.Split(s, indexSep)
	index := make([]int, len(segs))
	for i, seg := range segs {
		item, err := strconv.Atoi(seg)
		if err != nil {
			return nil, err
		}

		index[i] = item
	}
	return index, nil
}
