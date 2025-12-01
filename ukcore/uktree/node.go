package uktree

import (
	"iter"

	"github.com/oligarch316/ukase/ukcore"
)

type Node[T any] struct {
	Entry    T
	Children map[string]*Node[T]
}

func New[T any]() *Node[T] {
	children := make(map[string]*Node[T])
	return &Node[T]{Children: children}
}

func (n *Node[T]) Load() (T, error) { return n.Entry, nil }

func (n *Node[T]) List() iter.Seq2[string, ukcore.Tree[T]] {
	return func(yield func(string, ukcore.Tree[T]) bool) {
		for name, child := range n.Children {
			if !yield(name, child) {
				return
			}
		}
	}
}

func (n *Node[T]) Lookup(target ...string) (ukcore.Tree[T], bool) {
	tree, exists := n, true

	for _, name := range target {
		if tree, exists = tree.Children[name]; !exists {
			break
		}
	}

	return tree, exists
}
