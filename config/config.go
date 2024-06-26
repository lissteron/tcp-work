package config

import (
	"time"

	"github.com/spf13/viper"
)

const _defaultProofOfWorkDifficulty = 2

type Config struct {
	ServerTCPAddr         string
	ListenerTCPAddr       string
	ProofOfWorkDifficulty int
	WriteTimeout          time.Duration
	ReadTimeout           time.Duration
}

func initConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("LISTENER_TCP_ADDR", ":8080")
	viper.SetDefault("PROOF_OF_WORK_DIFFICULTY", _defaultProofOfWorkDifficulty)
	viper.SetDefault("WRITE_TIMEOUT", time.Second)
	viper.SetDefault("READ_TIMEOUT", time.Second)
}

func NewConfig() *Config {
	initConfig()

	return &Config{
		ServerTCPAddr:         viper.GetString("SERVER_TCP_ADDR"),
		ListenerTCPAddr:       viper.GetString("LISTENER_TCP_ADDR"),
		ProofOfWorkDifficulty: viper.GetInt("PROOF_OF_WORK_DIFFICULTY"),
		WriteTimeout:          viper.GetDuration("WRITE_TIMEOUT"),
		ReadTimeout:           viper.GetDuration("READ_TIMEOUT"),
	}
}
