// Trie route tree

package router

import (
	"github.com/fine-snow/finesnow/constant"
	"net/http"
)

var trieRouteTree = node{url: constant.NullStr, part: constant.NullStr, children: make([]*node, 0), isVar: false}

type node struct {
	url      string
	part     string
	children []*node
	isVar    bool
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
		if nd.isVar {
			if r.URL.RawQuery == constant.NullStr {
				r.URL.Query().Set(nd.part[1:], part)
				r.URL.RawQuery = nd.part[1:] + constant.EqualSign + part
			} else {
				r.URL.RawQuery = r.URL.RawQuery + constant.Ampersand + nd.part[1:] + constant.EqualSign + part
			}
		}
		if len(parts) == (depth + 1) {
			return nd.url
		}
		url := nd.search(parts, depth+1, r)
		if url != "" {
			return url
		}
	}
	return ""
}

func (n *node) matchNode(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isVar {
			return child
		}
	}
	return nil
}

func (n *node) matchNodes(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isVar {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
