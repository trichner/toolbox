package ast

import (
	"bytes"
	"encoding/json"
	"strconv"
)

//go:generate stringer -type=NodeType
type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeObject
	NodeTypeArray
	NodeTypeText
	NodeTypeNumber
	NodeTypeNull
	NodeTypeBoolean
)

type Node interface {
	json.Marshaler
	Type() NodeType
}

type Property struct {
	Name  string
	Value Node
}

type ObjectNode interface {
	Node
	Properties() []*Property
}

type ArrayNode interface {
	Node
	Items() []Node
}

type BooleanNode interface {
	Node
	Value() bool
}

type TextNode interface {
	Node
	Value() string
}

type NumberNode interface {
	Node
	Value() string
	ToInt() (int, error)
	ToInt64() (int64, error)
	ToFloat32() (float32, error)
	ToFloat64() (float64, error)
}

type NullNode interface {
	Node
}

func NewNullNode() NullNode {
	return NullNode(&nullNode{node: node{nodeType: NodeTypeNull}})
}

func NewBooleanNode(value bool) BooleanNode {
	return BooleanNode(&boolValueNode{
		node:  node{nodeType: NodeTypeBoolean},
		value: value,
	})
}

func NewTextNode(value string) TextNode {
	return TextNode(&textNode{
		node:  node{nodeType: NodeTypeText},
		value: value,
	})
}

func NewArrayNode(items []Node) ArrayNode {
	return ArrayNode(&arrayNode{
		node:  node{nodeType: NodeTypeArray},
		items: items,
	})
}

func NewNumberNode(value string) NumberNode {
	return NumberNode(&numberValueNode{
		node:  node{nodeType: NodeTypeNumber},
		value: value,
	})
}

func NewObjectNode(properties []*Property) ObjectNode {
	return ObjectNode(&objectNode{
		node:       node{nodeType: NodeTypeObject},
		properties: properties,
	})
}

type node struct {
	nodeType NodeType
}

func (n *node) Type() NodeType {
	return n.nodeType
}

type nullNode struct {
	node
}

func (n *nullNode) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

type textNode struct {
	node
	value string
}

func (s *textNode) Value() string {
	return s.value
}

func (s *textNode) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	writeEscaped(&buf, s.value)
	return buf.Bytes(), nil
}

type boolValueNode struct {
	node
	value bool
}

func (b *boolValueNode) Value() bool {
	return b.value
}

func (b *boolValueNode) MarshalJSON() ([]byte, error) {
	if b.value {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

type numberValueNode struct {
	node
	value string
}

func (n *numberValueNode) Value() string {
	return n.value
}

func (n *numberValueNode) ToInt() (int, error) {
	return strconv.Atoi(n.value)
}

func (n *numberValueNode) ToInt64() (int64, error) {
	return strconv.ParseInt(n.value, 10, 64)
}

func (n *numberValueNode) ToFloat32() (float32, error) {
	v, err := strconv.ParseFloat(n.value, 32)
	return float32(v), err
}

func (n *numberValueNode) ToFloat64() (float64, error) {
	return strconv.ParseFloat(n.value, 64)
}

func (n *numberValueNode) MarshalJSON() ([]byte, error) {
	return []byte(n.value), nil
}

type objectNode struct {
	node
	properties []*Property
}

func (o *objectNode) Properties() []*Property {
	return o.properties
}

func (o *objectNode) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.WriteRune('{')
	if err != nil {
		return nil, err
	}

	for i, p := range o.properties {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(`"`)
		buf.WriteString(p.Name)
		buf.WriteString(`":`)

		if p.Value == nil {
			buf.WriteString("null")
		} else {
			v, err := p.Value.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(v)
		}

	}

	buf.WriteRune('}')
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type arrayNode struct {
	node
	items []Node
}

func (a *arrayNode) Items() []Node {
	return a.items
}

func (a *arrayNode) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('[')

	for i, p := range a.items {
		if i != 0 {
			buf.WriteByte(',')
		}
		v, err := p.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(v)
	}

	buf.WriteRune(']')
	return buf.Bytes(), nil
}
