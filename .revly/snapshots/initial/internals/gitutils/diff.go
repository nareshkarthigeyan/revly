package gitutils

import (
	"bytes"
	"os"
	"os/exec"
)

func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD~1", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

// GetStagedDiff returns the diff of staged changes.
func GetStagedDiff() ([]byte, error) {
	return exec.Command("git", "diff", "--staged").Output()
}

// GetWorkingDiff returns the diff of unstaged changes.
func GetWorkingDiff() ([]byte, error) {
	return exec.Command("git", "diff").Output()
}

func CreateSnapshot(path string) error {
	// Create the directory if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	// Use rsync to copy the files
	cmd := exec.Command("rsync", "-a", "--exclude", ".git", "--exclude", ".revly", ".", path)
	return cmd.Run()
}

func DiffSnapshots(path1, path2 string) ([]byte, error) {
	cmd := exec.Command("git", "diff", "--no-index", path1, path2)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit code 1 means there are differences, which is not an error for us.
			if exitError.ExitCode() == 1 {
				return output, nil
			}
		}
		return nil, err
	}
	return output, nil
}