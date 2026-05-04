package uktree

import (
	"errors"

	"github.com/oligarch316/ukase/ukcore"
)

func Read[T any](tree ukcore.Tree[T], target ...string) (T, error) {
	if node, exists := tree.Lookup(target...); exists {
		return node.Load()
	}

	return *new(T), errors.New("[TODO Read] target not found")
}

func ReadChildren[T any](tree ukcore.Tree[T], target ...string) (map[string]T, error) {
	node, exists := tree.Lookup(target...)
	if !exists {
		return nil, errors.New("[TODO ReadChildren] target not found")
	}

	children := make(map[string]T)
	errs := make([]error, 0)

	for name, child := range node.List() {
		entry, err := child.Load()
		children[name] = entry
		errs = append(errs, err)
	}

	return children, errors.Join(errs...)
}

func Update[T any](node *Node[T], update func(*T) error, target ...string) error {
	node = ensure(node, target)
	return update(&node.Entry)
}

func Write[T any](node *Node[T], entry T, target ...string) {
	node = ensure(node, target)
	node.Entry = entry
}

func ensure[T any](node *Node[T], target []string) *Node[T] {
	for _, name := range target {
		child, exists := node.Children[name]
		if !exists {
			child = New[T]()
			node.Children[name] = child
		}

		node = child
	}

	return node
}
