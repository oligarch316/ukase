package ukase

import (
	"context"
	"log/slog"
	"os"

	"github.com/oligarch316/ukase/internal/ilog"
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukmeta/ukgen"
	"github.com/oligarch316/ukase/ukmeta/ukhelp"
)

// =============================================================================
// Config
// =============================================================================

var cfgDefault = Config{
	Log:            ilog.Discard,
	HelpCommand:    "help",
	InputProgram:   os.Args[0],
	InputArguments: os.Args[1:],
	CLI:            nil,
	Help:           nil,
	Gen:            nil,
}

type Option interface{ UkaseApplyApp(*Config) }

type Config struct {
	// TODO: Document
	Log *slog.Logger

	// TODO: Document
	HelpCommand string

	// TODO: Document
	InputProgram string

	// TODO: Document
	InputArguments []string

	// TODO: Document
	CLI []ukcli.Option

	// TODO: Document
	Help []ukhelp.Option

	// TODO: Document
	Gen []ukgen.Option
}

func newConfig(opts []Option) appConfig {
	config := cfgDefault
	for _, opt := range opts {
		opt.UkaseApplyApp(&config)
	}
	return appConfig{Config: config}
}

// =============================================================================
// Application
// =============================================================================

type Application struct {
	config  appConfig
	runtime *ukcli.Runtime
}

func New(opts ...Option) *Application {
	config := newConfig(opts)
	runtime := ukcli.NewRuntime(config)

	return &Application{config: config, runtime: runtime}
}

func (a *Application) Add(directives ...ukcli.Directive) {
	a.runtime.Add(directives...)
}

func (a *Application) Run(ctx context.Context) error {
	values := []string{a.config.InputProgram}
	values = append(values, a.config.InputArguments...)

	return a.runtime.Execute(ctx, values)
}

// =============================================================================
// Directives
// =============================================================================

func Command[Params any](handler func(context.Context, Params) error, info any) ukcli.Command[Params] {
	return ukcli.NewCommand(handler, info)
}

func Exec[Params any](handler func(context.Context, Params) error) ukcli.Exec[Params] {
	return ukcli.NewExec(handler)
}

func Info(info any) ukcli.Info {
	return ukcli.NewInfo(info)
}

func Rule[Params any](rule func(*Params)) ukcli.Rule[Params] {
	return ukcli.NewRule(rule)
}
