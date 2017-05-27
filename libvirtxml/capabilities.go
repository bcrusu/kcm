package libvirtxml

type Capabilities struct {
	doc  *Document
	root *Node
}

func NewCapabilitiesForXML(xmlDoc string) (*Capabilities, error) {
	doc := &Document{}
	if err := doc.Unmarshal(xmlDoc); err != nil {
		return nil, err
	}

	if doc.Root == nil {
		doc.Root = NewNode(nameForLocal("capabilities"))
	}

	return &Capabilities{
		doc:  doc,
		root: doc.Root,
	}, nil
}

func (s Capabilities) Host() CapabilitiesHost {
	node := s.root.ensureNode(nameForLocal("host"))
	return newCapabilitiesHost(node)
}

func (s Capabilities) Guests() []CapabilitiesGuest {
	var result []CapabilitiesGuest

	nodes := s.root.findNodes(nameForLocal("guest"))
	for _, node := range nodes {
		result = append(result, newCapabilitiesGuest(node))
	}

	return result
}
