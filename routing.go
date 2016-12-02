package routing

import (
	"fmt"
	"regexp"
	"strings"
)

// Version is this package's version number.
const Version = "0.0.1"

// Node represents a node in a trie.
type Node struct {
	name     string
	str      string
	reg      *regexp.Regexp
	parent   *Node
	child    map[string]*Node
	children []*Node
}

// New returns a new root Node.
func New() *Node {
	return &Node{parent: nil}
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
	if opt.str != "" {
		c := n.child[opt.str]

		if c.str == opt.str && c.name == opt.name {
			return c
		}

		return nil
	}

	if opt.name != "" {
		for _, c := range n.children {
			if c.name == opt.name {
				return c
			}
		}
	}

	return nil
}

func (n *Node) attach(opt option) *Node {
	node := &Node{
		name:   opt.name,
		str:    opt.str,
		reg:    opt.reg,
		parent: n,
	}

	if opt.str != "" {
		node.child = map[string]*Node{opt.str: node}
	} else {
		node.children = []*Node{node}
	}

	return node
}
