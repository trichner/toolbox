package gitexec

import (
	"bytes"
	"os"
	"os/exec"
)

func execGitClone(vcsRoot, sshUrl string) error {
	cmd := exec.Command("git", "clone", sshUrl, ".")
	cmd.Dir = vcsRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}

func execGitBranchRemoteMerged(clonePath, ref string) (string, error) {
	return execute(clonePath, "git", "branch", "-r", "--merged", ref)
}

func execGitLogInRange(clonePath, commitRange string) (string, error) {
	// git log develop 'origin/release/v4.13..origin/release/v4.14' --format=oneline
	return execute(clonePath, "git", "log", commitRange, "--format=format:%H (%aE) (%P) (%D) %s")
}

func execGitLsRemoteTags(clonePath string) (string, error) {
	return execute(clonePath, "git", "ls-remote", "--tags", "--refs")
}

func execute(workingDir string, prog string, args ...string) (string, error) {
	cmd := exec.Command(prog, args...)
	cmd.Dir = workingDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
