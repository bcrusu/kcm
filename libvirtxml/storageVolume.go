package libvirtxml

type StorageVolume struct {
	doc  *Document
	root *Node
}

func NewStorageVolume() StorageVolume {
	doc := &Document{}
	doc.Root = NewNode(nameForLocal("volume"))

	return StorageVolume{
		doc:  doc,
		root: doc.Root,
	}
}

func NewStorageVolumeForXML(xmlDoc string) (*StorageVolume, error) {
	doc := &Document{}
	if err := doc.Unmarshal(xmlDoc); err != nil {
		return nil, err
	}

	if doc.Root == nil {
		doc.Root = NewNode(nameForLocal("volume"))
	}

	return &StorageVolume{
		doc:  doc,
		root: doc.Root,
	}, nil
}

func (s StorageVolume) MarshalToXML() (string, error) {
	return s.doc.Marshal()
}

func (s StorageVolume) Type() string {
	return s.root.getAttribute(nameForLocal("type"))
}

func (s StorageVolume) SetType(value string) {
	s.root.setAttribute(nameForLocal("type"), value)
}

func (s StorageVolume) Name() string {
	return s.root.ensureNode(nameForLocal("name")).CharData
}

func (s StorageVolume) SetName(value string) {
	s.root.ensureNode(nameForLocal("name")).CharData = value
}

func (s StorageVolume) Key() string {
	return s.root.ensureNode(nameForLocal("key")).CharData
}

func (s StorageVolume) SetKey(value string) {
	s.root.ensureNode(nameForLocal("key")).CharData = value
}

func (s StorageVolume) Capacity() StorageVolumeSize {
	node := s.root.ensureNode(nameForLocal("capacity"))
	return newStorageVolumeSize(node)
}

func (s StorageVolume) Allocation() StorageVolumeSize {
	node := s.root.ensureNode(nameForLocal("allocation"))
	return newStorageVolumeSize(node)
}

func (s StorageVolume) Target() StorageVolumeTarget {
	node := s.root.ensureNode(nameForLocal("target"))
	return newStorageVolumeTarget(node)
}

func (s StorageVolume) BackingStore() StorageVolumeBackingStore {
	node := s.root.ensureNode(nameForLocal("backingStore"))
	return newStorageVolumeBackingStore(node)
}
