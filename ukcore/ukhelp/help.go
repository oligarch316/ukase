package ukhelp

import (
	"errors"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukvalue"
)

// =============================================================================
// Help
// =============================================================================

type Context ukcore.Context[Entry]

type Entry struct {
	Command EntryCommand
	Fields  []EntryField
}

type EntryCommand struct {
	Info  string
	Usage string
	Exec  bool
}

type EntryField struct {
	Info          string
	FlagLabel     string
	ArgumentLabel string
}

type Values struct {
	Info          ukcore.SpecValue[string]
	FlagLabel     ukcore.SpecValue[string]
	ArgumentLabel ukcore.SpecValue[string]
}

func (v Values) Load(spec ukcore.Spec, index []int) (EntryField, error) {
	info, errInfo := ukvalue.OrConst(v.Info, "").Load(spec, index)
	flag, errFlag := ukvalue.OrConst(v.FlagLabel, "").Load(spec, index)
	arg, errArg := ukvalue.OrConst(v.ArgumentLabel, "").Load(spec, index)

	entry := EntryField{Info: info, FlagLabel: flag, ArgumentLabel: arg}
	return entry, errors.Join(errInfo, errArg, errFlag)
}
