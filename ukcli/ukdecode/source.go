package ukdecode

import (
	"github.com/oligarch316/ukase/ukcli"
	"github.com/oligarch316/ukase/ukcore/ukparams"
)

// =============================================================================
// Source
// =============================================================================

type Source func(ukcli.Values) ukparams.Source

var SourceEnv Source = sourceEnv
var SourceInit Source = sourceInit
var SourceInput Source = sourceInput

func sourceEnv(vals ukcli.Values) ukparams.Source {
	return ukparams.SourceEnv{EnvNames: vals.EnvNames}
}

func sourceInit(vals ukcli.Values) ukparams.Source {
	return ukparams.SourceInit{Initial: vals.Initial}
}

func sourceInput(vals ukcli.Values) ukparams.Source {
	return ukparams.SourceList{
		ukparams.SourceMeta{MetaNames: vals.MetaNames},
		ukparams.SourceArgs{ArgRange: vals.ArgRange},
		ukparams.SourceFlags{FlagNames: vals.FlagNames},
	}
}
