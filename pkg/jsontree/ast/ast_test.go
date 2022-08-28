package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshallJSON(t *testing.T) {
	n := NewObjectNode([]*Property{{Name: "a", Value: NewTextNode("hello")}, {Name: "b", Value: NewNullNode()}, {
		Name: "c", Value: NewArrayNode([]Node{NewNumberNode("2"), NewNumberNode("4"), NewNumberNode("8")}),
	}})

	txt, _ := n.MarshalJSON()

	assert.Equal(t, `{"a":"hello","b":null,"c":[2,4,8]}`, string(txt))
}
