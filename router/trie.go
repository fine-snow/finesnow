// Trie route tree

package router

import (
	"github.com/fine-snow/finesnow/constant"
	"net/http"
)

// TrieRouteTreeAbstract Trie route tree node abstract
type trieRouteTreeAbstract interface {
	insert([]string, int)
	search([]string, int, *http.Request) string
	matchNode(part string) trieRouteTreeAbstract
	matchNodes(part string) []trieRouteTreeAbstract
	getUrl() string
	getPart() string
	getIsVar() bool
}

// trieRouteTree Global trie route tree
var trieRouteTree trieRouteTreeAbstract = &node{
	url:      constant.NullStr,
	part:     constant.NullStr,
	children: make([]trieRouteTreeAbstract, 0),
	isVar:    false}

// node Trie route tree node achieve
type node struct {
	url      string
	part     string
	children []trieRouteTreeAbstract
	isVar    bool
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

func (n *node) insert(parts []string, depth int) {
	part := parts[depth]
	nd := n.matchNode(part)
	if nd == nil {
		nd = &node{url: n.url + constant.Slash + part, part: part, isVar: part[0] == constant.Colon}
		n.children = append(n.children, nd)
	}
	if len(parts) == (depth + 1) {
		return
	}
	nd.insert(parts, depth+1)
}

func (n *node) search(parts []string, depth int, r *http.Request) string {
	part := parts[depth]
	nodes := n.matchNodes(part)
	if len(nodes) > 0 {
		nd := nodes[0]
		if nd.getIsVar() {
			if r.URL.RawQuery == constant.NullStr {
				r.URL.Query().Set(nd.getPart()[1:], part)
				r.URL.RawQuery = nd.getPart()[1:] + constant.EqualSign + part
			} else {
				r.URL.RawQuery = r.URL.RawQuery + constant.Ampersand + nd.getPart()[1:] + constant.EqualSign + part
			}
		}
		if len(parts) == (depth + 1) {
			return nd.getUrl()
		}
		url := nd.search(parts, depth+1, r)
		if url != constant.NullStr {
			return url
		}
	}
	return constant.NullStr
}

func (n *node) matchNode(part string) trieRouteTreeAbstract {
	for _, child := range n.children {
		if child.getPart() == part || child.getIsVar() {
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
