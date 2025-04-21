package ukload

import (
	"errors"
	"fmt"

	"golang.org/x/tools/go/packages"
)

// =============================================================================
// Load
// =============================================================================

type Option interface{ UkaseApplyLoad(*packages.Config) }

func load(mode packages.LoadMode, pattern string, opts []Option) (*packages.Package, error) {
	config := &packages.Config{Mode: mode}
	for _, opt := range opts {
		opt.UkaseApplyLoad(config)
	}

	list, err := packages.Load(config, pattern)
	if err != nil {
		return nil, err
	}

	if len(list) != 1 {
		return nil, fmt.Errorf("[TODO loadPackage] unexpected number of load results, count: %d, pattern: %s", len(list), pattern)
	}

	return list[0], nil
}

// -----------------------------------------------------------------------------
// Module
// -----------------------------------------------------------------------------

type Module struct{ Dir, Path string }

func NewModule(opts ...Option) (Module, error) {
	mode := packages.NeedModule
	data, err := load(mode, "", opts)
	if err != nil {
		return Module{}, err
	}

	if data.Module == nil {
		return Module{}, fmt.Errorf("[TODO NewModule] got a nil module")
	}

	if data.Module.Error != nil {
		return Module{}, fmt.Errorf("[TODO NewModule] module error: %s", data.Module.Error.Err)
	}

	mod := Module{Dir: data.Module.Dir, Path: data.Module.Path}
	return mod, nil
}

// -----------------------------------------------------------------------------
// Package
// -----------------------------------------------------------------------------

type Package struct{ Dir, Name, Path string }

func NewPackage(pattern string, opts ...Option) (Package, error) {
	mode := packages.NeedName | packages.NeedFiles
	data, err := load(mode, pattern, opts)
	if err != nil {
		return Package{}, err
	}

	pkgErrs := make([]error, len(data.Errors))
	for i, pkgErr := range data.Errors {
		pkgErrs[i] = pkgErr
	}

	pkg := Package{Dir: data.Dir, Name: data.Name, Path: data.PkgPath}
	return pkg, errors.Join(pkgErrs...)
}
