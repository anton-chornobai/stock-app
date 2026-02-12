package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	ServiceName string `yaml:"service_name"`
	DB DBConfig `yaml:"db"`
}

//conn string. Deriving string from the params ↓

func (d DBConfig) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

func GetConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		return  nil, fmt.Errorf("CONFIG_PATH environment variable is not set")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
