package ukcmd

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukmeta/ukgen"
	"golang.org/x/tools/go/packages"
)

// =============================================================================
// Command
// =============================================================================

type GenerateParams struct {
	Out string `ukflag:"o out"`

	PatternCore string `ukflag:"c core"`
	PatternMeta string `ukflag:"m meta"`

	LoadBuild []string `ukflag:"load-build"`
	LoadDir   string   `ukflag:"load-dir"`
	LoadEnv   []string `ukflag:"load-env"`
}

func Generate(info any, opts ...ukgen.Option) ukcli.Command[GenerateParams] {
	handle := func(ctx context.Context, params GenerateParams) ([]ukgen.Option, error) {
		loadOpt := func(config *packages.Config) {
			config.BuildFlags = params.LoadBuild
			config.Context = ctx
			config.Dir = params.LoadDir
			config.Env = params.LoadEnv
		}

		modRoot, err := ukgen.LoadModule(loadOpt)
		if err != nil {
			return nil, err
		}

		pkgCore, err := ukgen.LoadPackage(params.PatternCore, loadOpt)
		if err != nil {
			return nil, err
		}

		pkgMeta, err := ukgen.LoadPackage(params.PatternMeta, loadOpt)
		if err != nil {
			return nil, err
		}

		// TODO: Remove me
		_ = modRoot
		_ = pkgCore
		_ = pkgMeta

		creator, err := genBuildCreator(modRoot)
		if err != nil {
			return nil, err
		}

		transformer, err := genBuildTransformer(pkgCore, loadOpt)
		if err != nil {
			return nil, err
		}

		renderer, err := genBuildRenderer()
		if err != nil {
			return nil, err
		}

		opts = append(
			opts,
			genOpt(func(c *ukgen.Config) { c.Creator = creator }),
			genOpt(func(c *ukgen.Config) { c.Transformer = transformer }),
			genOpt(func(c *ukgen.Config) { c.Renderer = renderer }),
		)

		return opts, nil
	}

	return ukcli.Command[GenerateParams]{
		Exec: ukgen.NewExec(handle),
		Info: ukcli.NewInfo(info),
	}
}

type genOpt func(*ukgen.Config)

func (o genOpt) UkaseApplyGen(c *ukgen.Config) { o(c) }

func genBuildCreator(rootModule ukgen.Module) (ukgen.Creator, error) {
	// TODO
	fmt.Printf("rootModule: %+v\n", rootModule)
	return ukgen.FSRoot(rootModule.Dir), nil
	// return ukgen.CreateStdout, nil
}

func genBuildTransformer(corePkg ukgen.Package, opts ...func(*packages.Config)) (ukgen.Transformer, error) {

	abc := func(pkgPath string) (ukgen.Package, error) {
		sourcePkg, err := ukgen.LoadPackage(pkgPath, opts...)
		if err != nil {
			return ukgen.Package{}, err
		}

		sinkPkg := ukgen.Package{
			Dir:  filepath.Join(corePkg.Dir, "slug", sourcePkg.Name),
			Name: sourcePkg.Name,
			Path: path.Join(corePkg.Path, "slug", sourcePkg.Name),
		}

		return sinkPkg, nil
	}

	return ukgen.TransformPackage(abc), nil
}

func genBuildRenderer() (ukgen.Renderer, error) {
	// TODO
	return ukgen.RenderCore("info.go"), nil
	// return ukgen.RenderJSON("info.json"), nil
}
