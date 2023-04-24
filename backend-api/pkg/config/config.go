package config

import (
	"bytes"
	"encoding/json"
	"os"
)

type Config struct {
	DB     `json:"database"`
	Server `json:"http"`
}

type DB struct {
	DBhost string `json:"database_host"`
	DBport int    `json:"database_port"`
	DBuser string `json:"database_user"`
	DBpass string `json:"database_password"`
	DBname string `json:"database_name"`
}

type Server struct {
	Address      string `json:"address"`
	Port         string `json:"port"`
	IdleTimeout  int    `json:"idle_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	ReadTimeout  int    `json:"read_timeout"`
}

func Init(path string) (*Config, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
