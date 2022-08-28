// Code generated by "stringer -type=NodeType"; DO NOT EDIT.

package ast

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NodeTypeUnknown-0]
	_ = x[NodeTypeObject-1]
	_ = x[NodeTypeArray-2]
	_ = x[NodeTypeText-3]
	_ = x[NodeTypeNumber-4]
	_ = x[NodeTypeNull-5]
	_ = x[NodeTypeBoolean-6]
}

const _NodeType_name = "NodeTypeUnknownNodeTypeObjectNodeTypeArrayNodeTypeTextNodeTypeNumberNodeTypeNullNodeTypeBoolean"

var _NodeType_index = [...]uint8{0, 15, 29, 42, 54, 68, 80, 95}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
