package handlers

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/lissteron/tcp-work/internal/app/core/ports"
	"github.com/lissteron/tcp-work/pkg/tlog"
)

type ProofOfWork struct {
	logger             tlog.Logger
	proofOfWorkService ports.ProofOfWorkService
	quoteService       ports.QuoteService
	writeTimeout       time.Duration
	readTimeout        time.Duration
}

func NewProofOfWork(
	logger tlog.Logger,
	proofOfWorkService ports.ProofOfWorkService,
	quoteService ports.QuoteService,
	writeTimeout time.Duration,
	readTimeout time.Duration,
) *ProofOfWork {
	return &ProofOfWork{
		logger:             logger,
		proofOfWorkService: proofOfWorkService,
		quoteService:       quoteService,
		writeTimeout:       writeTimeout,
		readTimeout:        readTimeout,
	}
}

func (h *ProofOfWork) Handle(conn net.Conn) {
	const _bufSize = 1024

	defer conn.Close()

	ctx := context.Background()

	challenge, err := h.proofOfWorkService.Generate()
	if err != nil {
		h.logger.Errorw(ctx, "generate error", err)

		// тут можно добавить отправку ответа что произошла ошибка

		return
	}

	if err := conn.SetWriteDeadline(time.Now().Add(h.writeTimeout)); err != nil {
		h.logger.Errorw(ctx, "set write deadline error", err)

		return
	}

	if _, err := conn.Write([]byte("CHALLENGE:" + challenge)); err != nil {
		h.logger.Errorw(ctx, "write error", err)

		return
	}

	buf := make([]byte, _bufSize)

	if err := conn.SetReadDeadline(time.Now().Add(h.readTimeout)); err != nil {
		h.logger.Errorw(ctx, "set read deadline error", err)

		return
	}

	n, err := conn.Read(buf)
	if err != nil {
		h.logger.Errorw(ctx, "read error", err)

		return
	}

	if err := conn.SetWriteDeadline(time.Now().Add(h.writeTimeout)); err != nil {
		h.logger.Errorw(ctx, "set write deadline error", err)

		return
	}

	solution := strings.TrimSpace(string(buf[:n]))

	if h.proofOfWorkService.Validate(challenge, solution) {
		if _, err := conn.Write([]byte(h.quoteService.GetRandom() + "\n")); err != nil {
			h.logger.Errorw(ctx, "write error", err)

			return
		}

		return
	}

	if _, err := conn.Write([]byte("ERROR: Invalid PoW")); err != nil {
		h.logger.Errorw(ctx, "write error", err)

		return
	}
}
