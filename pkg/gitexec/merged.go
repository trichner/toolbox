package gitexec

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-git/go-git/v5"
)

type Repository struct {
	clonePath  string
	repository *git.Repository
}

func Clone(sshUrl string) (*Repository, error) {
	dir, err := createTempDir()
	if err != nil {
		return nil, err
	}

	err = execGitClone(dir, sshUrl)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	return &Repository{clonePath: dir, repository: repo}, nil
}

func CloneTreeless(sshUrl string) (*Repository, error) {
	dir, err := createTempDir()
	if err != nil {
		return nil, err
	}

	err = execGitClone(dir, sshUrl)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	return &Repository{clonePath: dir, repository: repo}, nil
}

func (r *Repository) ListCommitsInRange(fromRef, toRef string) ([]Commit, error) {
	if fromRef == toRef {
		return []Commit{}, nil
	}

	var commitRange string
	if fromRef == "" {
		commitRange = toRef
	} else {
		commitRange = fmt.Sprintf("%s..%s", fromRef, toRef)
	}

	log, err := execGitLogInRange(r.clonePath, commitRange)
	if err != nil {
		return nil, err
	}
	return parseLogs(log)
}

// ListMergedBranches lists all remote branches that were merged into a given ref (branch, tag, sha, ...)
func (r *Repository) ListMergedBranches(ref string) ([]string, error) {
	return listMergedBranches(r.clonePath, ref)
}

// ListTags lists all remote tags of a repository
func (r *Repository) ListTags() ([]string, error) {
	return listRemoteTags(r.clonePath)
}

func createTempDir() (string, error) {
	return ioutil.TempDir("", "gitclone")
}

func listMergedBranches(clonePath, ref string) ([]string, error) {
	output, err := execGitBranchRemoteMerged(clonePath, ref)
	if err != nil {
		return nil, fmt.Errorf("cannot list merged branches for ref '%s': %w", ref, err)
	}

	splits := strings.Split(output, "\n")

	trimmed := make([]string, 0, len(splits))
	for _, s := range splits {
		t := strings.TrimSpace(s)
		if t == "" {
			continue
		}
		trimmed = append(trimmed, t)
	}
	return trimmed, nil
}

func listRemoteTags(clonePath string) ([]string, error) {
	output, err := execGitLsRemoteTags(clonePath)
	if err != nil {
		return nil, fmt.Errorf("cannot list remote tags: %w", err)
	}

	splits := strings.Split(output, "\n")

	tags := make([]string, 0, len(splits))
	// looks like: c3c4a34a7cb59dc04ab2d7112fe7640529f9fdbb	Refs/tags/mailer-api/v1.0.0
	for _, s := range splits {
		if s == "" {
			continue
		}
		fields := strings.Fields(s)
		if len(fields) != 2 {
			return nil, fmt.Errorf("listed remote tags have bad format: '%s'", s)
		}
		ref := strings.TrimSpace(fields[1])
		tag := strings.TrimPrefix(ref, "Refs/tags/")
		tags = append(tags, tag)
	}

	return tags, nil
}
