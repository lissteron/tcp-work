package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/lissteron/tcp-work/config"
	"github.com/lissteron/tcp-work/internal/app/core/ports"
	"github.com/lissteron/tcp-work/internal/app/core/services"
	"github.com/lissteron/tcp-work/internal/app/handlers"
	"github.com/lissteron/tcp-work/pkg/tlog"
)

type App struct {
	logger tlog.Logger

	listener           net.Listener
	proofOfWorkHandler ports.TCPHandler

	errGroup    *errgroup.Group
	errGroupCtx context.Context //nolint:containedctx // need context in struct.
}

func NewApp(_ context.Context, logger tlog.Logger, cfg *config.Config) (*App, error) {
	listener, err := net.Listen("tcp", cfg.ListenerTCPAddr)
	if err != nil {
		return nil, fmt.Errorf("net listen: %w", err)
	}

	app := &App{
		logger:   logger,
		listener: listener,
	}

	var (
		proofOfWorkService = services.NewProofOfWork(logger, cfg.ProofOfWorkDifficulty)
		quoteService       = services.NewQuote(logger)
	)

	app.proofOfWorkHandler = handlers.NewProofOfWork(
		logger,
		proofOfWorkService,
		quoteService,
		cfg.WriteTimeout,
		cfg.ReadTimeout,
	)

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Infow(ctx, "start application")

	// error group
	a.errGroup, a.errGroupCtx = errgroup.WithContext(ctx)

	// start application.
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// start http
	a.errGroup.Go(func() error {
		var wg sync.WaitGroup

		for {
			conn, err := a.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break
				}

				a.logger.Errorw(ctx, "accept connection error", "error", err)

				continue
			}

			wg.Add(1)

			go func() {
				defer wg.Done()

				a.proofOfWorkHandler.Handle(conn)
			}()
		}

		wg.Wait()

		return nil
	})

	// wait exit signal
	select {
	case <-exit:
		a.logger.Infow(ctx, "stopping application")
	case <-a.errGroupCtx.Done():
		a.logger.Errorw(ctx, "stopping application with error")
	case <-ctx.Done():
		a.logger.Errorw(ctx, "stopping application with context canceled")
	}

	signal.Stop(exit)

	return nil
}

func (a *App) Close(ctx context.Context) {
	const (
		closeTimeout = 100 * time.Millisecond
		exitTimeout  = 30 * time.Second
	)

	// wait for orchestration
	time.Sleep(closeTimeout)

	ctx, cancel := context.WithTimeout(ctx, exitTimeout)
	defer cancel()

	if err := a.listener.Close(); err != nil {
		a.logger.Errorw(ctx, "close listener error", "error", err)
	}

	if err := a.errGroup.Wait(); err != nil {
		a.logger.Errorw(ctx, "wait error group error", "error", err)
	}

	a.logger.Infow(ctx, "service exited")
}
