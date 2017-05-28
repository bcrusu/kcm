package libvirtxml

type DomainFilesystem struct {
	node *Node
}

func newDomainFilesystem(node *Node) DomainFilesystem {
	return DomainFilesystem{
		node: node,
	}
}

func (s DomainFilesystem) Type() string {
	return s.node.getAttribute(nameForLocal("type"))
}

func (s DomainFilesystem) SetType(value string) {
	s.node.setAttribute(nameForLocal("type"), value)
}

func (s DomainFilesystem) Accessmode() string {
	return s.node.getAttribute(nameForLocal("accessmode"))
}

func (s DomainFilesystem) SetAccessmode(value string) {
	s.node.setAttribute(nameForLocal("accessmode"), value)
}

func (s DomainFilesystem) SourceDir() string {
	node := s.node.ensureNode(nameForLocal("source"))
	return node.getAttribute(nameForLocal("dir"))
}

func (s DomainFilesystem) SetSourceDir(value string) {
	node := s.node.ensureNode(nameForLocal("source"))
	node.setAttribute(nameForLocal("dir"), value)
}

func (s DomainFilesystem) TargetDir() string {
	node := s.node.ensureNode(nameForLocal("target"))
	return node.getAttribute(nameForLocal("dir"))
}

func (s DomainFilesystem) SetTargetDir(value string) {
	node := s.node.ensureNode(nameForLocal("target"))
	node.setAttribute(nameForLocal("dir"), value)
}

func (s DomainFilesystem) Readonly() bool {
	return s.node.hasNode(nameForLocal("readonly"))
}

func (s DomainFilesystem) SetReadonly(value bool) {
	if value {
		s.node.ensureNode(nameForLocal("readonly"))
	} else {
		s.node.removeNodes(nameForLocal("readonly"))
	}
}
