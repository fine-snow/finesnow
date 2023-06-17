// Prefix route tree

package router

import (
	"github.com/fine-snow/finesnow/constant"
)

// PrefixRouteTreeAbstract Prefix route tree node abstract
type prefixRouteTreeAbstract interface {
	insert([]string, int)
	search([]string, int) string
	matchNode(part string) prefixRouteTreeAbstract
	matchNodes(part string) []prefixRouteTreeAbstract
	getUrl() string
	getPart() string
	getIsVar() bool
	setIsExist(bool)
}

// prefixRouteTree Global prefix route tree
var prefixRouteTree prefixRouteTreeAbstract = &node{
	url:      constant.NullStr,
	part:     constant.NullStr,
	children: make([]prefixRouteTreeAbstract, constant.Zero),
	isVar:    false}

// node prefix route tree node achieve
// url Node Real Routing Path
// part Node Partial Routing Path
// children Node Sub-node Collection
// isVar Is the node part path a variable
// isExist Does this segment of the route really exist
type node struct {
	url      string
	part     string
	children []prefixRouteTreeAbstract
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

// insert Add nodes to the routing tree
func (n *node) insert(parts []string, depth int) {
	part := parts[depth]
	nd := n.matchNode(part)
	if nd == nil {
		nd = &node{url: n.url + constant.Slash + part, part: part, isVar: part[constant.Zero] == constant.Colon}
		n.children = append(n.children, nd)
	}
	if len(parts) == (depth + constant.One) {
		nd.setIsExist(true)
		return
	}
	nd.insert(parts, depth+constant.One)
}

// search Query the real URL through route tree matching
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
		url := nd.search(parts, depth+constant.One)
		if url != constant.NullStr {
			return url
		}
	}
	return constant.NullStr
}

// matchNode Matches a single node when a node is inserted
func (n *node) matchNode(part string) prefixRouteTreeAbstract {
	for _, child := range n.children {
		if child.getPart() == part {
			return child
		}
	}
	return nil
}

// matchNodes Multiple nodes are matched when querying for real URLs through the routing tree
func (n *node) matchNodes(part string) []prefixRouteTreeAbstract {
	nodes := make([]prefixRouteTreeAbstract, constant.Zero)
	for _, child := range n.children {
		if child.getPart() == part || child.getIsVar() {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
