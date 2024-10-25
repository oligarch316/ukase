package ukgen

import (
	"reflect"

	"github.com/oligarch316/ukase/internal/ispec"
	"github.com/oligarch316/ukase/ukcli/ukinfo"
	"github.com/oligarch316/ukase/ukmeta"
)

const tagKeyIndex = "ukidx"

var (
	pkgPathUkgen  = reflect.TypeFor[Config]().PkgPath()
	pkgPathUkinfo = reflect.TypeFor[ukinfo.Any]().PkgPath()
	pkgPathUkmeta = reflect.TypeFor[ukmeta.Input]().PkgPath()
)

// =============================================================================
// Data
// =============================================================================

type coreData struct {
	Names    coreNamesData
	Packages corePackagesData
	Types    coreTypesData
}

type coreNamesData struct {
	Package            string
	EncoderConstructor string
	EncoderDefault     string
	EncoderType        string
	TagKeyIndex        string
	TagKeyInline       string
}

type corePackagesData struct {
	Ukgen  string
	Ukinfo string
	Ukmeta string
}

type coreTypesData struct {
	ArgumentInfo typeData
	FlagInfo     typeData
}

// =============================================================================
// Generate
// =============================================================================

func (g *Generator) generateCore() (coreData, error) {
	names := g.generateCoreNames()

	packages, err := g.generateCorePackages()
	if err != nil {
		return coreData{}, err
	}

	types, err := g.generateCoreTypes()
	if err != nil {
		return coreData{}, err
	}

	data := coreData{Names: names, Packages: packages, Types: types}
	return data, nil
}

func (g *Generator) generateCoreNames() coreNamesData {
	return coreNamesData{
		Package:            g.config.Names.Package,
		EncoderConstructor: g.config.Names.EncoderConstructor,
		EncoderDefault:     g.config.Names.EncoderDefault,
		EncoderType:        g.config.Names.EncoderType,
		TagKeyIndex:        tagKeyIndex,
		TagKeyInline:       ispec.TagKeyInline,
	}
}

func (g *Generator) generateCorePackages() (corePackagesData, error) {
	pkgUkgen, err := g.loadImportName(pkgPathUkgen)
	if err != nil {
		return corePackagesData{}, err
	}

	pkgUkinfo, err := g.loadImportName(pkgPathUkinfo)
	if err != nil {
		return corePackagesData{}, err
	}

	pkgUkmeta, err := g.loadImportName(pkgPathUkmeta)
	if err != nil {
		return corePackagesData{}, err
	}

	data := corePackagesData{
		Ukgen:  pkgUkgen,
		Ukinfo: pkgUkinfo,
		Ukmeta: pkgUkmeta,
	}

	return data, nil
}

func (g *Generator) generateCoreTypes() (coreTypesData, error) {
	argumentInfo, err := g.loadImport(g.config.Types.ArgumentInfo)
	if err != nil {
		return coreTypesData{}, err
	}

	flagInfo, err := g.loadImport(g.config.Types.FlagInfo)
	if err != nil {
		return coreTypesData{}, err
	}

	data := coreTypesData{ArgumentInfo: argumentInfo, FlagInfo: flagInfo}
	return data, nil
}

// =============================================================================
// Utility
// =============================================================================

func (g Generator) reservedNames() map[string]struct{} {
	return map[string]struct{}{
		g.config.Names.EncoderConstructor: {},
		g.config.Names.EncoderDefault:     {},
		g.config.Names.EncoderType:        {},
	}
}
