package libvirtxml

type Node struct {
	Name       Name
	Attributes []*Attribute
	Nodes      []*Node
	CharData   string
	Comments   string
}

type Attribute struct {
	Name  Name
	Value string
}

func NewNode(name Name) *Node {
	return &Node{
		Name:       name,
		Attributes: make([]*Attribute, 0),
		Nodes:      make([]*Node, 0),
	}
}

func (n *Node) findAttribute(name Name) *Attribute {
	for _, attr := range n.Attributes {
		if attr.Name == name {
			return attr
		}
	}

	return nil
}

func (n *Node) findNodes(name Name) []*Node {
	var result []*Node

	for _, node := range n.Nodes {
		if node.Name == name {
			result = append(result, node)
		}
	}

	return result
}

func (n *Node) hasNode(name Name) bool {
	for _, node := range n.Nodes {
		if node.Name == name {
			return true
		}
	}

	return false
}

func (n *Node) ensureNode(name Name) *Node {
	for _, node := range n.Nodes {
		if node.Name == name {
			return node
		}
	}

	newNode := NewNode(name)
	n.Nodes = append(n.Nodes, newNode)
	return newNode
}

func (n *Node) setAttribute(name Name, value string) {
	attr := n.findAttribute(name)
	if attr == nil {
		attr = &Attribute{
			Name: name,
		}
		n.Attributes = append(n.Attributes, attr)
	}

	attr.Value = value
}

func (n *Node) getAttribute(name Name) string {
	attr := n.findAttribute(name)
	if attr != nil {
		return attr.Value
	}

	return ""
}

func (n *Node) removeAttribute(name Name) {
	var filtered []*Attribute

	for _, attr := range n.Attributes {
		if attr.Name != name {
			filtered = append(filtered, attr)
		}
	}

	n.Attributes = filtered
}

func (n *Node) removeNodes(name Name) {
	var filtered []*Node

	for _, node := range n.Nodes {
		if node.Name != name {
			filtered = append(filtered, node)
		}
	}

	n.Nodes = filtered
}

func (n *Node) addNode(node *Node) {
	n.Nodes = append(n.Nodes, node)
}
