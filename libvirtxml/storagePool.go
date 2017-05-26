package libvirtxml

type StoragePool struct {
	doc  *Document
	root *Node
}

func NewStoragePool() StoragePool {
	doc := &Document{}
	doc.Root = NewNode(nameForLocal("pool"))

	return StoragePool{
		doc:  doc,
		root: doc.Root,
	}
}

func NewStoragePoolForXML(xmlDoc string) (*StoragePool, error) {
	doc := &Document{}
	if err := doc.Unmarshal(xmlDoc); err != nil {
		return nil, err
	}

	if doc.Root == nil {
		doc.Root = NewNode(nameForLocal("pool"))
	}

	return &StoragePool{
		doc:  doc,
		root: doc.Root,
	}, nil
}

func (s StoragePool) MarshalToXML() (string, error) {
	return s.doc.Marshal()
}

func (s StoragePool) Type() string {
	return s.root.getAttribute(nameForLocal("type"))
}

func (s StoragePool) SetType(value string) {
	s.root.setAttribute(nameForLocal("type"), value)
}

func (s StoragePool) Target() StoragePoolTarget {
	node := s.root.ensureNode(nameForLocal("target"))
	return newStoragePoolTarget(node)
}
