package server

import (
	"context"
	"fmt"

	"github.com/lissteron/tcp-work/config"
	"github.com/lissteron/tcp-work/pkg/tlog"
)

func run(ctx context.Context, cfg *config.Config) error {
	logger, err := tlog.New()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}

	app, err := NewApp(ctx, logger, cfg)
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	defer app.Close(ctx)

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("app run: %w", err)
	}

	return nil
}
