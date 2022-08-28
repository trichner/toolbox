package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/trichner/toolbox/pkg/jira"
	"github.com/trichner/toolbox/pkg/jira/credentials"
)

func main() {
	// read args
	var key string
	flag.StringVar(&key, "key", "", "key to search for")

	var fixVersion string
	flag.StringVar(&fixVersion, "fixVersion", "", "fixVersion to search for")

	var query string
	flag.StringVar(&query, "query", "", "search with given query")

	var label string
	flag.StringVar(&label, "label", "", "label to search for")
	flag.Parse()

	clientCredentials, err := credentials.FindCredentials()
	if err != nil {
		log.Fatal(err)
	}

	service, err := jira.NewJiraService(clientCredentials.Baseurl, clientCredentials.Username, clientCredentials.Token)
	if err != nil {
		panic(err)
	}

	if fixVersion != "" {
		issues, err := service.SearchInFixVersion(fixVersion)
		if err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(issues, "", " ")
		fmt.Printf("%s\n", s)
		return
	}

	if key != "" {
		issue, err := service.GetByKey(key)
		if err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(issue, "", " ")
		fmt.Printf("%s\n", s)
		return
	}

	if label != "" {
		issues, err := service.SearchByLabel(label)
		if err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(issues, "", " ")
		fmt.Printf("%s\n", s)
		return
	}

	if query != "" {
		issues, err := service.SearchByQuery(query)
		if err != nil {
			panic(err)
		}

		s, _ := json.MarshalIndent(issues, "", " ")
		fmt.Printf("%s\n", s)
		return
	}

	flag.Usage()
	os.Exit(-1)
}
