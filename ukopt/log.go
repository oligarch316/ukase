package ukopt

import (
	"log/slog"

	"github.com/oligarch316/ukase"
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore/ukdec"
	"github.com/oligarch316/ukase/ukcore/ukexec"
	"github.com/oligarch316/ukase/ukcore/ukinit"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

var (
	_ ukase.Option  = Log{}
	_ ukcli.Option  = Log{}
	_ ukdec.Option  = Log{}
	_ ukexec.Option = Log{}
	_ ukinit.Option = Log{}
	_ ukspec.Option = Log{}
	_ ukase.Option  = Log{}
)

const logKey = "ukase"

type Log struct{ *slog.Logger }

func (o Log) UkaseApplyDec(c *ukdec.Config)   { /* TODO */ }
func (o Log) UkaseApplyExec(c *ukexec.Config) { c.Log = o.with("exec") }
func (o Log) UkaseApplyInit(c *ukinit.Config) { /* TODO */ }
func (o Log) UkaseApplySpec(c *ukspec.Config) { /* TODO */ }

func (o Log) UkaseApplyCLI(c *ukcli.Config) {
	c.Log = o.with("cli")

	c.Exec = append(c.Exec, o)
	c.Decode = append(c.Decode, o)
	c.Init = append(c.Init, o)
	c.Spec = append(c.Spec, o)
}

func (o Log) UkaseApplyApp(c *ukase.Config) {
	c.Log = o.with("app")

	c.CLI = append(c.CLI, o)
}

func (o Log) with(name string) *slog.Logger { return o.Logger.With(logKey, name) }

func LogDefault() Log                     { return Log{Logger: slog.Default()} }
func LogLogger(logger *slog.Logger) Log   { return Log{Logger: logger} }
func LogHandler(handler slog.Handler) Log { return Log{Logger: slog.New(handler)} }
