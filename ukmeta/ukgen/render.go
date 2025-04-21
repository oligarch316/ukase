package ukgen

import (
	"bytes"
	_ "embed"
	"go/format"
	"path/filepath"

	"cmp"
	"encoding/json"
	"errors"
	"io"
	"slices"
	"text/template"
)

// =============================================================================
// Render
// =============================================================================

var _ Renderer = Render(nil)
var _ Renderer = RenderCore("")
var _ Renderer = RenderMeta("")

// TODO: Document
type Renderer interface{ Render(Creator, Sink) error }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

// TODO: Document
type Render func(Creator, Sink) error

func (r Render) Render(fs Creator, sink Sink) error { return r(fs, sink) }

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type Encode func(io.Writer, Sink) error

func (e Encode) To(name string) Render {
	return func(fs Creator, sink Sink) error {
		w, err := fs.Create(name)
		if err != nil {
			return err
		}

		defer w.Close()
		return e(w, sink)
	}
}

func Encoder[T interface{ Encode(any) error }](f func(io.Writer) T, opts ...func(T)) Encode {
	return func(w io.Writer, sink Sink) error {
		enc := f(w)
		for _, opt := range opts {
			opt(enc)
		}

		encodeErr := enc.Encode(sink)
		closeErr := encoderClose(enc)
		return errors.Join(encodeErr, closeErr)
	}
}

func encoderClose(enc any) error {
	if closer, ok := enc.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func RenderJSON(name string, opts ...func(*json.Encoder)) Render {
	return Encoder(json.NewEncoder, opts...).To(name)
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

//go:embed stuff.tmpl
var renderCoreText string

type RenderCore string

func (rc RenderCore) Render(fs Creator, sink Sink) error {
	pkgMap := rc.collect(sink)

	for pkg, paramsList := range pkgMap {
		src, err := rc.execute(pkg, paramsList)
		if err != nil {
			return err
		}

		if err = rc.write(fs, pkg, src); err != nil {
			return err
		}
	}

	return nil
}

func (RenderCore) collect(sink Sink) map[Package][]Parameters {
	// TODO:
	// Are we concerned with conflicting Name/Dir info for the same PkgPath?

	pathMap := make(map[string]Package)
	pkgMap := make(map[Package][]Parameters)

	for _, params := range sink {
		pkgPath := params.Type.Package.Path

		pkg, ok := pathMap[pkgPath]
		if !ok {
			pkg = params.Type.Package
			pathMap[pkgPath] = pkg
		}

		pkgMap[pkg] = append(pkgMap[pkg], params)
	}

	return pkgMap
}

func (RenderCore) execute(pkg Package, paramsList []Parameters) ([]byte, error) {
	type Data struct {
		PackageName    string
		ImportMap      map[string]string
		ParametersList []Parameters
	}

	// Load dependencies
	deps := NewDependencies(pkg.Path)
	for _, params := range paramsList {
		for _, field := range params.Fields {
			deps.AddType(field.Type)
		}
	}

	// Sort parameters (lexicographic order of type name)
	compare := func(a, b Parameters) int { return cmp.Compare(a.Type.Name, b.Type.Name) }
	slices.SortFunc(paramsList, compare)

	// TODO: Sort each parameter's fields by name

	// Prepare template values
	data := Data{PackageName: pkg.Name, ImportMap: deps.Imports(), ParametersList: paramsList}
	funcs := template.FuncMap{"renderType": deps.RenderType}

	// Parse template
	t, err := template.New("ukgencore").Funcs(funcs).Parse(renderCoreText)
	if err != nil {
		return nil, err
	}

	// Execute template
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.Bytes(), err
}

func (rc RenderCore) write(fs Creator, pkg Package, src []byte) error {
	p, err := format.Source(src)
	if err != nil {
		return err
	}

	name := filepath.Join(pkg.Dir, string(rc))
	w, err := fs.Create(name)
	if err != nil {
		return err
	}

	defer w.Close()
	_, err = w.Write(p)
	return err
}

// -----------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

type RenderMeta string

func (RenderMeta) Render(Creator, Sink) error {
	return errors.New("[TODO RenderMeta] not yet implemented")
}
