package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	Port      string `env:"PORT,required"`
	UserPort  string `env:"USER_PORT,required"`
	UserHost  string `env:"USER_HOST,required"`
	SecretKey string `env:"SECRET_KEY,required"`
}

func (c Config) String() string {
	var sb strings.Builder

	sb.WriteString("Auth Service Settings:\n")
	sb.WriteString(fmt.Sprintf("PORT: %s\n", c.Port))
	sb.WriteString(fmt.Sprintf("USER_PORT: %s\n", c.UserPort))
	sb.WriteString(fmt.Sprintf("USER_HOST: %s\n", c.UserHost))

	return sb.String()
}

func MustLoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(fmt.Errorf("parsing config error: %s", err.Error()))
	}
	return cfg
}
