package config

import (
	"github.com/jinzhu/configor"
)

// Config - Application configuration
type Config struct {
	Meilisearch struct {
		Host   string `yaml:"host" default:"http://localhost:7700" env:"MEILISEARCH_HOST"`
		APIKey string `yaml:"api_key" default:"" env:"MEILISEARCH_API_KEY"`
	} `yaml:"meilisearch"`
}

// LoadConfig - Load configuration file
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := configor.New(&configor.Config{
		Debug:      false,
		Verbose:    false,
		Silent:     true,
		AutoReload: false,
	}).Load(cfg, path)
	return cfg, err
}
