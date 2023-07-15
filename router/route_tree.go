// Prefix route tree function

package router

// PrefixRouteTree Prefix route tree node abstract
type PrefixRouteTree interface {
	insert([]string, int)
	search([]string, int) string
	matchNode(part string) PrefixRouteTree
	matchNodes(part string) []PrefixRouteTree
	getUrl() string
	getPart() string
	getIsVar() bool
	setIsExist(bool)
}

// prefixRouteTree Global prefix route tree
var prefixRouteTree PrefixRouteTree = &treeNode{
	url:      "",
	part:     "",
	children: make([]PrefixRouteTree, 0),
	isVar:    false,
}

// treeNode prefix route tree node achieve
// url Node Real Routing Path
// part Node Partial Routing Path
// children Node Sub-node Collection
// isVar Is the node part path a variable
// isExist Does this segment of the route really exist
type treeNode struct {
	url      string
	part     string
	children []PrefixRouteTree
	isVar    bool
	isExist  bool
}

func (n *treeNode) getUrl() string {
	return n.url
}

func (n *treeNode) getPart() string {
	return n.part
}

func (n *treeNode) getIsVar() bool {
	return n.isVar
}

func (n *treeNode) setIsExist(b bool) {
	n.isExist = b
}

// insert Add nodes to the routing tree
func (n *treeNode) insert(parts []string, depth int) {
	part := parts[depth]
	nd := n.matchNode(part)
	if nd == nil {
		nd = &treeNode{url: n.url + "/" + part, part: part, isVar: part[0] == ':'}
		n.children = append(n.children, nd)
	}
	if len(parts) == (depth + 1) {
		nd.setIsExist(true)
		return
	}
	nd.insert(parts, depth+1)
}

// search Query the real URL through route tree matching
func (n *treeNode) search(parts []string, depth int) string {
	if len(parts) == depth {
		if n.isExist {
			return n.url
		} else {
			return ""
		}
	}
	part := parts[depth]
	nodes := n.matchNodes(part)
	for _, nd := range nodes {
		url := nd.search(parts, depth+1)
		if url != "" {
			return url
		}
	}
	return ""
}

// matchNode Matches a single node when a node is inserted
func (n *treeNode) matchNode(part string) PrefixRouteTree {
	for _, child := range n.children {
		if child.getPart() == part {
			return child
		}
	}
	return nil
}

// matchNodes Multiple nodes are matched when querying for real URLs through the routing tree
func (n *treeNode) matchNodes(part string) []PrefixRouteTree {
	nodes := make([]PrefixRouteTree, 0)
	for _, child := range n.children {
		if child.getPart() == part || child.getIsVar() {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
