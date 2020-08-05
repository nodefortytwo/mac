package clone

import (
	"github.com/nodefortytwo/mac/internal/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func GetCommand() *cli.Command {
	root := config.New().GetCodeRoot()
	return &cli.Command{
		Name:      "clone",
		Usage:     "Clone a repo from a url",
		UsageText: "Automatically download a repo to {CodeRoot}/{RepoPath}\n\t eg. https://github.com/nodefortytwo/mac -> " + root + "/nodefortytwo/mac",
		ArgsUsage: "[git repo url]",
		Action:    cloneHandler,
	}
}

func cloneHandler(c *cli.Context) error {
	var gitRepo string
	if gitRepo = c.Args().First(); gitRepo == "" {
		return errors.New("provide git repo as first arg")
	}

	if !strings.HasSuffix(gitRepo, ".git") {
		gitRepo += ".git"
	}

	repo, err := url.Parse(gitRepo)
	if err != nil {
		return err
	}

	if repo.Host != "github.com" {
		return errors.Errorf("this command only supports github repos, %s not supported", repo.Host)
	}

	path := getPath(repo.Path)

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return errors.Wrap(err, "error creating directory")
	}

	if _, err := os.Stat(path + "/.git"); !os.IsNotExist(err) {
		log.Info("already cloned")
		return nil
	}

	cmd := exec.Command("git", "clone", gitRepo, path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := strings.TrimSpace(string(out))
		log.Fatal(errors.Wrap(err, errMsg))
	}
	return nil
}

func getPath(repo string) string {
	path := strings.ReplaceAll(repo, ".git", "")
	return config.New().GetCodeRoot() + path
}
