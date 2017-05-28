package libvirtxml

type DomainDiskSource struct {
	node *Node
}

func newDomainDiskSource(node *Node) DomainDiskSource {
	return DomainDiskSource{
		node: node,
	}
}

func (s DomainDiskSource) File() string {
	return s.node.getAttribute(nameForLocal("file"))
}

func (s DomainDiskSource) SetFile(value string) {
	s.node.setAttribute(nameForLocal("file"), value)
}

func (s DomainDiskSource) Pool() string {
	return s.node.getAttribute(nameForLocal("pool"))
}

func (s DomainDiskSource) SetPool(value string) {
	s.node.setAttribute(nameForLocal("pool"), value)
}

func (s DomainDiskSource) Volume() string {
	return s.node.getAttribute(nameForLocal("volume"))
}

func (s DomainDiskSource) SetVolume(value string) {
	s.node.setAttribute(nameForLocal("volume"), value)
}

func (s DomainDiskSource) Mode() string {
	return s.node.getAttribute(nameForLocal("mode"))
}

func (s DomainDiskSource) SetMode(value string) {
	s.node.setAttribute(nameForLocal("mode"), value)
}
