package lock

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
)

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:   "lock",
		Usage:  "lock your mac",
		Action: lockHandler,
	}
}

func lockHandler(c *cli.Context) error {
	cmd := exec.Command(`/System/Library/CoreServices/Menu Extras/User.menu/Contents/Resources/CGSession`, "-suspend")
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := strings.TrimSpace(string(out))
		log.Fatal(errors.Wrap(err, errMsg))
	}
	return nil
}
