package json2sheet

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteToNewSheet(t *testing.T) {
	buf := strings.NewReader(`
	{"a":"hello","b":"world"}
	{"b":2,"a":1,"c":3}
	{"d":4,"a":1,"c":3}
	`)
	url, err := WriteToNewSheet(buf)
	fmt.Println(url)
	assert.NoError(t, err)
}
