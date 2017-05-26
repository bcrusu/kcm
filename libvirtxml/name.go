package libvirtxml

import "encoding/xml"

type Name struct {
	Local string
	Space string
}

func NewName(namespace string, local string) Name {
	return Name{
		Local: local,
		Space: namespace,
	}
}

func nameForLocal(local string) Name {
	return Name{
		Local: local,
	}
}

func nameForXMLName(name xml.Name) Name {
	return Name{
		Local: name.Local,
		Space: name.Space,
	}
}

func (n Name) toXMLName() xml.Name {
	return xml.Name{
		Local: n.Local,
		Space: n.Space,
	}
}
