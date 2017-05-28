package libvirtxml

type DomainDevices struct {
	node *Node
}

func newDomainDevices(node *Node) DomainDevices {
	return DomainDevices{
		node: node,
	}
}

func (s DomainDevices) Emulator() string {
	return s.node.ensureNode(nameForLocal("emulator")).CharData
}

func (s DomainDevices) SetEmulator(value string) {
	s.node.ensureNode(nameForLocal("emulator")).CharData = value
}

func (s DomainDevices) Disks() []DomainDisk {
	var result []DomainDisk

	nodes := s.node.findNodes(nameForLocal("disk"))
	for _, node := range nodes {
		result = append(result, newDomainDisk(node))
	}

	return result
}

func (s DomainDevices) NewDisk() DomainDisk {
	node := NewNode(nameForLocal("disk"))
	s.node.addNode(node)
	return newDomainDisk(node)
}

func (s DomainDevices) SetDisks(disks []DomainDisk) {
	s.node.removeNodes(nameForLocal("disk"))

	for _, disk := range disks {
		s.node.addNode(disk.node)
	}
}

func (s DomainDevices) Graphics() []DomainGraphic {
	var result []DomainGraphic

	nodes := s.node.findNodes(nameForLocal("graphics"))
	for _, node := range nodes {
		result = append(result, newDomainGraphic(node))
	}

	return result
}

func (s DomainDevices) Interfaces() []DomainInterface {
	var result []DomainInterface

	nodes := s.node.findNodes(nameForLocal("interface"))
	for _, node := range nodes {
		result = append(result, newDomainInterface(node))
	}

	return result
}

func (s DomainDevices) NewInterface() DomainInterface {
	node := NewNode(nameForLocal("interface"))
	s.node.addNode(node)
	return newDomainInterface(node)
}

func (s DomainDevices) Channels() []DomainChannel {
	var result []DomainChannel

	nodes := s.node.findNodes(nameForLocal("channel"))
	for _, node := range nodes {
		result = append(result, newDomainChannel(node))
	}

	return result
}

func (s DomainDevices) Filesystems() []DomainFilesystem {
	var result []DomainFilesystem

	nodes := s.node.findNodes(nameForLocal("filesystem"))
	for _, node := range nodes {
		result = append(result, newDomainFilesystem(node))
	}

	return result
}

func (s DomainDevices) NewFilesystem() DomainFilesystem {
	node := NewNode(nameForLocal("filesystem"))
	s.node.addNode(node)
	return newDomainFilesystem(node)
}
