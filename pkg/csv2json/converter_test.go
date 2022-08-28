package csv2json

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed cities.csv
var csvFixture string

//go:embed cities.ndjson
var ndjsonFixture string

func TestConvert(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Convert(strings.NewReader(csvFixture), buf)

	assert.NoError(t, err)
	assert.Equal(t, chomp(ndjsonFixture), chomp(buf.String()))
}

func chomp(s string) string {
	return strings.TrimSpace(s)
}
