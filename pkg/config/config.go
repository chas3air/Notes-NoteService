package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	Grpc     GrpcConfig     `yaml:"grpc"`
	Rest     RestConfig     `yaml:"rest"`
	Metrics  MetricsConfig  `yaml:"metrics"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type GrpcConfig struct {
	Port int `yaml:"port"`
}

type RestConfig struct {
	Port int `yaml:"port"`
}

type MetricsConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Database  string `yaml:"dbname"`
	SSLMode   string `yaml:"sslmode"`
	TableName string `yaml:"table"`
}

func Load() (*Config, error) {
	configPath := fetchConfigPath()
	if configPath == "" {
		return nil, errors.New("config path is empty")
	}

	return LoadPath(configPath)
}

func LoadPath(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.New("cannot read config: " + err.Error())
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

// Data source name
func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}
