package ukinput

import (
	"errors"
	"fmt"

	"github.com/oligarch316/ukase/ukcore"
)

// =============================================================================
// Go Style
// =============================================================================

type StyleGo struct{}

// -----------------------------------------------------------------------------
// Format
// -----------------------------------------------------------------------------

func (StyleGo) Format(flagName string) string { return "-" + flagName }

// -----------------------------------------------------------------------------
// Parse
// -----------------------------------------------------------------------------

func (sg StyleGo) Parse(ctx Context, args []string) (ukcore.Input, error) {
	state := &styleState{Args: args, Entries: ctx}

	// ----- Program
	if err := sg.parseProgram(state); err != nil {
		return state.Input, err
	}

	// ----- Targets + Flags
	for tkn, more := sg.lex(state); more; tkn, more = sg.lex(state) {
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
			if err := sg.parseTarget(state); err != nil {
				return state.Input, err
			}

			continue
		}

		// Flag ⇒ consume and continue
		if tkn == tokenFlag {
			if err := sg.parseFlag(state); err != nil {
				return state.Input, err
			}

			continue
		}

		// Unknown ⇒ error
		return state.Input, fmt.Errorf("[TODO StyleGo.Parse] unknown token: %d", tkn)
	}

	// ----- Arguments
	sg.parseArguments(state)

	return state.Input, nil
}

func (sg StyleGo) lex(state *styleState) (token, bool) {
	peekStr, peekOk := state.peek()

	switch {
	case !peekOk:
		return tokenInvalid, false
	case peekStr == "":
		return tokenInvalid, true
	case peekStr == "--":
		return tokenDelim, true
	case len(peekStr) > 1 && peekStr[0] == '-':
		return tokenFlag, true
	}

	if _, targetOk := state.Entries.Lookup(peekStr); targetOk {
		return tokenTarget, true
	}

	return tokenArgument, true
}

func (StyleGo) parseProgram(state *styleState) error {
	program, ok := state.pop()
	if !ok {
		return errors.New("[TODO StyleGo.parseProgram] missing program")
	}

	state.Input.Program = program
	return nil
}

func (StyleGo) parseTarget(state *styleState) error {
	target, targetOk := state.pop()
	if !targetOk {
		return errors.New("[TODO StyleGo.parseTarget] missing target")
	}

	child, childOk := state.Entries.Lookup(target)
	if !childOk {
		return errors.New("[TODO StyleGo.parseTarget] unknown target")
	}

	state.Input.Target = append(state.Input.Target, target)
	state.Entries = child
	return nil
}

func (StyleGo) parseFlag(state *styleState) error {
	// TODO
	return errors.New("[TODO StyleGo.parseFlag] not yet implemented")
}

func (StyleGo) parseArguments(state *styleState) {
	for pos, val := range state.Args {
		argument := ukcore.InputArgument{Position: pos, Value: val}
		state.Input.Arguments = append(state.Input.Arguments, argument)
	}

	state.Args = nil
}
