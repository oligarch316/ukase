package ukgen

import (
	"maps"
	"reflect"
	"slices"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore/ukspec"
	"github.com/oligarch316/ukase/ukmeta"
)

// =============================================================================
// Extract
// =============================================================================

var _ Extractor = Extract(nil)

type Extractor interface {
	Extract(ukcli.Context) (Source, error)
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type Extract func(ukcli.Context) (Source, error)

func (e Extract) Extract(ctx ukcli.Context) (Source, error) { return e(ctx) }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

func ExtractSpecs(specs *[]ukspec.Parameters) Extract {
	return func(ctx ukcli.Context) (Source, error) { return extractSpecs(ctx, *specs) }
}

func extractSpecs(ctx ukcli.Context, specs []ukspec.Parameters) (Source, error) {
	specMap := make(map[reflect.Type]ukspec.Parameters)
	inlines := make([]ukspec.Inline, 0)

	for _, spec := range specs {
		if _, seen := specMap[spec.Type]; seen {
			continue
		}

		specMap[spec.Type] = spec
		inlines = append(inlines, spec.Inlines...)
	}

	for len(inlines) > 0 {
		inline := inlines[0]
		inlines = inlines[1:]

		if _, seen := specMap[inline.FieldType]; seen {
			continue
		}

		spec, err := ctx.LoadSpec(inline.FieldType)
		if err != nil {
			return nil, err
		}

		specMap[spec.Type] = spec
		inlines = append(inlines, spec.Inlines...)
	}

	specVals := maps.Values(specMap)
	return slices.Collect(specVals), nil
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

func ExtractRecord[E ukmeta.Event](record *ukmeta.Record[E]) Extract {
	return func(ctx ukcli.Context) (Source, error) { return extractRecord(ctx, *record) }
}

func extractRecord[E ukmeta.Event](ctx ukcli.Context, record ukmeta.Record[E]) (Source, error) {
	var specs []ukspec.Parameters

	for _, event := range record {
		eventExec, ok := any(event).(ukmeta.EventExec)
		if !ok {
			continue
		}

		specs = append(specs, eventExec.Spec)
	}

	return extractSpecs(ctx, specs)
}
