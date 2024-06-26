package services

import (
	"math/rand"

	"github.com/lissteron/tcp-work/pkg/tlog"
)

type Quote struct {
	logger tlog.Logger
	quotes []string
}

func NewQuote(logger tlog.Logger) *Quote {
	return &Quote{
		logger: logger,
		quotes: []string{
			"Guard well your thoughts when alone and your words when accompanied.",
			"I like to listen. I have learned a great deal from listening carefully. Most people never listen.",
			"I think, that if the world were a bit more like ComicCon, it would be a better place.",
			"Quit being so hard on yourself. We are what we are; we love what we love. We don't need to justify it to anyone...",
			"Voice is not just the sound that comes from your throat, but the feelings that come from your words.",
		},
	}
}

func (s *Quote) GetRandom() string {
	return s.quotes[rand.Intn(len(s.quotes))] //nolint:gosec // not used for security purposes.
}
