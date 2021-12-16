package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Branch struct {
	Name    string
	Current bool
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

func ChangeBranch(branch string) (string, error) {
	gitCmd := exec.Command("git", "checkout", branch)
	output, err := gitCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to checkout branch %s: %w", branch, err)
	}
	return string(output), nil
}
