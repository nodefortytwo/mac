package prompt

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
	"os"
	"time"

	. "github.com/logrusorgru/aurora"
)

var statusMap = map[string]string{
	"*": "Your branch is ahead of",
	">": "renamed:",
	"!": "modified:",
	"+": "new file:",
	"x": "deleted:",
	"?": "Untracked files:",
}

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:      "prompt",
		Usage:     "makes a your prompt.. obviously",
		ArgsUsage: "",
		Action:    promptHandler,
	}
}

func promptHandler(c *cli.Context) error {
	t := time.Now().Format("15:04:05")
	path, _ := os.Getwd()
	gitStatus := printGitStatus()
	pointer := "->"
	fmt.Printf("%s: %s %s \n%s ", Bold(Green(t)), Bold(Cyan(path)), Bold(gitStatus), BrightWhite(pointer))

	return nil
}

func printGitStatus() string {
	repo, err := GetGitRepo()
	if repo == nil || err != nil {
		return ""
	}

	branch, err := getBranchName(repo)
	if repo == nil || err != nil {
		return ""
	}

	var flags []string

	status, err := status()
	if status == "" || err != nil {
		return ""
	}

	for flag, str := range statusMap {
		if strings.Contains(status, str) {
			flags = append(flags, flag)
		}
	}

	if len(flags) > 0 {
		branch += " " + strings.Join(flags, "")
	}

	return fmt.Sprintf("[%s]", branch)
}

func status() (string, error) {
	cmd := exec.Command("git", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := strings.TrimSpace(string(out))
		return "", errors.Wrap(err, errMsg)
	}
	return strings.TrimSpace(string(out)), nil
}

func GetGitRootPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func GetGitRepo() (*git.Repository, error) {
	path, err := GetGitRootPath()

	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	return repo, err
}

func getBranchName(repo *git.Repository) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	refName := ref.Name()
	if refName.IsBranch() {
		branchName := strings.TrimPrefix(refName.String(), "refs/heads/")
		return branchName, nil
	}

	return "", nil
}
