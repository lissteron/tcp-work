package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"github.com/lissteron/tcp-work/cmd/tcp-work/client"
	"github.com/lissteron/tcp-work/cmd/tcp-work/server"
	"github.com/lissteron/tcp-work/config"
	"github.com/lissteron/tcp-work/pkg/tlog"
)

const (
	exitCodeOk = iota
	exitCodeAppError
	exitCodeLoggerError
	exitCodePanic
)

func main() {
	exitCode := exitCodeOk
	defer func() { os.Exit(exitCode) }()

	logger, err := tlog.New()
	if err != nil {
		exitCode = exitCodeLoggerError

		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Errorw(context.Background(), "service panic", "error", r, "stack_trace", string(debug.Stack()))

			exitCode = exitCodePanic
		}
	}()

	newApp := &cli.App{
		Commands: []*cli.Command{
			server.BuildCmd(),
			client.BuildCmd(),
			ver(),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := newApp.RunContext(ctx, os.Args); err != nil {
		logger.Errorw(ctx, "run failed", "error", err)

		exitCode = exitCodeAppError
	}
}

func ver() *cli.Command {
	return &cli.Command{
		Name:        "version",
		Aliases:     []string{"ver", "v"},
		Description: "Show build info",
		Action: func(_ *cli.Context) error {
			log.Printf("ServiceName: %s", config.ServiceName)
			log.Printf("AppName: %s", config.AppName)
			log.Printf("GitHash: %s", config.GitHash)
			log.Printf("Version: %s", config.Version)
			log.Printf("BuildAt: %s", config.BuildAt)
			log.Printf("ReleaseID: %s", config.ReleaseID)

			return nil
		},
	}
}
