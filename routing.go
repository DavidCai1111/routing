package routing

import (
	"fmt"
	"strings"
)

// Version is this package's version number.
const Version = "0.0.1"

// Node represents a node in a trie.
type Node struct {
	option
	parent   *Node
	children map[option]*Node
}

// New returns a new root Node.
func New() *Node {
	return &Node{
		parent:   nil,
		children: map[option]*Node{},
	}
}

// Define defines a url
func (n *Node) Define(url string) {
	if url[0] != '/' {
		panic(fmt.Sprintf("routing: %v should start with '/'\n", url))
	}

	n.define(strings.Split(url, "/")[1:])
}

func (n *Node) define(frags []string) {
	options := parse(frags[0])

	nodes := []*Node{}

	for _, opt := range options {
		node := n.find(opt)

		if node == nil {
			node = n.attach(opt)
		}

		nodes = append(nodes, node)
	}

	if len(frags) == 1 {
		return
	}

	for _, node := range nodes {
		node.define(frags[1:])
	}
}

func (n *Node) find(opt option) *Node {
	if c, ok := n.children[opt]; ok {
		return c
	}

	return nil
}

func (n *Node) attach(opt option) *Node {
	node := &Node{
		option:   opt,
		parent:   n,
		children: map[option]*Node{},
	}

	n.children[opt] = node

	return node
}
