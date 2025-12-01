package ukinput

import "github.com/oligarch316/ukase/ukcore"

type Style interface {
	Format(string) string
	Parse(Context, []string) (ukcore.Input, error)
}

// -----

type token int

const (
	tokenInvalid token = iota
	tokenDelim
	tokenFlag
	tokenTarget
	tokenArgument
)

type styleState struct {
	Input   ukcore.Input
	Args    []string
	Entries ukcore.Tree[Entry]
}

func (ss *styleState) peek() (string, bool) {
	if len(ss.Args) == 0 {
		return "", false
	}

	return ss.Args[0], true
}

func (ss *styleState) pop() (str string, ok bool) {
	if str, ok = ss.peek(); ok {
		ss.Args = ss.Args[1:]
	}

	return
}
