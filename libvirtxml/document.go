package libvirtxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

type Document struct {
	ProcInst *xml.ProcInst
	Root     *Node
	CharData string
	Comments string
}

func (d *Document) Unmarshal(xmlDoc string) error {
	reader := strings.NewReader(xmlDoc)
	decoder := xml.NewDecoder(reader)

	nodes, charData, comments, procInst, err := decodeNodes(decoder)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return errors.New("libvirtxml: invalid XML document - no root element")
	}

	if len(nodes) > 1 {
		return errors.New("libvirtxml: invalid XML document - conatains multiple root elements")
	}

	d.ProcInst = procInst
	d.Root = nodes[0]
	d.CharData = charData
	d.Comments = comments

	return nil
}

func (d *Document) Marshal() (string, error) {
	var buffer bytes.Buffer

	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", "  ")

	if d.ProcInst != nil {
		if err := encoder.EncodeToken(d.ProcInst); err != nil {
			return "", err
		}
	}

	if d.Comments != "" {
		if err := encoder.EncodeToken(xml.Comment(d.Comments)); err != nil {
			return "", err
		}
	}

	if d.CharData != "" {
		if err := encoder.EncodeToken(xml.CharData(d.CharData)); err != nil {
			return "", err
		}
	}

	if err := encodeNode(encoder, d.Root); err != nil {
		return "", err
	}

	encoder.Flush()
	return buffer.String(), nil
}

func decodeNodes(decoder *xml.Decoder) ([]*Node, string, string, *xml.ProcInst, error) {
	var nodes []*Node
	var charData string
	var comments string
	var procInst *xml.ProcInst
	var err error

loop:
	for {
		var token xml.Token
		token, err = decoder.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			var node *Node
			node, err = decodeNode(decoder, t)
			if err != nil {
				break loop
			}

			nodes = append(nodes, node)
		case xml.EndElement:
			break loop
		case xml.CharData:
			str := string(t)
			str = strings.TrimSpace(str)
			if str != "" {
				charData = strings.Join([]string{charData, str}, "")
			}
		case xml.Comment:
			comments = strings.Join([]string{comments, string(t)}, "")
		case xml.ProcInst:
			procInst = &t
		}
	}

	if err != nil && err != io.EOF {
		return nil, "", "", nil, err
	}

	return nodes, charData, comments, procInst, nil
}

func decodeNode(decoder *xml.Decoder, element xml.StartElement) (*Node, error) {
	attributes := make([]*Attribute, len(element.Attr), len(element.Attr))

	for i, attr := range element.Attr {
		attributes[i] = &Attribute{
			Name:  nameForXMLName(attr.Name),
			Value: attr.Value,
		}
	}

	result := NewNode(nameForXMLName(element.Name))
	result.Attributes = attributes

	nodes, charData, comments, _, err := decodeNodes(decoder)
	if err != nil {
		return nil, err
	}

	result.Nodes = nodes
	result.CharData = charData
	result.Comments = comments

	return result, nil
}

func encodeNode(encoder *xml.Encoder, node *Node) error {
	startElement := xml.StartElement{
		Name: node.Name.toXMLName(),
	}

	for _, attribute := range node.Attributes {
		if attribute.Value == "" {
			continue
		}

		startElement.Attr = append(startElement.Attr, xml.Attr{
			Name:  attribute.Name.toXMLName(),
			Value: attribute.Value,
		})
	}

	if err := encoder.EncodeToken(startElement); err != nil {
		return err
	}

	if node.Comments != "" {
		if err := encoder.EncodeToken(xml.Comment(node.Comments)); err != nil {
			return err
		}
	}

	if node.CharData != "" {
		if err := encoder.EncodeToken(xml.CharData(node.CharData)); err != nil {
			return err
		}
	}

	for _, node := range node.Nodes {
		if err := encodeNode(encoder, node); err != nil {
			return err
		}
	}

	endElement := xml.EndElement{
		Name: node.Name.toXMLName(),
	}

	if err := encoder.EncodeToken(endElement); err != nil {
		return err
	}

	return nil
}
