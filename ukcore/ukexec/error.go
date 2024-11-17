package ukexec

import (
	"errors"
	"fmt"

	"github.com/oligarch316/ukase/internal/ierror"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

var errIsTagged = ierror.IsTaggedFunc(ierror.ErrDec)

// =============================================================================
// Conflict
// =============================================================================

var ErrConflict = errors.New("conflict error")

type ExecConflictError struct {
	Target           []string
	Original, Update ukspec.Parameters
	err              error
}

type InfoConflictError struct {
	Target           []string
	Original, Update any
	err              error
}

type FlagConflictError struct {
	Target           []string
	Name             string
	Original, Update ukspec.Flag
	err              error
}

func (ExecConflictError) Is(t error) bool { return errIsTagged(t, ErrConflict) }
func (InfoConflictError) Is(t error) bool { return errIsTagged(t, ErrConflict) }
func (FlagConflictError) Is(t error) bool { return errIsTagged(t, ErrConflict) }

func (e ExecConflictError) Unwrap() error { return e.err }
func (e InfoConflictError) Unwrap() error { return e.err }
func (e FlagConflictError) Unwrap() error { return e.err }

func (e ExecConflictError) Error() string {
	return fmt.Sprintf("exec conflict for target '%s': %s", e.Target, e.err)
}

func (e InfoConflictError) Error() string {
	return fmt.Sprintf("info conflict for target '%s': %s", e.Target, e.err)
}

func (e FlagConflictError) Error() string {
	return fmt.Sprintf("flag (%s) conflict for target '%s': %s", e.Name, e.Target, e.err)
}

// =============================================================================
// Parse
// =============================================================================

var ErrParse = errors.New("parse error")

type ParseError struct {
	Target   []string
	Position int
	err      error
}

func (ParseError) Is(t error) bool { return errIsTagged(t, ErrParse) }

func (e ParseError) Unwrap() error { return e.err }

func (e ParseError) Error() string { return e.err.Error() }
