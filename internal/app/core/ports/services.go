package ports

type ProofOfWorkService interface {
	Generate() (string, error)
	Validate(challenge, solution string) bool
}

type QuoteService interface {
	GetRandom() string
}
