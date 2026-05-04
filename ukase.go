package ukase

import (
	"context"
	"fmt"

	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcli/ukdecode"
	"github.com/oligarch316/ukase/ukcli/ukdirect"
	"github.com/oligarch316/ukase/ukcli/ukparse"
	"github.com/oligarch316/ukase/ukcore/uktree"
)

// =============================================================================
// App
// =============================================================================

type App struct {
	Decoder ukcli.Decoder
	Parser  ukcli.Parser
	runtime *ukcli.Runtime
}

func New(opts ...Option) *App {
	config := newConfig(opts)

	return &App{
		Decoder: ukdecode.New(config.Decode...),
		Parser:  ukparse.New(config.Parse...),
		runtime: ukcli.New(config.CLI...),
	}
}

func (a *App) Run(ctx context.Context, args []string) error {
	cliCtx, err := a.runtime.Build(ctx)
	if err != nil {
		return err
	}

	input, err := a.Parser.Parse(cliCtx, args)
	if err != nil {
		return err
	}

	entry, err := uktree.Read(cliCtx, input.Target...)
	if err != nil {
		return err
	}

	if entry.Exec == nil {
		return fmt.Errorf("[TODO App.execute] nil entry.Exec at target '%v'", input.Target)
	}

	return entry.Exec.Execute(cliCtx, input)
}

// -----------------------------------------------------------------------------
// Config
// -----------------------------------------------------------------------------

type Option interface{ UkaseApply(*Config) }

type Config struct {
	CLI    []ukcli.Option
	Decode []ukdecode.Option
	Parse  []ukparse.Option
}

func newConfig(opts []Option) Config {
	config := cfgDefault
	for _, opt := range opts {
		opt.UkaseApply(&config)
	}
	return config
}

var cfgDefault = Config{}

// -----------------------------------------------------------------------------
// Directives
// -----------------------------------------------------------------------------

func Add(app *App, directives ...ukcli.Directive) {
	app.runtime.Add(directives...)
}

func AddAuto(app *App, operation ukdirect.Operation, name string) {
	directive := ukdirect.Auto(operation, name)
	app.runtime.Add(directive)
}

func AddGlobal(app *App, operation ukdirect.Operation) {
	directive := ukdirect.Global(operation)
	app.runtime.Add(directive)
}

func AddLocal(app *App, operation ukdirect.Operation, target ...string) {
	directive := ukdirect.Local(operation, target...)
	app.runtime.Add(directive)
}

func AddInfo(app *App, info string, target ...string) {
	operation := ukdirect.Info(info)
	directive := ukdirect.Local(operation, target...)
	app.runtime.Add(directive)
}

func AddHandler[Params any](app *App, handler func(context.Context, Params) error, target ...string) {
	operation := ukdirect.Handler(app.Decoder, handler)
	directive := ukdirect.Local(operation, target...)
	app.runtime.Add(directive)
}

func AddCommand[Params any](app *App, info string, handler func(context.Context, Params) error, target ...string) {
	operation := ukdirect.Command[Params]{
		Info: ukdirect.Info(info),
		Exec: ukdirect.Handler(app.Decoder, handler),
	}

	directive := ukdirect.Local(operation, target...)
	app.runtime.Add(directive)
}
