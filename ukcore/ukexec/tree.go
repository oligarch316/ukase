package ukexec

import (
	"context"
	"fmt"
	"maps"

	"github.com/oligarch316/ukase/internal/ierror"
	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/ukspec"
)

// =============================================================================
// Tree
// =============================================================================

// TODO: Document
type Tree struct {
	config Config
	root   *entry
}

// TODO: Document
func NewTree(opts ...Option) *Tree {
	return &Tree{
		config: newConfig(opts),
		root:   newEntry(),
	}
}

// -----------------------------------------------------------------------------
// ❭ Load
// -----------------------------------------------------------------------------

// TODO: Document
func (t *Tree) LoadEntry(target ...string) (Entry, bool) {
	return t.root.Child(target...)
}

// -----------------------------------------------------------------------------
// ❭ Add
// -----------------------------------------------------------------------------

// TODO: Document
func (t *Tree) AddExec(exec ukcore.Exec, spec ukspec.Parameters, target ...string) error {
	t.config.Log.Debug("adding exec", "target", target, "specType", spec.Type)

	flags := make(map[string]ukspec.Flag)

	for _, update := range spec.Flags {
		for _, name := range update.Names {
			flags[name] = update

			original, conflict := t.root.flags[name]
			if !conflict {
				continue
			}

			if err := t.config.FlagConflict(original, update); err != nil {
				return FlagConflictError{
					Target:   target,
					Name:     name,
					Original: original,
					Update:   update,
					err:      ierror.D(err),
				}
			}
		}
	}

	entry := t.ensureTarget(target, flags)

	if entry.exec == nil {
		entry.exec, entry.spec = exec, spec
		return nil
	}

	overwrite, err := t.config.ExecConflict(entry.spec, spec)
	if err != nil {
		return ExecConflictError{
			Target:   target,
			Original: entry.spec,
			Update:   spec,
			err:      ierror.D(err),
		}
	}

	if overwrite {
		entry.exec, entry.spec = exec, spec
	}

	return nil
}

// TODO: Document
func (t *Tree) AddInfo(info any, target ...string) error {
	t.config.Log.Debug("adding info", "target", target, "infoType", fmt.Sprintf("%T", info))

	entry := t.ensureTarget(target, nil)

	if entry.info == nil {
		entry.info = info
		return nil
	}

	overwrite, err := t.config.InfoConflict(entry.info, info)
	if err != nil {
		return InfoConflictError{
			Target:   target,
			Original: entry.info,
			Update:   info,
			err:      ierror.D(err),
		}
	}

	if overwrite {
		entry.info = info
	}

	return nil
}

func (t *Tree) ensureTarget(target []string, flags map[string]ukspec.Flag) *entry {
	entry := t.root
	maps.Copy(entry.flags, flags)

	for _, name := range target {
		child, ok := entry.children[name]
		if !ok {
			child = newEntry()
			entry.children[name] = child
		}

		entry = child
		maps.Copy(entry.flags, flags)
	}

	return entry
}

// -----------------------------------------------------------------------------
// ❭ Execute
// -----------------------------------------------------------------------------

// TODO: Document
func (t *Tree) Execute(ctx context.Context, values []string) error {
	parser := newParser(values)

	program, ok := parser.ConsumeValue()
	if !ok {
		return ParseError{err: ierror.NewD("missing program")}
	}

	input := ukcore.Input{Program: program}
	entry := t.root

	for {
		// Consume all flags for the current entry
		flags, err := parser.ConsumeFlags(entry.flags)
		if err != nil {
			return ParseError{Target: input.Target, Position: parser.Position, err: err}
		}

		input.Flags = append(input.Flags, flags...)

		// Consume the next token of kind ...
		token := parser.ConsumeToken()

		// ... ❬Delim❭ or ❬EOF❭ ⇒ break out to argument parsing
		if token.Kind == kindDelim || token.Kind == kindEOF {
			break
		}

		// ... non-subcommand ⇒ set as 1st argument and break out to argument parsing
		child, ok := entry.children[token.Value]
		if !ok {
			input.Arguments = t.appendArguments(input.Arguments, token.Value)
			break
		}

		// ... subcommand ⇒ append command name to target and continue
		input.Target = append(input.Target, token.Value)
		entry = child
	}

	// All remaining unconsumed values are treated as arguments
	input.Arguments = t.appendArguments(input.Arguments, parser.Values...)

	t.config.Log.Info("executing", "target", input.Target)

	if entry.exec == nil {
		return t.config.ExecUnspecified(ctx, input)
	}

	return entry.exec(ctx, input)
}

func (Tree) appendArguments(args []ukcore.Argument, values ...string) []ukcore.Argument {
	pos := len(args)
	for _, value := range values {
		args = append(args, ukcore.Argument{Position: pos, Value: value})
		pos += 1
	}
	return args
}
