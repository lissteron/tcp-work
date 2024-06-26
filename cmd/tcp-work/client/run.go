package client

import (
	"context"
	"fmt"
	"time"

	"github.com/lissteron/tcp-work/config"
	"github.com/lissteron/tcp-work/pkg/tlog"
)

func run(ctx context.Context, cfg *config.Config) error {
	logger, err := tlog.New()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}

	app := NewApp(ctx, logger, cfg)

	time.Sleep(time.Second)

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}
