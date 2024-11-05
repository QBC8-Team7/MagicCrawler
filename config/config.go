package config

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Server struct {
	Host string
	Port string
}

type Bot struct {
	Token string
}

type Database struct {
	ConnectionURI string
}

type Crawler struct {
	CrawlTime uint
}

type Config struct {
	Server
	Bot
	Database
	Crawler
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
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
