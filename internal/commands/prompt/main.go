package prompt

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/nodefortytwo/isgit"
	"github.com/nodefortytwo/slice"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

var statusMap = map[string]git.StatusCode{
	"*":  git.UpdatedButUnmerged,
	">":  git.Renamed,
	"!":  git.Modified,
	"+":  git.Added,
	"x":  git.Deleted,
	"?":  git.Untracked,
	"||": git.Copied,
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
	fmt.Printf("%s: %s %s \n%s ", Bold(Green(t)), Bold(Cyan(path)), Bold(gitStatus), pointer)

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

	ref, err := repo.Worktree()
	if err != nil {
		return ""
	}
	stat, err := ref.Status()
	if err != nil {
		return ""
	}

	for _, status := range stat {
		for flag, statusCode := range statusMap {
			if status.Worktree == statusCode {
				flags = append(flags, flag)
			}
		}
	}

	flags = slice.String(flags).Unique()
	if len(flags) > 0 {
		branch += " " + strings.Join(flags, "")
	}

	return fmt.Sprintf("[%s]", branch)
}

func GetGitRepo() (*git.Repository, error) {
	path, err := isgit.GetRootDirWD()

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
