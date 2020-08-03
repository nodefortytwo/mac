package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/urfave/cli/v2"
)

func GetCommand() *cli.Command {
	return &cli.Command{
		Name:  "docker",
		Usage: "some docker shortcuts",
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:  "clean",
				Usage: "free up some resources",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:   "procs",
						Usage:  "force kill all processes",
						Action: procHandler,
					},
					&cli.Command{
						Name:   "images",
						Usage:  "remove any dangling images",
						Action: imageHandler,
					},
				},
			},
		},
	}
}

func procHandler(c *cli.Context) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, container := range containers {
		log.Warnf("Stopping, Killing and Removing: %s", container.Image)
		to := time.Second * 5
		err = cli.ContainerStop(ctx, container.ID, &to)
		if err != nil {
			return err
		}

		err = cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func imageHandler(c *cli.Context) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	args := filters.NewArgs()
	args.Add("dangling", "true")

	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return err
	}

	for _, image := range images {
		log.Warnf("Removing %s", image.ID)
		cli.ImageRemove(ctx, image.ID, types.ImageRemoveOptions{})
	}

	return nil
}
