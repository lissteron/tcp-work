package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/lissteron/tcp-work/config"
	"github.com/lissteron/tcp-work/internal/app/core/domains"
	"github.com/lissteron/tcp-work/pkg/tlog"
)

type App struct {
	logger tlog.Logger
	cfg    *config.Config
}

func NewApp(_ context.Context, logger tlog.Logger, cfg *config.Config) *App {
	return &App{
		logger: logger,
		cfg:    cfg,
	}
}

func (a *App) Run(ctx context.Context) error {
	const _bufSize = 1024

	conn, err := net.Dial("tcp", a.cfg.ServerTCPAddr)
	if err != nil {
		a.logger.Errorw(ctx, "dial error", "error", err)

		return fmt.Errorf("dial: %w", err)
	}

	defer conn.Close()

	buf := make([]byte, _bufSize)

	n, err := conn.Read(buf)
	if err != nil {
		a.logger.Errorw(ctx, "read error", "error", err)

		return fmt.Errorf("read: %w", err)
	}

	challengeMsg := string(buf[:n])

	if !strings.HasPrefix(challengeMsg, "CHALLENGE:") {
		a.logger.Errorw(ctx, "invalid challenge", "challenge", challengeMsg)

		return fmt.Errorf("%w: %s", domains.ErrInvalidChallenge, challengeMsg)
	}

	var (
		challenge = strings.TrimPrefix(challengeMsg, "CHALLENGE:")
		solution  = solvePoW(challenge, a.cfg.ProofOfWorkDifficulty)
	)

	if _, err := conn.Write([]byte(solution + "\n")); err != nil {
		a.logger.Errorw(ctx, "write error", "error", err)

		return fmt.Errorf("write: %w", err)
	}

	n, err = conn.Read(buf)
	if err != nil {
		a.logger.Errorw(ctx, "read error", "error", err)
	}

	a.logger.Infow(ctx, string(buf[:n]))

	return nil
}

func solvePoW(challenge string, difficulty int) string {
	var solution string

	for {
		solution = strconv.FormatInt(time.Now().UnixNano(), 10)

		var (
			hash    = sha256.Sum256([]byte(challenge + solution))
			hashStr = hex.EncodeToString(hash[:])
		)

		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			break
		}
	}

	return solution
}
