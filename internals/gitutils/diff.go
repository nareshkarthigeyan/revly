package gitutils

import (
	"bytes"
	"os/exec"
)

func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD~1", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}