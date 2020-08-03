package main

import (
	"context"
	"github.com/nodefortytwo/mac/internal/commands/clone"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
)

func main() {

	app := &cli.App{
		Name:  "mac",
		Usage: "your friendly utility CLI",
		Commands: []*cli.Command{
			clone.GetCommand(),
		},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.mac")

	err := viper.ReadInConfig() // Find and read the config file

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("no config file found at ~/.mac/config")
		} else {
			log.Fatal(errors.Wrap(err, "error with config parsing"))
		}
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err = app.RunContext(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
