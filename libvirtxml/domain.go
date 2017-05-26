package libvirtxml

type Domain struct {
	doc  *Document
	root *Node
}

func NewDomainForXML(xmlDoc string) (Domain, error) {
	doc := &Document{}
	if err := doc.Unmarshal(xmlDoc); err != nil {
		return Domain{}, err
	}

	if doc.Root == nil {
		doc.Root = NewNode(nameForLocal("domain"))
	}

	return Domain{
		doc:  doc,
		root: doc.Root,
	}, nil
}

func (s Domain) MarshalToXML() (string, error) {
	return s.doc.Marshal()
}

func (s Domain) Name() string {
	return s.root.ensureNode(nameForLocal("name")).CharData
}

func (s Domain) SetName(value string) {
	s.root.ensureNode(nameForLocal("name")).CharData = value
}

func (s Domain) UUID() string {
	return s.root.ensureNode(nameForLocal("uuid")).CharData
}

func (s Domain) SetUUID(value string) {
	s.root.ensureNode(nameForLocal("uuid")).CharData = value
}

func (s Domain) ID() string {
	return s.root.getAttribute(nameForLocal("id"))
}

func (s Domain) SetID(value string) {
	s.root.setAttribute(nameForLocal("id"), value)
}

func (s Domain) Devices() DomainDevices {
	node := s.root.ensureNode(nameForLocal("devices"))
	return newDomainDevices(node)
}

func (s Domain) Metadata() DomainMetadata {
	node := s.root.ensureNode(nameForLocal("metadata"))
	return newDomainMetadata(node)
}
