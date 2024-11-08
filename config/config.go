package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Logger struct {
	Level string
	Path  string
}
type Server struct {
	Host       string
	Port       string
	Mode       string // development or production
	AppVersion string
}

type Bot struct {
	Token string
}

type Postgres struct {
	ConnectionURI      string
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSslmode  bool
}

type Crawler struct {
	CrawlTime uint
}

type Config struct {
	Server
	Bot
	Postgres
	Crawler
	Logger
}

func LoadConfig() (*Config, error) {
	var c Config
	v := viper.New()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("unable to get the current directory")
	}
	configDir := filepath.Dir(filename)

	v.SetConfigFile(filepath.Join(configDir, "config.yml"))
	v.SetConfigType("yaml")

	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
