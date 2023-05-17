package json2sheet

import (
	"fmt"
	"github.com/alecthomas/kong"
	"log"
	"os"
	"strings"

	"github.com/trichner/toolbox/pkg/json2sheet"
)

var cli struct {
	SpreadsheetUrl string `help:"complete URL to the spreadsheet"`
}

func Exec(args []string) {
	parser := kong.Must(&cli, kong.Name(args[0]))
	_, err := parser.Parse(args[1:])
	parser.FatalIfErrorf(err)

	spreadsheetUrl := strings.TrimSpace(cli.SpreadsheetUrl)
	if spreadsheetUrl != "" {
		url, err := json2sheet.UpdateSheet(spreadsheetUrl, os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(url)
	} else {
		url, err := json2sheet.WriteToNewSheet(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(url)
	}

}
