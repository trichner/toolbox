package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Credentials struct {
	Baseurl  string `json:"baseurl"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func FindCredentials() (*Credentials, error) {
	var err error
	var contents []byte
	c, err := readFromEnvironment("JIRA_CREDENTIALS_PATH")
	if err != nil {
		contents = c
	}

	if contents == nil {
		contents, err = readFromHomeDirectory(".config/jira/credentials.json")
	}

	if contents == nil {
		c, err := readFile("jira-credentials.json")
		if err != nil {
			contents = c
		}
	}

	if contents == nil {
		return nil, fmt.Errorf("no Credentials found")
	}

	var clientCredentials Credentials
	if err := json.Unmarshal(contents, &clientCredentials); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &clientCredentials, nil
}

func readFromHomeDirectory(relativePath string) ([]byte, error) {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		return nil, fmt.Errorf("environment variable 'HOME' not defined")
	}

	p := path.Join(homePath, relativePath)
	return ioutil.ReadFile(p)
}

func readFromEnvironment(name string) ([]byte, error) {
	p := os.Getenv(name)
	if p == "" {
		return nil, fmt.Errorf("environment variable %q not defined", name)
	}
	return readFile(p)
}

func readFile(path string) ([]byte, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return c, nil
}
