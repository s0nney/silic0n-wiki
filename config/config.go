package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	Secret   string         `yaml:"secret"`
	Media    MediaConfig    `yaml:"media"`
}

type MediaConfig struct {
	UploadDir    string   `yaml:"upload_dir"`
	MaxFileSize  int64    `yaml:"max_file_size"`
	AllowedTypes []string `yaml:"allowed_types"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

var AppConfig *Config

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	AppConfig = &cfg
	return nil
}
