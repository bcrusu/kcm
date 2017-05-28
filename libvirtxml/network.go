package libvirtxml

type Network struct {
	doc  *Document
	root *Node
}

func NewNetwork() Network {
	doc := &Document{}
	doc.Root = NewNode(nameForLocal("network"))

	return Network{
		doc:  doc,
		root: doc.Root,
	}
}

func NewNetworkForXML(xmlDoc string) (*Network, error) {
	doc := &Document{}
	if err := doc.Unmarshal(xmlDoc); err != nil {
		return nil, err
	}

	if doc.Root == nil {
		doc.Root = NewNode(nameForLocal("network"))
	}

	return &Network{
		doc:  doc,
		root: doc.Root,
	}, nil
}

func (s Network) MarshalToXML() (string, error) {
	return s.doc.Marshal()
}

func (s Network) Name() string {
	return s.root.ensureNode(nameForLocal("name")).CharData
}

func (s Network) SetName(value string) {
	s.root.ensureNode(nameForLocal("name")).CharData = value
}

func (s Network) UUID() string {
	return s.root.ensureNode(nameForLocal("uuid")).CharData
}

func (s Network) SetUUID(value string) {
	s.root.ensureNode(nameForLocal("uuid")).CharData = value
}

func (s Network) MACAddress() string {
	node := s.root.ensureNode(nameForLocal("mac"))
	return node.getAttribute(nameForLocal("address"))
}

func (s Network) Forward() NetworkForward {
	node := s.root.ensureNode(nameForLocal("forward"))
	return newNetworkForward(node)
}

func (s Network) Bridge() NetworkBridge {
	node := s.root.ensureNode(nameForLocal("bridge"))
	return newNetworkBridge(node)
}

func (s Network) IPs() []NetworkIP {
	var result []NetworkIP

	nodes := s.root.findNodes(nameForLocal("ip"))
	for _, node := range nodes {
		result = append(result, newNetworkIP(node))
	}

	return result
}

func (s Network) SetIPs(ips []NetworkIP) {
	s.root.removeNodes(nameForLocal("ip"))

	for _, ip := range ips {
		s.root.addNode(ip.node)
	}
}

func (s Network) NewIP() NetworkIP {
	node := NewNode(nameForLocal("ip"))
	s.root.addNode(node)
	return newNetworkIP(node)
}

func (s Network) Metadata() Metadata {
	node := s.root.ensureNode(nameForLocal("metadata"))
	return newMetadata(node)
}
