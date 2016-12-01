package routing

import "regexp"

// Version is this package's version number.
const Version = "0.0.1"

// Node represents a node in a trie.
type Node struct {
	name     string
	str      string
	reg      regexp.Regexp
	parent   *Node
	child    map[string]*Node
	children []*Node
}

// New returns a new root Node.
func New() *Node {
	return &Node{parent: nil}
}
