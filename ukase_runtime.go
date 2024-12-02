//go:build !ukgen

package ukase

import "github.com/oligarch316/ukase/ukcli"

type appConfig struct{ Config }

func (ac appConfig) UkaseApplyCLI(c *ukcli.Config) {
	ac.cliApplyHelpAuto(c)
	ac.cliApplyUser(c)
}

func (ac appConfig) cliApplyHelpAuto(c *ukcli.Config) {
	// TODO
}

func (ac appConfig) cliApplyUser(c *ukcli.Config) {
	for _, opt := range ac.CLI {
		opt.UkaseApplyCLI(c)
	}
}
