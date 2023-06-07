package csv2json

import (
	"log"
	"os"

	c2j "github.com/trichner/toolbox/pkg/csv2json"
)

func Exec(args []string) {
	err := c2j.Convert(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalf("cannot convert csv to json: %v", err)
	}
}
