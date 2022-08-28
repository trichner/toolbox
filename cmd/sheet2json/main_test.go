package sheet2json

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlToSpreadsheet(t *testing.T) {
	u := "https://docs.google.com/spreadsheets/d/1dAN8MO9NDVPqVIoOxC9H_j4Ir5c1viQ97igGdXOyXsU/edit#gid=886605725"
	spreadsheetId, sheetId, err := urlToSpreadsheetID(u)

	assert.NoError(t, err)
	fmt.Printf("%s / %d", spreadsheetId, sheetId)
}

func TestExec(t *testing.T) {
	Exec([]string{"sheet2json", "--spreadsheet-id=https://docs.google.com/spreadsheets/d/1dAN8MO9NDVPqVIoOxC9H_j4Ir5c1viQ97igGdXOyXsU/edit#gid=886605725"})
}
