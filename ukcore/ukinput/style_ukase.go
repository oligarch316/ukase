package ukinput

import (
	"errors"
	"fmt"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Ukase Style
// =============================================================================

type StyleUkase struct {
	// TODO: Document
	Delimiter string
}

// -----------------------------------------------------------------------------
// Format
// -----------------------------------------------------------------------------

func (StyleUkase) Format(flagName string) string {
	// TODO: Need to check rune length instead
	switch len(flagName) {
	case 0:
		return flagName
	case 1:
		return "-" + flagName
	default:
		return "--" + flagName
	}
}

// -----------------------------------------------------------------------------
// Parse
// -----------------------------------------------------------------------------

func (su StyleUkase) Parse(ctx Context, args []string) (ukcore.Input, error) {
	state := &styleState{Args: args, Entries: ctx}

	// ----- Program
	if err := su.parseProgram(state); err != nil {
		return state.Input, err
	}

	// ----- Targets + Flags
	for tkn, more := su.lex(state); more; tkn, more = su.lex(state) {
		// Argument ⇒ break
		if tkn == tokenArgument {
			break
		}

		// Delimiter ⇒ consume and break
		if tkn == tokenDelim {
			_, _ = state.pop()
			break
		}

		// Target ⇒ consume and continue
		if tkn == tokenTarget {
			if err := su.parseTarget(state); err != nil {
				return state.Input, err
			}

			continue
		}

		// Flag ⇒ consume and continue
		if tkn == tokenFlag {
			if err := su.parseFlag(state); err != nil {
				return state.Input, err
			}

			continue
		}

		// Unknown ⇒ error
		return state.Input, fmt.Errorf("[TODO StyleUkase.Parse] unknown token: %d", tkn)
	}

	// ----- Arguments
	su.parseArguments(state)

	return state.Input, nil
}

func (su StyleUkase) lex(state *styleState) (token, bool) {
	peekStr, peekOk := state.peek()

	switch {
	case !peekOk:
		return tokenInvalid, false
	case peekStr == "":
		return tokenInvalid, true
	case peekStr == su.Delimiter:
		return tokenDelim, true
	case len(peekStr) > 1 && peekStr[0] == '-':
		return tokenFlag, true
	}

	if _, targetOk := state.Entries.Lookup(peekStr); targetOk {
		return tokenTarget, true
	}

	return tokenArgument, true
}

func (StyleUkase) parseProgram(state *styleState) error {
	program, ok := state.pop()
	if !ok {
		return errors.New("[TODO StyleUkase.parseProgram] missing program")
	}

	state.Input.Program = program
	return nil
}

func (StyleUkase) parseTarget(state *styleState) error {
	target, targetOk := state.pop()
	if !targetOk {
		return errors.New("[TODO StyleUkase.parseTarget] missing target")
	}

	child, childOk := state.Entries.Lookup(target)
	if !childOk {
		return errors.New("[TODO StyleUkase.parseTarget] unknown target")
	}

	state.Input.Target = append(state.Input.Target, target)
	state.Entries = child
	return nil
}

func (su StyleUkase) parseFlag(state *styleState) error {
	name, err := su.parseFlagName(state)
	if err != nil {
		return err
	}

	value, err := su.parseFlagValue(state, name)
	if err != nil {
		return err
	}

	flag := ukcore.InputFlag{Name: name, Value: value}
	state.Input.Flags = append(state.Input.Flags, flag)
	return nil
}

func (StyleUkase) parseFlagName(state *styleState) (string, error) {
	name, nameOk := state.pop()
	if !nameOk {
		return "", errors.New("[TODO StyleUkase.parseFlagName] missing flag name")
	}

	switch nameLen := len(name); {
	case nameLen == 2 && name[0] == '-' && name[1] != '-':
		// Short Flag ⇒ -x
		return name[1:], nil
	case nameLen > 2 && name[0] == '-' && name[1] == '-' && name[2] != '-':
		// Long Flag  ⇒ --xx…
		return name[2:], nil
	default:
		// Invalid Flag
		return "", fmt.Errorf("[TODO StyleUkase.parseFlagName] invalid flag name '%s'", name)
	}
}

func (StyleUkase) parseFlagValue(state *styleState, name string) (string, error) {
	entry, err := state.Entries.Load()
	if err != nil {
		return "", err
	}

	kind, kindOk := entry[name]
	if !kindOk {
		return "", fmt.Errorf("[TODO StyleUkase.parseFlagValue] unknown flag name '%s'", name)
	}

	peekStr, peekOk := state.peek()
	usable := peekOk && kind.Allowed(peekStr)

	// Value is unusable and required ⇒ error
	if kind.Required() && !usable {
		return "", fmt.Errorf("[TODO StyleUkase.parseFlagValue] missing value for flag '%s'", name)
	}

	// Value is unusable and NOT required ⇒ ignore and use assumed
	if !kind.Required() && !usable {
		return kind.Assumed(), nil
	}

	// Value is usable ⇒ consume and use
	value, valueOk := state.pop()
	if !valueOk {
		return "", fmt.Errorf("[TODO StyleUkase.parseFlagValue] missing value for flag '%s'", name)
	}

	return value, nil
}

func (StyleUkase) parseArguments(state *styleState) {
	for pos, val := range state.Args {
		argument := ukcore.InputArgument{Position: pos, Value: val}
		state.Input.Arguments = append(state.Input.Arguments, argument)
	}

	state.Args = nil
}
