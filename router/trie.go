// Trie route tree

package router

import "strings"

var trieRouteTree = node{url: slash, part: slash, children: make([]*node, 0), isVar: false}

type node struct {
	url      string
	part     string
	children []*node
	isVar    bool
}

func (n *node) insert(parts []string, url string, depth int) {
	part := parts[depth]
	nd := n.matchNode(part)
	if nd == nil {
		nd = &node{url: strings.Join([]string{n.url, part}, slash), part: part, isVar: part[0] == ':'}
		n.children = append(n.children, nd)
	}
	nd.insert(parts, url, depth+1)
}

func (n *node) search(parts []string, depth int) *node {
	if len(parts) == depth {
		if n.part == "" {
			return nil
		}
		return n
	}
	part := parts[depth]
	nodes := n.matchNodes(part)
	for _, nd := range nodes {
		next := nd.search(parts, depth+1)
		if next != nil {
			return next
		}
	}
	return nil
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
