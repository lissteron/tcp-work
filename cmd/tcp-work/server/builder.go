package server

import (
	"github.com/urfave/cli/v2"

	"github.com/lissteron/tcp-work/config"
)

func BuildCmd() *cli.Command {
	cfg := config.NewConfig()

	return &cli.Command{
		Name:        "server",
		Description: "proof of work server",
		Action: func(ctx *cli.Context) error {
			return run(ctx.Context, cfg)
		},
	}
}
