package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerMode string

const (
	Development ServerMode = "development"
	Production  ServerMode = "production"
)

type Logger struct {
	Level   string
	Path    string
	SysPath string
}
type Server struct {
	Host       string
	Port       string
	Mode       ServerMode
	AppVersion string
}

type Bot struct {
	Token string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SslMode  bool
}

type Crawler struct {
	Time           int
	Seeds          []map[string]string
	GeneralLogPath string
	MetricLogPath  string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Config struct {
	Server
	Bot
	Postgres
	Crawler
	Logger
	Redis
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
