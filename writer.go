package main

import (
	"fmt"
	"strings"
)

const (
	HTRoot = "root"
	HTBody = "body"
	HTDiv  = "div"
	HTH1   = "h1"
	HTH2   = "h2"
	HTH3   = "h3"
	HTH4   = "h4"
	HTH5   = "h5"
	HTUL   = "ul"
	SPAN   = "span"
	HTLI   = "li"
	HTBR   = "br"
	HTA    = "a"
	TABLE  = "table"
	TR     = "tr"
	TD     = "td"
)

type ContentDelegate func(element *Element)
type ElementAttributes map[string]string

func (s ElementAttributes) Format() string {
	output := strings.Builder{}
	for key, value := range s {
		output.WriteString(fmt.Sprintf("%s=\"%s\"", key, value))
	}

	return output.String()
}

type Element struct {
	Tag        string
	Text       string
	Attributes ElementAttributes
	Children   []*Element
}

func NewElement(tag string) *Element {
	return &Element{
		Tag:        tag,
		Attributes: ElementAttributes{},
	}
}

func (s *Element) Element(tag string) *Element {
	element := NewElement(tag)
	s.Children = append(s.Children, element)

	return element
}

func (s *Element) Do(delegate ContentDelegate) {
	delegate(s)
}

func (s *Element) openingTag() string {
	tagContent := s.Tag
	if attrs := s.Attributes.Format(); len(attrs) > 0 {
		tagContent = fmt.Sprintf("%s %s", s.Tag, attrs)
	}

	return fmt.Sprintf("<%s>", tagContent)
}

func (s *Element) closingTag() string {
	return fmt.Sprintf("</%s>", s.Tag)
}

func (s *Element) Output(builder *strings.Builder) {
	builder.WriteString(s.openingTag())
	builder.WriteRune('\n')

	if len(s.Text) > 0 {
		builder.WriteString(s.Text)
		builder.WriteRune('\n')
	}

	for _, child := range s.Children {
		child.Output(builder)
	}

	builder.WriteString(s.closingTag())
	builder.WriteRune('\n')
}
