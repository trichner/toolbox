package json2sheet

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/trichner/toolbox/pkg/json2sheet"
)

func TestExec(t *testing.T) {
	url, err := json2sheet.WriteToNewSheet(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
}
