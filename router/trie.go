// Trie route tree

package router

import (
	"github.com/fine-snow/finesnow/constant"
)

// TrieRouteTreeAbstract Trie route tree node abstract
type trieRouteTreeAbstract interface {
	insert([]string, int)
	search([]string, int) string
	matchNode(part string) trieRouteTreeAbstract
	matchNodes(part string) []trieRouteTreeAbstract
	getUrl() string
	getPart() string
	getIsVar() bool
	setIsExist(bool)
}

// trieRouteTree Global trie route tree
var trieRouteTree trieRouteTreeAbstract = &node{
	url:      constant.NullStr,
	part:     constant.NullStr,
	children: make([]trieRouteTreeAbstract, 0),
	isVar:    false}

// node Trie route tree node achieve
// url Node Real Routing Path
// part Node Partial Routing Path
// children Node Sub-node Collection
// isVar Is the node part path a variable
// isExist Does this segment of the route really exist
type node struct {
	url      string
	part     string
	children []trieRouteTreeAbstract
	isVar    bool
	isExist  bool
}

func (n *node) getUrl() string {
	return n.url
}

func (n *node) getPart() string {
	return n.part
}

func (n *node) getIsVar() bool {
	return n.isVar
}

func (n *node) setIsExist(b bool) {
	n.isExist = b
}

func (n *node) insert(parts []string, depth int) {
	part := parts[depth]
	nd := n.matchNode(part)
	if nd == nil {
		nd = &node{url: n.url + constant.Slash + part, part: part, isVar: part[0] == constant.Colon}
		n.children = append(n.children, nd)
	}
	if len(parts) == (depth + 1) {
		nd.setIsExist(true)
		return
	}
	nd.insert(parts, depth+1)
}

func (n *node) search(parts []string, depth int) string {
	if len(parts) == depth {
		if n.isExist {
			return n.url
		} else {
			return constant.NullStr
		}
	}
	part := parts[depth]
	nodes := n.matchNodes(part)
	for _, nd := range nodes {
		url := nd.search(parts, depth+1)
		if url != constant.NullStr {
			return url
		}
	}
	return constant.NullStr
}

func (n *node) matchNode(part string) trieRouteTreeAbstract {
	for _, child := range n.children {
		if child.getPart() == part {
			return child
		}
	}
	return nil
}

func (n *node) matchNodes(part string) []trieRouteTreeAbstract {
	nodes := make([]trieRouteTreeAbstract, 0)
	for _, child := range n.children {
		if child.getPart() == part || child.getIsVar() {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
