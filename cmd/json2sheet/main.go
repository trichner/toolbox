package json2sheet

import (
	"fmt"
	"log"
	"os"

	"github.com/trichner/toolbox/pkg/json2sheet"
)

func Exec(args []string) {
	url, err := json2sheet.WriteToNewSheet(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
}
