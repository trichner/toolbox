package json2sheet

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSheetWriter struct {
	invocations [][][]string
}

func (m *mockSheetWriter) UpdateValues(data [][]string) error {
	m.invocations = append(m.invocations, data)
	return nil
}

func TestWriteArraysTo(t *testing.T) {
	src := `
	["hello", "world"]
	["whats", "up", 55.88]
	[true, true, false]
	["wow"]
	`
	m := &mockSheetWriter{}
	WriteArraysTo(m, strings.NewReader(src))

	rows := m.invocations[0]
	assert.Equal(t, 4, len(rows))
	assert.Equal(t, 2, len(rows[0]))
}

func TestWriteObjectsTo(t *testing.T) {
	src := `
	{"a":"hello","b":"world"}
	{"b":2,"a":1,"c":3}
	{"d":4,"a":1,"c":3}
	`
	m := &mockSheetWriter{}
	WriteObjectsTo(m, strings.NewReader(src))

	rows := m.invocations[0]
	assert.Equal(t, 4, len(rows))
	assert.Equal(t, 4, len(rows[0]))
}
