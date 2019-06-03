package main

import (
	"fmt"
	"strings"
)

type HTMLTag string

func (s HTMLTag) EmitNewlines() bool {
	for _, tagValue := range nonNewlineTags {
		if tagValue == s {
			return false
		}
	}

	return true
}

func (s HTMLTag) String() string {
	return string(s)
}

const (
	HTP           HTMLTag = "p"
	HTRoot        HTMLTag = "root"
	HTBody        HTMLTag = "body"
	Division      HTMLTag = "div"
	H1            HTMLTag = "h1"
	H2            HTMLTag = "h2"
	H3            HTMLTag = "h3"
	H4            HTMLTag = "h4"
	H5            HTMLTag = "h5"
	UnorderedList HTMLTag = "ul"
	Span          HTMLTag = "span"
	ListItem      HTMLTag = "li"
	BR            HTMLTag = "br"
	Anchor        HTMLTag = "a"
	Table         HTMLTag = "table"
	TableHeaders  HTMLTag = "thead"
	TableRow      HTMLTag = "tr"
	TableCell     HTMLTag = "td"
)

var nonNewlineTags = []HTMLTag{
	Span,
	Anchor,
}

type ContentDelegate func(element *DocumentElement)
type ElementAttributes map[string]string

func (s ElementAttributes) Format() string {
	output := strings.Builder{}
	for key, value := range s {
		output.WriteString(fmt.Sprintf("%s=\"%s\"", key, value))
	}

	return output.String()
}

type DocumentElement struct {
	Tag        HTMLTag
	Text       string
	Attributes ElementAttributes
	Children   []*DocumentElement
}

func Element(tag HTMLTag) *DocumentElement {
	return &DocumentElement{
		Tag:        tag,
		Attributes: ElementAttributes{},
	}
}

func NewElement(tag HTMLTag, text string, attrs ElementAttributes) *DocumentElement {
	return &DocumentElement{
		Tag:        tag,
		Text:       text,
		Attributes: attrs,
	}
}

func (s *DocumentElement) Push(child *DocumentElement) *DocumentElement {
	s.Children = append(s.Children, child)
	return s
}

func (s *DocumentElement) Element(tag HTMLTag) *DocumentElement {
	element := Element(tag)
	s.Children = append(s.Children, element)

	return element
}

func (s *DocumentElement) Do(delegate ContentDelegate) {
	delegate(s)
}

func (s *DocumentElement) openingTag() string {
	tagContent := s.Tag.String()
	if attrs := s.Attributes.Format(); len(attrs) > 0 {
		tagContent = fmt.Sprintf("%s %s", s.Tag, attrs)
	}

	return fmt.Sprintf("<%s>", tagContent)
}

func (s *DocumentElement) closingTag() string {
	return fmt.Sprintf("</%s>", s.Tag)
}

func (s *DocumentElement) Output(builder *strings.Builder) {
	emitNewlines := s.Tag.EmitNewlines()

	if builder.WriteString(s.openingTag()); emitNewlines {
		builder.WriteRune('\n')
	}

	if len(s.Text) > 0 {
		if builder.WriteString(s.Text); emitNewlines {
			builder.WriteRune('\n')
		}
	}

	for _, child := range s.Children {
		child.Output(builder)
	}

	if builder.WriteString(s.closingTag()); emitNewlines {
		builder.WriteRune('\n')
	}
}

func (s *DocumentElement) String() string {
	builder := &strings.Builder{}
	s.Output(builder)

	return builder.String()
}
