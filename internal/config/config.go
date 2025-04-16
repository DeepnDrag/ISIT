package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type DB struct {
	Password       string `yaml:"password"`
	Host           string `yaml:"host"`
	SSLMode        string `yaml:"ssl_mode"`
	User           string `yaml:"user"`
	Port           string `yaml:"port"`
	MigrationsPath string `yaml:"migrations_path"`
	Name           string `yaml:"db_name"`
}

type Logger struct {
	Sink  string `yaml:"sink"`
	Level string `yaml:"level"`
}

type Server struct {
	URL string `yaml:"url"`
}

type JWT struct {
	SecretKey      string `mapstructure:"secret_key"`
	ExpirationTime int    `mapstructure:"expiration_time"`
}

type AppConfig struct {
	DB     DB     `yaml:"database"`
	Logger Logger `yaml:"logger"`
	Server Server `yaml:"server"`
	JWT    JWT    `yaml:"jwt"`
}

func New(path string) (*AppConfig, error) {
	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var appConfig AppConfig

	if err := yaml.Unmarshal(yamlConfig, &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
