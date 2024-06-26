package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/lissteron/tcp-work/pkg/tlog"
)

type ProofOfWork struct {
	logger                tlog.Logger
	proofOfWorkDifficulty int
}

func NewProofOfWork(
	logger tlog.Logger,
	proofOfWorkDifficulty int,
) *ProofOfWork {
	return &ProofOfWork{
		logger:                logger,
		proofOfWorkDifficulty: proofOfWorkDifficulty,
	}
}

func (h *ProofOfWork) Generate() (string, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return "", fmt.Errorf("rand int: %w", err)
	}

	return strconv.FormatInt(num.Int64(), 10), nil
}

func (h *ProofOfWork) Validate(challenge, solution string) bool {
	var (
		hash    = sha256.Sum256([]byte(challenge + solution)) // Выбран как один из самых распространенных алгоритмов для PoW
		hashStr = hex.EncodeToString(hash[:])
	)

	h.logger.Infow(context.Background(), "server test", "solution", hashStr)

	return strings.HasPrefix(hashStr, strings.Repeat("0", h.proofOfWorkDifficulty))
}
