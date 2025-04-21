package ukgen

import (
	"errors"
	"fmt"

	"golang.org/x/tools/go/packages"
)

type Module struct{ Dir, Path string }

func LoadModule(opts ...func(*packages.Config)) (Module, error) {
	mode, pattern := packages.NeedModule, ""
	data, err := loadPackage(mode, pattern, opts...)
	if err != nil {
		return Module{}, err
	}

	if data.Module == nil {
		return Module{}, fmt.Errorf("[TODO LoadModule] got a nil module")
	}

	if data.Module.Error != nil {
		return Module{}, fmt.Errorf("[TODO LoadModule] module error: %s", data.Module.Error.Err)
	}

	mdl := Module{Dir: data.Module.Dir, Path: data.Module.Path}
	return mdl, nil
}

type Package struct{ Dir, Name, Path string }

func LoadPackage(pattern string, opts ...func(*packages.Config)) (Package, error) {
	mode := packages.NeedName | packages.NeedFiles
	data, err := loadPackage(mode, pattern, opts...)
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

func loadPackage(mode packages.LoadMode, pattern string, opts ...func(*packages.Config)) (*packages.Package, error) {
	config := &packages.Config{Mode: mode}
	for _, opt := range opts {
		opt(config)
	}

	list, err := packages.Load(config, pattern)
	if err != nil {
		return nil, err
	}

	if len(list) != 1 {
		return nil, fmt.Errorf("[TODO load] unexpected number of load results, count: %d, pattern: %s", len(list), pattern)
	}

	return list[0], nil
}
