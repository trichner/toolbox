package gitexec

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_gitClone(t *testing.T) {
	dir := t.TempDir()
	log.Println(dir)

	err := execGitClone(dir, "git@github.com:octokit/go-octokit.git")
	assert.NoError(t, err)
}

func Test_gitListTags(t *testing.T) {
	dir := t.TempDir()
	log.Println(dir)
	repo := "tb"
	err := execGitClone(dir, fmt.Sprintf("git@github.com:trichner/%s.git", repo))
	assert.NoError(t, err)

	tags, err := listRemoteTags(dir)
	assert.NoError(t, err)

	for _, b := range tags {
		fmt.Printf("- %s\n", b)
	}
}
