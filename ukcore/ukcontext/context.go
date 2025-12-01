package ukcontext

import (
	"context"
	"time"

	"github.com/oligarch316/ukase/ukcore"
	"github.com/oligarch316/ukase/ukcore/uktree"
)

// =============================================================================
// Context
// =============================================================================

type treeContext[T any] struct {
	context.Context
	ukcore.Tree[T]
}

func New[T any](ctx context.Context, tree ukcore.Tree[T]) ukcore.Context[T] {
	return treeContext[T]{Context: ctx, Tree: tree}
}

// -----------------------------------------------------------------------------
// Derive
// > Persist underlying context
// > Update underlying tree
// -----------------------------------------------------------------------------

func Derive[T, U any](parent ukcore.Context[T], derive func(T) (U, error)) ukcore.Context[U] {
	tree := uktree.Derive(parent, derive)
	return treeContext[U]{Context: parent, Tree: tree}
}

func DeriveTarget[T, U any](parent ukcore.Context[T], derive func([]string, T) (U, error)) ukcore.Context[U] {
	tree := uktree.DeriveTarget(parent, derive)
	return treeContext[U]{Context: parent, Tree: tree}
}

// -----------------------------------------------------------------------------
// Cancellation
// > Update underlying context
// > Persist underlying tree
// -----------------------------------------------------------------------------

func WithCancel[T any](parent ukcore.Context[T]) (ukcore.Context[T], context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	return treeContext[T]{Context: ctx, Tree: parent}, cancel
}

func WithCancelCause[T any](parent ukcore.Context[T]) (ukcore.Context[T], context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(parent)
	return treeContext[T]{Context: ctx, Tree: parent}, cancel
}

func WithDeadline[T any](parent ukcore.Context[T], d time.Time) (ukcore.Context[T], context.CancelFunc) {
	ctx, cancel := context.WithDeadline(parent, d)
	return treeContext[T]{Context: ctx, Tree: parent}, cancel
}

func WithDeadlineCause[T any](parent ukcore.Context[T], d time.Time, cause error) (ukcore.Context[T], context.CancelFunc) {
	ctx, cancel := context.WithDeadlineCause(parent, d, cause)
	return treeContext[T]{Context: ctx, Tree: parent}, cancel
}

func WithTimeout[T any](parent ukcore.Context[T], timeout time.Duration) (ukcore.Context[T], context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	return treeContext[T]{Context: ctx, Tree: parent}, cancel
}

func WithTimeoutCause[T any](parent ukcore.Context[T], timeout time.Duration, cause error) (ukcore.Context[T], context.CancelFunc) {
	ctx, cancel := context.WithTimeoutCause(parent, timeout, cause)
	return treeContext[T]{Context: ctx, Tree: parent}, cancel
}
