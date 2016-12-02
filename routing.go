package routing

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Version is this package's version number.
const Version = "0.0.1"

// Node represents a node in a trie.
type Node struct {
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
func (n *Node) Define(u string) {
	n.define(strings.Split(checkURL(u), "/")[1:])
}

// Match matchs a url
func (n *Node) Match(u string) (map[string]string, bool) {
	u, err := url.QueryUnescape(u)

	if err != nil {
		return nil, false
	}

	return n.match(map[string]string{}, strings.Split(checkURL(u), "/")[1:])
}

func (n *Node) match(p map[string]string, frags []string) (map[string]string, bool) {
	frag := frags[0]

	for opt, child := range n.children {
		if opt.str == frag ||
			(opt.reg != "" && regexp.MustCompile(opt.reg).MatchString(frag)) {
			if opt.name != "" {
				p[opt.name] = frag
			}

			if len(frags) == 1 {
				return p, true
			}

			if p, ok := child.match(p, frags[1:]); ok {
				return p, ok
			}
		}
	}

	return nil, false
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
		parent:   n,
		children: map[option]*Node{},
	}

	n.children[opt] = node

	return node
}

func checkURL(url string) string {
	if url[0] != '/' {
		panic(fmt.Sprintf("routing: %v should start with '/'\n", url))
	}

	return url
}
