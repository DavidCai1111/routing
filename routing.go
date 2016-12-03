package routing

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Version is this package's version number.
const Version = "1.1.0"

var (
	normalStrReg    = regexp.MustCompile(`^[\w\.-]+$`)
	separatedStrReg = regexp.MustCompile(`^[\w\.\-][\w\.\-\|]+[\w\.\-]$`)
	nameStrReg      = regexp.MustCompile(`^\:\w+\b`)
	surroundStrReg  = regexp.MustCompile(`^\(.+\)$`)
)

type option struct {
	name string
	str  string
	reg  string
}

// Node represents a node in the trie tree. In this package, you
// can only get the root node by the New() func. And then define and
// match URL with it.
type Node struct {
	callback interface{}
	parent   *Node
	children map[option]*Node
}

// New returns a new root Node. you can define and match URL with it.
func New() *Node {
	return &Node{
		parent:   nil,
		children: map[option]*Node{},
	}
}

// Define defines a URL and some meta infomation (interface{}) on it.
//
// support signatures:
//
// - string: `/hello`
//
// - separated string: `/a|b|c`
//
// - regex: `/([0-9a-f]{24})`
//
// - named parameter: `/:id`
//
// - named separated string: `/:id(a|b|c)`
//
// - named regex: `/:id([0-9a-f]{24})`
//
func (n *Node) Define(u string, callback interface{}) {
	n.define(strings.Split(checkURL(u), "/")[1:], callback)
}

func (n *Node) define(frags []string, callback interface{}) {
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
		for _, node := range nodes {
			node.callback = callback
		}

		return
	}

	for _, node := range nodes {
		node.define(frags[1:], callback)
	}
}

// Match matchs a URL and get the meta infomation defined on it. You should
// always check the third return value at first to know whether the URL is
// matched. The first return value is the URL's defined meta infomation and the
// second return value is the map of named parameters.
func (n *Node) Match(u string) (interface{}, map[string]string, bool) {
	u, err := url.QueryUnescape(u)

	if err != nil {
		return nil, nil, false
	}

	return n.match(map[string]string{}, strings.Split(checkURL(u), "/")[1:])
}

func (n *Node) match(p map[string]string, frags []string) (interface{}, map[string]string, bool) {
	frag := frags[0]

	for opt, child := range n.children {
		if opt.str == frag ||
			(opt.reg != "" && regexp.MustCompile(opt.reg).MatchString(frag)) ||
			opt.name != "" {
			if opt.name != "" {
				p[opt.name] = frag
			}

			if len(frags) == 1 {
				return child.callback, p, true
			}

			if cb, p, ok := child.match(p, frags[1:]); ok {
				return cb, p, ok
			}
		}
	}

	return nil, nil, false
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

func parse(frag string) []option {
	if frag == "" || normalStrReg.MatchString(frag) {
		return []option{option{str: frag}}
	}

	if separatedStrReg.MatchString(frag) {
		separated := strings.Split(frag, "|")
		options := make([]option, len(separated))

		for i, s := range separated {
			options[i].str = s
		}

		return options
	}

	var name string

	frag = nameStrReg.ReplaceAllStringFunc(frag, func(n string) string {
		name = n[1:]
		return ""
	})

	if len(frag) == 0 {
		return []option{option{name: name}}
	}

	if surroundStrReg.MatchString(frag) {
		frag = frag[1 : len(frag)-1]

		if separatedStrReg.MatchString(frag) {
			separated := strings.Split(frag, "|")
			options := make([]option, len(separated))

			for i, s := range separated {
				options[i].name = name
				options[i].str = s
			}

			return options
		}

		return []option{option{name: name, reg: regexp.MustCompile(frag).String()}}
	}

	panic(fmt.Sprintf("routing: Invalid frag: %v", frag))
}

func checkURL(url string) string {
	if url[0] != '/' {
		panic(fmt.Sprintf("routing: %v should start with '/'\n", url))
	}

	return url
}
