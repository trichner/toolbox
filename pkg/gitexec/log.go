package gitexec

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

var logLinePattern = regexp.MustCompile("^([0-9a-f]{40}) \\(([^)]*)\\) \\(([^)]*)\\) \\(([^)]*)\\) (.*)$")

type Commit struct {
	Hash         string
	AuthorEmail  string
	ParentHashes []string
	Refs         []string
	Message      string
}

func parseLogs(data string) ([]Commit, error) {
	var commits []Commit
	s := bufio.NewScanner(strings.NewReader(data))
	for s.Scan() {
		line := s.Text()
		commit, err := parseLogLine(line)
		if err != nil {
			return nil, err
		}
		commits = append(commits, commit)
	}
	return commits, nil
}

func parseLogLine(line string) (Commit, error) {
	matches := logLinePattern.FindStringSubmatch(line)
	if matches == nil {
		return Commit{}, fmt.Errorf("cannot parse line: '%s'", line)
	}
	return Commit{
		Hash:         matches[1],
		AuthorEmail:  matches[2],
		ParentHashes: strings.Split(matches[3], " "),
		Refs:         strings.Split(matches[4], ","),
		Message:      matches[5],
	}, nil
}
