package git

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Branch struct {
	Name    string
	Current bool
}

func GetRepoRoot(output io.Writer) (string, error) {
	gitCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	gitCmd.Stderr = output
	out, err := gitCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git top-level")
	}
	return string(out), nil
}

func GetBranches() ([]Branch, error) {
	gitCmd := exec.Command("git", "branch")
	gitCmd.Stderr = os.Stderr
	output, err := gitCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git branches: %w", err)
	}
	var branches []Branch
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		branchText := scanner.Text()
		branches = append(branches, Branch{
			Name:    strings.TrimSpace(strings.TrimPrefix(branchText, "*")),
			Current: strings.HasPrefix(branchText, "*"),
		})
	}
	return branches, nil
}

func ChangeBranch(branch string, output io.Writer) error {
	gitCmd := exec.Command("git", "checkout", branch)
	gitCmd.Stdout = output
	gitCmd.Stderr = output
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout branch %s: %w", branch, err)
	}
	return nil
}

func DeleteBranch(branch string, output io.Writer) ([]Branch, error) {
	gitCmd := exec.Command("git", "branch", "-D", branch)
	gitCmd.Stdout = output
	gitCmd.Stderr = output
	if err := gitCmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to delete branch %s: %w", branch, err)
	}
	return GetBranches()
}
