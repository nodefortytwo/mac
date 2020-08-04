package quit

import (
	"github.com/keybase/go-ps"
	"github.com/nodefortytwo/mac/internal/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os/exec"
	"strings"
	"syscall"
)

const (
	SIGTERM = 15
	SIGKILL = 9
)

func GetCommand() *cli.Command {

	return &cli.Command{
		Name:   "quit",
		Usage:  "kills all applications and closes all finder windows",
		Action: quitCommand,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "exclude",
				Aliases: []string{"e"},
				Usage:   "apps to exclude",
				Value:   cli.NewStringSlice(config.New().GetQuitExclusions()...),
			},
			&cli.BoolFlag{
				Name:    "hard",
				Aliases: []string{"f"},
				Usage:   "use SIGKILL instead of SIGTERM",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  "dryrun",
				Usage: "just print what you would brutally murder",
				Value: false,
			},
		},
	}
}

func quitCommand(c *cli.Context) error {

	procs, err := ps.Processes()

	for _, p := range procs {
		path, _ := p.Path()
		if p.PPid() != 1 {
			continue
		}
		if !strings.Contains(path, "/Applications") {
			continue
		}

		if shouldExclude(path, c.StringSlice("exclude")) {
			continue
		}

		signal := syscall.SIGTERM
		if c.Bool("hard") {
			signal = syscall.SIGKILL
		}

		log.Warnf("Killing: %s", path)

		if c.Bool("dryrun") {
			return nil
		}

		syscall.Kill(p.Pid(), signal)
	}

	cmd := exec.Command("osascript", "-e", `tell application "Finder" to close windows`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := strings.TrimSpace(string(out))
		log.Fatal(errors.Wrap(err, errMsg))
	}

	return nil
}

func shouldExclude(path string, exclusions []string) bool {
	for _, term := range exclusions {
		if strings.Contains(strings.ToLower(path), strings.ToLower(term)) {
			return true
		}
	}
	return false
}
